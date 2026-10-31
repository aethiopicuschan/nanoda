package strings

import (
	"unsafe"
)

// 先頭のアドレスから0が出るまでの文字列を取得する
func GoString(ptr *byte) string {
	length := 0
	for *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(length))) != 0 {
		length++
	}
	return unsafe.String((*byte)(unsafe.Pointer(ptr)), length)
}
