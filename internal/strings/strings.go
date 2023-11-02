package strings

import (
	"unsafe"
)

// 終端文字を付けてポインタにして返す
func CString(s string) uintptr {
	bytes := []byte(s)
	if len(bytes) == 0 || bytes[len(bytes)-1] != 0 {
		bytes = append(bytes, 0)
	}
	return *(*uintptr)(unsafe.Pointer(&bytes))
}

// 先頭のアドレスから0が出るまでの文字列を取得する
func GoString(ptr *byte) string {
	length := 0
	for *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + uintptr(length))) != 0 {
		length++
	}
	return unsafe.String((*byte)(unsafe.Pointer(ptr)), length)
}

// 先頭のアドレスから0が出るまでの文字列を取得する
func GoStringFromUintptr(ptr uintptr) string {
	return GoString((*byte)(unsafe.Pointer(ptr)))
}
