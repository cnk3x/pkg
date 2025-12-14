package speedx

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cnk3x/pkg/urlx"
)

// Downloader 是一个高级下载器
type Downloader struct {
	URL              string
	SaveDir          string
	FileName         string // 可选
	MaxThreads       int
	MaxChunkSize     int64 // 字节
	ProgressCallback func(progress float64, speed float64)
	MaxRetryTimes    int
	RetryInterval    time.Duration
	ContinueOnError  bool

	// 内部状态
	fileSize         int64
	tempFilePath     string
	infoFilePath     string
	downloadedChunks map[int]bool // 记录已下载的块索引
	mutex            sync.Mutex
	wg               sync.WaitGroup
	progressChan     chan int64
	speedChan        chan int64
	doneChan         chan struct{}
	err              error
}

// NewDownloader 创建一个新的下载器实例
func NewDownloader(url string, saveDir string, options ...func(*Downloader)) *Downloader {
	d := &Downloader{
		URL:              url,
		SaveDir:          saveDir,
		MaxThreads:       5,
		MaxChunkSize:     1024 * 1024, // 1MB
		MaxRetryTimes:    3,
		RetryInterval:    5 * time.Second,
		ContinueOnError:  true,
		downloadedChunks: make(map[int]bool),
		progressChan:     make(chan int64, 100),
		speedChan:        make(chan int64, 100),
		doneChan:         make(chan struct{}),
	}

	// 应用可选参数
	for _, option := range options {
		option(d)
	}

	return d
}

// Start 开始下载
func (d *Downloader) Start(ctx context.Context) error {
	// 1. 检查URL和服务器支持
	if err := d.checkURL(ctx); err != nil {
		return err
	}

	// // 2. 确定文件名
	// if err := d.determineFileName(); err != nil {
	// 	return err
	// }

	// 3. 准备文件路径
	if err := d.prepareFilePaths(); err != nil {
		return err
	}

	// 4. 加载已下载的块信息
	if err := d.loadDownloadInfo(); err != nil {
		return err
	}

	// 5. 启动进度报告goroutine
	go d.reportProgress()

	// 6. 开始下载
	if d.fileSize == 0 {
		// 空文件
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			d.mutex.Lock()
			defer d.mutex.Unlock()
			if err := os.WriteFile(d.tempFilePath, []byte{}, 0644); err != nil {
				d.err = err
			}
		}()
	} else if d.supportsResume() {
		// 支持断点续传，多线程下载
		d.multiThreadedDownload(ctx)
	} else {
		// 不支持断点续传，单线程下载
		d.singleThreadedDownload(ctx)
	}

	// 7. 等待下载完成
	d.wg.Wait()
	close(d.progressChan)
	close(d.speedChan)
	<-d.doneChan

	// 8. 检查错误
	if d.err != nil && !d.ContinueOnError {
		return d.err
	}

	// 9. 重命名临时文件为最终文件名
	if err := os.Rename(d.tempFilePath, filepath.Join(d.SaveDir, d.FileName)); err != nil {
		return err
	}

	// 10. 删除下载信息文件
	os.Remove(d.infoFilePath)

	return nil
}

// checkURL 检查URL是否有效并获取文件信息
func (d *Downloader) checkURL(ctx context.Context) error {
	var lastReq *http.Request
	return urlx.Windows().Method(urlx.MethodHead).Url(d.URL).
		Client(func(cli *http.Client) error {
			cli.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return errors.New("stopped after 10 redirects")
				}
				lastReq = req
				return nil
			}
			return nil
		}).
		Process(ctx, func(resp *http.Response) error {
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("HTTP request failed with status: %s", resp.Status)
			}
			// 获取文件大小
			sizeStr := resp.Header.Get("Content-Length")
			if sizeStr != "" {
				size, err := strconv.ParseInt(sizeStr, 10, 64)
				if err != nil {
					return err
				}
				d.fileSize = size
			}

			if lastReq != nil {
				resp.Request = lastReq
			}
			d.determineFileName(resp)
			return nil
		})
}

// determineFileName 确定保存的文件名
func (d *Downloader) determineFileName(resp *http.Response) {
	if d.FileName != "" {
		return
	}

	// 尝试从Content-Disposition头获取
	disposition := resp.Header.Get("Content-Disposition")
	if disposition != "" {
		parts := strings.Split(disposition, "filename=")
		if len(parts) > 1 {
			d.FileName = strings.Trim(parts[1], `"`)
			return
		}
	}

	// 尝试从URL路径中获取
	if filename := filepath.Base(resp.Request.URL.Path); filename != "" && filename != "." && filename != "/" {
		d.FileName = filename
		return
	}

	// 作为最后的手段，使用当前时间戳
	d.FileName = strconv.FormatInt(time.Now().UnixNano(), 10)
}

// prepareFilePaths 准备临时文件和信息文件的路径
func (d *Downloader) prepareFilePaths() error {
	// 创建保存目录
	if err := os.MkdirAll(d.SaveDir, 0755); err != nil {
		return err
	}

	// 设置临时文件路径
	d.tempFilePath = filepath.Join(d.SaveDir, d.FileName+".downloading")
	d.infoFilePath = d.tempFilePath + ".info"

	// 创建临时文件（如果不存在）
	if _, err := os.Stat(d.tempFilePath); os.IsNotExist(err) {
		file, err := os.Create(d.tempFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 如果知道文件大小，预先分配空间
		if d.fileSize > 0 {
			if err := file.Truncate(d.fileSize); err != nil {
				return err
			}
		}
	}

	return nil
}

// loadDownloadInfo 加载已下载的块信息
func (d *Downloader) loadDownloadInfo() error {
	if _, err := os.Stat(d.infoFilePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(d.infoFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		chunkIdx, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		d.downloadedChunks[chunkIdx] = true
	}

	return scanner.Err()
}

// saveDownloadInfo 保存已下载的块信息
func (d *Downloader) saveDownloadInfo(chunkIdx int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.downloadedChunks[chunkIdx] = true

	file, err := os.OpenFile(d.infoFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, chunkIdx)
	return err
}

// supportsResume 检查是否支持断点续传
func (d *Downloader) supportsResume() bool {
	resp, err := http.Head(d.URL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.Header.Get("Accept-Ranges") == "bytes" && d.fileSize > 0
}

// multiThreadedDownload 多线程下载
func (d *Downloader) multiThreadedDownload(ctx context.Context) {
	numChunks := int((d.fileSize + d.MaxChunkSize - 1) / d.MaxChunkSize)
	chunksToDownload := []int{}

	// 找出还未下载的块
	for i := range numChunks {
		if !d.downloadedChunks[i] {
			chunksToDownload = append(chunksToDownload, i)
		}
	}

	// 如果所有块都已下载，直接返回
	if len(chunksToDownload) == 0 {
		return
	}

	// 限制并发数
	sem := make(chan struct{}, d.MaxThreads)

	for _, chunkIdx := range chunksToDownload {
		sem <- struct{}{}
		d.wg.Add(1)
		go func(idx int) {
			defer func() {
				d.wg.Done()
				<-sem
			}()

			start := int64(idx) * d.MaxChunkSize
			end := start + d.MaxChunkSize - 1
			if end >= d.fileSize {
				end = d.fileSize - 1
			}

			// 下载并写入块
			if err := d.downloadAndWriteChunk(ctx, idx, start, end); err != nil {
				d.mutex.Lock()
				if d.err == nil {
					d.err = err
				}
				d.mutex.Unlock()
			}
		}(chunkIdx)
	}
}

// singleThreadedDownload 单线程下载
func (d *Downloader) singleThreadedDownload(ctx context.Context) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()

		// 如果文件已存在且大小匹配，直接返回
		if fileInfo, err := os.Stat(d.tempFilePath); err == nil && fileInfo.Size() == d.fileSize {
			return
		}

		err := urlx.Windows().Url(d.URL).
			Process(ctx, func(resp *http.Response) error {
				file, err := os.OpenFile(d.tempFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				defer file.Close()

				buf := make([]byte, 4096)
				var total int64
				for {
					n, err := resp.Body.Read(buf)
					if n > 0 {
						if _, e := file.Write(buf[:n]); err != nil {
							return e
						}
						total += int64(n)
						d.progressChan <- int64(n)
						d.speedChan <- int64(n)
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						return err
					}
				}
				return nil
			})

		if err != nil {
			d.mutex.Lock()
			d.err = err
			d.mutex.Unlock()
			return
		}
	}()
}

// downloadAndWriteChunk 下载并写入一个块
func (d *Downloader) downloadAndWriteChunk(ctx context.Context, chunkIdx int, start, end int64) error {
	var lastErr error

	for i := 0; i <= d.MaxRetryTimes; i++ {
		// 如果已经有其他错误且不允许继续，直接返回
		d.mutex.Lock()
		if d.err != nil && !d.ContinueOnError {
			d.mutex.Unlock()
			return d.err
		}
		d.mutex.Unlock()

		var errRetry = errors.New("retry")

		err := urlx.Windows().Url(d.URL).
			HeaderSet("Range", fmt.Sprintf("bytes=%d-%d", start, end)).
			Process(ctx, func(resp *http.Response) error {
				// 检查状态码
				if resp.StatusCode != http.StatusPartialContent {
					resp.Body.Close()
					lastErr = fmt.Errorf("unexpected status code: %s", resp.Status)
					return errRetry
				}

				// 打开文件
				file, err := os.OpenFile(d.tempFilePath, os.O_WRONLY, 0644)
				if err != nil {
					resp.Body.Close()
					lastErr = err
					return errRetry
				}

				// 定位到写入位置
				if _, err := file.Seek(start, io.SeekStart); err != nil {
					file.Close()
					resp.Body.Close()
					lastErr = err
					return errRetry
				}

				// 下载并写入数据
				buf := make([]byte, 4096)
				var total int64

				for {
					n, err := resp.Body.Read(buf)
					if n > 0 {
						if _, e := file.Write(buf[:n]); err != nil {
							file.Close()
							resp.Body.Close()
							lastErr = e
							return errRetry
						}
						total += int64(n)
						d.progressChan <- int64(n)
						d.speedChan <- int64(n)
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						file.Close()
						resp.Body.Close()
						lastErr = err
						return errRetry
					}
				}

				// 关闭文件和响应体
				file.Close()
				resp.Body.Close()

				// 检查是否下载了完整的块
				if total != end-start+1 {
					lastErr = fmt.Errorf("incomplete chunk download: expected %d bytes, got %d", end-start+1, total)
					return errRetry
				}

				// 保存块信息
				if err := d.saveDownloadInfo(chunkIdx); err != nil {
					lastErr = err
					return errRetry
				}

				// 下载成功
				return nil
			})

		if err == errRetry {
			time.Sleep(d.RetryInterval)
		} else {
			return nil
		}
	}

	return fmt.Errorf("failed to download chunk after %d retries: %v", d.MaxRetryTimes, lastErr)
}

// reportProgress 报告下载进度和速度
func (d *Downloader) reportProgress() {
	if d.ProgressCallback == nil {
		return
	}

	var totalDownloaded int64
	var speedSamples []int64
	var lastReportTime time.Time

	for {
		select {
		case bytes, ok := <-d.progressChan:
			if !ok {
				return
			}
			totalDownloaded += bytes

		case bytes, ok := <-d.speedChan:
			if !ok {
				return
			}
			speedSamples = append(speedSamples, bytes)

		case <-time.After(1 * time.Second):
			// 每秒计算一次速度
			now := time.Now()
			duration := now.Sub(lastReportTime)
			lastReportTime = now

			// 计算平均速度
			var totalSpeed int64
			for _, s := range speedSamples {
				totalSpeed += s
			}
			speed := float64(totalSpeed) / duration.Seconds()

			// 计算进度
			progress := float64(totalDownloaded) / float64(d.fileSize) * 100
			if progress > 100 {
				progress = 100
			}

			// 调用回调函数
			d.ProgressCallback(progress, speed)

			// 清空速度样本
			speedSamples = speedSamples[:0]

		case <-d.doneChan:
			return
		}
	}
}

// 可选参数设置函数
func WithFileName(name string) func(*Downloader) {
	return func(d *Downloader) {
		d.FileName = name
	}
}

func WithMaxThreads(threads int) func(*Downloader) {
	return func(d *Downloader) {
		d.MaxThreads = threads
	}
}

func WithMaxChunkSize(size int64) func(*Downloader) {
	return func(d *Downloader) {
		d.MaxChunkSize = size
	}
}

func WithProgressCallback(callback func(float64, float64)) func(*Downloader) {
	return func(d *Downloader) {
		d.ProgressCallback = callback
	}
}

func WithMaxRetryTimes(times int) func(*Downloader) {
	return func(d *Downloader) {
		d.MaxRetryTimes = times
	}
}

func WithRetryInterval(interval time.Duration) func(*Downloader) {
	return func(d *Downloader) {
		d.RetryInterval = interval
	}
}

func WithContinueOnError(continueOnError bool) func(*Downloader) {
	return func(d *Downloader) {
		d.ContinueOnError = continueOnError
	}
}
