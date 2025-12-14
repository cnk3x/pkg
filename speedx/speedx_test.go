package speedx

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDownload(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// 示例用法
	url := "https://ftp.sjtu.edu.cn/ubuntu-cd/24.04.3/ubuntu-24.04.3-live-server-amd64.iso"
	saveDir := "./downloads"

	// 创建下载器
	downloader := NewDownloader(url, saveDir,
		WithMaxThreads(10),
		WithMaxChunkSize(2*1024*1024), // 2MB
		WithProgressCallback(func(progress float64, speed float64) {
			t.Logf("\r\rProgress: %.2f%%, Speed: %.2f MB/s", progress, speed/(1024*1024))
		}),
		WithMaxRetryTimes(5),
		WithRetryInterval(3*time.Second),
		WithContinueOnError(true),
	)

	// 开始下载
	fmt.Println("Starting download...")
	err := downloader.Start(ctx)
	if err != nil {
		fmt.Printf("\nDownload failed: %v\n", err)
		return
	}

	fmt.Println("\nDownload completed successfully!")
}
