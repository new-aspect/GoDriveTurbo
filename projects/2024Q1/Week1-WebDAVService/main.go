package main

import (
	"golang.org/x/net/webdav"
	"log"
	"net/http"
)

func main() {
	// 指定要共享的本地文件夹路径
	dir := "/Users/zhaon/Movies"

	// 创建一个新的WebDAV处理器
	dav := &webdav.Handler{
		Prefix:     "/dav",
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}

	// 设置HTTP服务器使用WebDAV处理器
	http.HandleFunc("/", dav.ServeHTTP)

	// 启动HTTP服务器在端口3100上监听
	log.Println("WebDAV server running on http://localhost:15244/")
	err := http.ListenAndServe(":15244", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
