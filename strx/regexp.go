package strx

import (
	"regexp"
	"sync"
)

func createRegPool() func(expr string) (*regexp.Regexp, error) {
	type CompileResult struct {
		re  *regexp.Regexp
		err error
	}
	regPool := map[string]CompileResult{}
	mu := &sync.Mutex{}
	return func(expr string) (*regexp.Regexp, error) {
		mu.Lock()
		defer mu.Unlock()
		if v, find := regPool[expr]; find {
			return v.re, v.err
		}
		re, err := regexp.Compile(expr)
		regPool[expr] = CompileResult{re: re, err: err}
		return re, nil
	}
}

var RegexpE = createRegPool()

func Regexp(expr string) *regexp.Regexp { re, _ := RegexpE(expr); return re }
