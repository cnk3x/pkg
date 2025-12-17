package arrx

// 可变参数固定
func Variadic[T any](arr ...T) []T { return arr }

// Len 返回切片 s 的长度
//
// 参数:
//   - s: 类型为 S 的切片，其中 S 是 []T 的类型别名，T 可以为任意类型
//
// 返回值:
//   - int: 切片 s 的元素个数
func Len[S ~[]T, T any](s S) int { return len(s) }
