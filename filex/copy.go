package filex

import (
	"context"
	"errors"
	"io"
)

var ErrInvalidWrite = errors.New("invalid write result")
var ErrShortWrite = io.ErrShortWrite

func Copy(ctx context.Context, dst io.Writer, src io.Reader, progress func(int64)) (err error) {
	size := 32 * 1024
	if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
		if l.N < 1 {
			size = 1
		} else {
			size = int(l.N)
		}
	}
	buf := make([]byte, size)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if nw < 0 || nr < nw {
					nw = 0
					if ew == nil {
						ew = ErrInvalidWrite
					}
				}

				if progress != nil {
					progress(int64(nw))
				}

				if ew != nil {
					return ew
				}
				if nr != nw {
					return ErrShortWrite
				}
			}

			if er != nil {
				if er != io.EOF {
					return er
				}
				return nil
			}
		}
	}
}
