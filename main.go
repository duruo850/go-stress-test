// Package main go 实现的压测工具
package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"go-stress-test/conf"
	"go-stress-test/model"
	"go-stress-test/server"
)

// main go 实现的压测工具
// 编译可执行文件
//go:generate go build main.go
func main() {
	runtime.GOMAXPROCS(1)

	if conf.Concurrency == 0 || conf.TotalNumber == 0 || (conf.RequestURL == "" && conf.Path == "") {
		fmt.Printf("示例: go run main.go -c 1 -n 1 -u https://www.baidu.com/ \n")
		fmt.Printf("压测地址或curl路径必填 \n")
		fmt.Printf("当前请求参数: -c %d -n %d -d %v -u %s \n", conf.Concurrency, conf.TotalNumber, conf.DebugStr, conf.RequestURL)
		flag.Usage()
		return
	}
	debug := strings.ToLower(conf.DebugStr) == "true"
	request, err := model.NewRequest(conf.RequestURL, conf.Verify, conf.Code, 0, debug, conf.Path, conf.Headers, conf.Body, conf.MaxCon, conf.Http2, conf.Keepalive)
	if err != nil {
		fmt.Printf("参数不合法 %v \n", err)
		return
	}
	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", conf.Concurrency, conf.TotalNumber)
	request.Print()
	// 开始处理
	server.Dispose(conf.Concurrency, conf.TotalNumber, request)
	return
}
