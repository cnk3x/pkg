package unitx

import "fmt"

// Mebibyte 1024 进制
const (
	ki = 1024 << (10 * iota)
	mi
	gi
	ti
	pi
	ei
)

// Megabyte 1000 进制
const (
	k = 1000
	m = k * 1000
	g = m * 1000
	t = g * 1000
	p = t * 1000
	e = p * 1000
)

var (
	n_mebibyte = []uint64{1, ki, mi, gi, ti, pi, ei}
	n_megabyte = []uint64{1, k, m, g, t, p, e}
	u_mebibyte = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	u_megabyte = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
)

type N interface {
	~int64 | ~int32 | ~int16 | ~int8 | ~int |
		~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint |
		~float64 | ~float32
}

// Mebibyte 1024 进制
func Mebibyte[T N](size T) string {
	for i := len(n_mebibyte) - 1; i >= 0; i-- {
		if uint64(size) >= n_mebibyte[i] {
			return fmt.Sprintf("%.2f %s", float64(size)/float64(n_mebibyte[i]), u_mebibyte[i])
		}
	}
	return fmt.Sprintf("%d B", uint64(size))
}

// Megabyte 1000 进制
func Megabyte[T N](size T) string {
	for i := len(n_megabyte) - 1; i >= 0; i-- {
		if uint64(size) >= n_megabyte[i] {
			return fmt.Sprintf("%.2f %s", float64(size)/float64(n_megabyte[i]), u_megabyte[i])
		}
	}
	return fmt.Sprintf("%d B", uint64(size))
}
