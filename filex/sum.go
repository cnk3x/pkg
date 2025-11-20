package filex

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/cnk3x/gopkg/errx"
)

func CheckSum(ctx context.Context, fn, digest string) bool {
	return errx.May(CheckSumE(ctx, fn, digest))
}

func CheckSumE(ctx context.Context, fn, digest string) (pass bool, err error) {
	if digest == "" {
		return
	}

	hashType, digest, ok := strings.Cut(digest, ":")
	if !ok {
		digest = hashType
		switch l := len(digest); l {
		case 32:
			hashType = "md5" //202cb962ac59075b964b07152d234b70
		case 40:
			hashType = "sha1" //04aab3efdead9c5276644e352e8d11858f06f0bb
		case 64:
			hashType = "sha256" //a835c4f79afae98d8537d7a83928d95975985f6c731788ae8d08d64c58fdfc5d
		default:
			err = errx.Errorf("unknown digest length: %d", l)
			return
		}
	}

	d, e := CalcSum(ctx, fn, hashType)
	if err = e; err != nil {
		return
	}

	pass = digest == d
	return
}

func CalcSum(ctx context.Context, fn, hashType string) (digest string, err error) {
	var h hash.Hash
	switch hashType {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	default:
		err = errx.Errorf("unknown hash type: %s", hashType)
		return
	}

	var f *os.File
	if f, err = os.Open(fn); err != nil {
		err = errx.Errorf("open file: %w", err)
		return
	}

	if _, err = io.Copy(h, ReaderContext(ctx, f)); err != nil {
		err = errx.Errorf("calc file: %w", err)
		return
	}

	digest = hex.EncodeToString(h.Sum(nil))
	return
}

func ReaderContext(ctx context.Context, r io.Reader) io.Reader {
	return rFunc(func(b []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, errx.Errorf("break with context done: %w", ctx.Err())
		default:
			return r.Read(b)
		}
	})
}

type rFunc func([]byte) (int, error)

func (r rFunc) Read(b []byte) (int, error) { return r(b) }
