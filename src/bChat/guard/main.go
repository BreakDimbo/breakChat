package main

import (
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// TODO: 启动日志

	// TODO: 启动 gRPC 服务
	ListenAndServe()
}
