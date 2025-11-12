package strx

import (
	"unsafe"
)

// Atob string to bytes (unsafe zero copy)
func Atob(s string) []byte { return unsafe.Slice(unsafe.StringData(s), len(s)) }

// Btoa bytes to string (unsafe zero copy)
func Btoa(s string) []byte { return unsafe.Slice(unsafe.StringData(s), len(s)) }
