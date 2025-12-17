package arrx

// Sum 对切片中的数值元素进行求和运算
//
// 参数:
//   - s: 包含数值类型的切片
//
// 返回值:
//   - R: 求和结果，类型为R
func Sum[T Num](s []T) T {
	return Reduce(s, func(r T, t T) T { return r + t }, 0)
}

// SumBy 根据指定函数对切片中每个元素转换为数值后进行求和运算
//
// 参数:
//   - s: 包含任意类型的切片
//   - atoi: 转换函数，将元素T转换为数值类型R
//
// 返回值:
//   - R: 求和结果，类型为R
func SumBy[T any, R Num](s []T, atoi func(T) R) R {
	return Reduce(s, func(r R, t T) R { return r + atoi(t) }, 0)
}

// SumIndex 根据指定函数（包含元素索引）对切片中每个元素转换为数值后进行求和运算
//
// 参数:
//   - s: 包含任意类型的切片
//   - atoi: 转换函数，将元素T和其索引转换为数值类型R
//
// 返回值:
//   - R: 求和结果，类型为R
func SumIndex[T any, R Num](s []T, atoi func(T, int) R) R {
	return ReduceIndex(s, func(r R, t T, i int) R { return r + atoi(t, i) }, 0)
}
