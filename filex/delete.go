package filex

import "os"

// Delete 删除文件或者目录，不存在不报错
func Delete(fn string) (err error) {
	err = os.Remove(fn)
	if os.IsNotExist(err) {
		err = nil
	}
	return
}

type DeleteOptions struct {
	RemoveChildren bool //删除
}
