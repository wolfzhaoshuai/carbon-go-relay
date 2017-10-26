package utils

import "unsafe"

//StringToBytes convert string to bytes without new memory allocation
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//BytesToStr convert bytes to string without new memory allocation
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
