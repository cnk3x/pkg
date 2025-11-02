package filex

import "io"

// ProgressCopy 返回一个 ProcessFunc，在拷贝过程中通过 progress 回调报告已拷贝字节数
func ProgressCopy(w io.Writer, r io.Reader, progress ...ProgressFunc) (err error) {
	_, err = io.Copy(ProgressWriter(w, progress...), r)
	return
}

// ProgressWriter 返回一个 io.Writer，在写入过程中通过 progress 回调报告已写入字节数
func ProgressWriter(w io.Writer, progress ...ProgressFunc) io.Writer {
	if len(progress) == 0 {
		return w
	}
	return &progressWriter{w: w, progress: progress}
}

type progressWriter struct {
	w        io.Writer
	progress []ProgressFunc
}

func (p *progressWriter) Write(b []byte) (n int, err error) {
	n, err = p.w.Write(b)
	if err == nil {
		for _, progress := range p.progress {
			progress(int64(n))
		}
	}
	return
}

// ProgressReader 返回一个 io.Reader，在读取过程中通过 progress 回调报告已读取字节数
func ProgressReader(r io.Reader, progress ...ProgressFunc) io.Reader {
	if len(progress) == 0 {
		return r
	}
	return &progressReader{r: r, progress: progress}
}

type progressReader struct {
	r        io.Reader
	progress []ProgressFunc
}

func (p *progressReader) Read(b []byte) (n int, err error) {
	n, err = p.r.Read(b)
	if err == nil {
		for _, progress := range p.progress {
			progress(int64(n))
		}
	}
	return
}
