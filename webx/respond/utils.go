package respond

import (
	"bytes"
	"encoding/json"

	"github.com/cnk3x/pkg/syncx"
)

var bufPool = syncx.NewPool[bytes.Buffer]()

func Escape(src []byte) []byte {
	buf := bufPool.Get()
	defer bufPool.Put(buf)
	json.HTMLEscape(buf, src)
	return buf.Bytes()
}
