package cache

import (
	"unsafe"
)

const (
	sizeOfInt64 = int(unsafe.Sizeof(int64(0)))
	sizeOfPtr   = int(unsafe.Sizeof(uintptr(0)))
)

func sizeOfString(s string) int {
	return len(s) + sizeOfPtr
}

func sizeOfMap(m map[string]Item) int {
	var size int
	for key, item := range m {
		size += sizeOfString(key)
		size += sizeOfString(item.Value) + sizeOfInt64
	}
	size += sizeOfPtr
	return size
}
