package main

import "fmt"

// 假设的drive.Drive 接口
type Drive interface{}

// StorageHook 函数类型
type StorageHook func(typ string, storage Drive)

// 存储所有注册的 storageHook 函数
var storageHooks = make([]StorageHook, 0)

// 注册 storageHook 函数
func RegisterStorageHook(hook StorageHook) {
	storageHooks = append(storageHooks, hook)
}

// 调用所有注册的 StorageHook 函数
func callStorageHooks(typ string, storage Drive) {
	for _, hook := range storageHooks {
		hook(typ, storage)
	}
}

// 实现一个StorageHook函数

func myStorageHook(typ string, storage Drive) {
	fmt.Printf("StorageHook called with type: %s and storage: %v\n", typ, storage)
}

func main() {
	// 注册StorageHook
	RegisterStorageHook(myStorageHook)

	// 假设的 storage 对象
	var storage Drive

	// 调用 StorageHooks
	callStorageHooks("example", storage)
}
