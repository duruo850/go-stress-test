// Package server 压测启动
package server

import (
	"fmt"
	"go-stress-test/conf"
	"sync"
	"time"

	httplongclinet "go-stress-test/server/client/http_longclinet"

	"go-stress-test/model"
	"go-stress-test/server/client"
	"go-stress-test/server/golink"
	"go-stress-test/server/statistics"
	"go-stress-test/server/verify"
)

// init 注册验证器
func init() {

	// http
	model.RegisterVerifyHTTP("statusCode", verify.HTTPStatusCode)
	model.RegisterVerifyHTTP("json", verify.HTTPJson)

	// webSocket
	model.RegisterVerifyWebSocket("json", verify.WebSocketJSON)
}

// Dispose 处理函数
func Dispose(concurrency, totalNumber uint64, request *model.Request) {
	// 设置接收数据缓存
	ch := make(chan *model.RequestResults, 1000)
	var (
		wg          sync.WaitGroup // 发送数据完成
		wgReceiving sync.WaitGroup // 数据处理完成
	)
	wgReceiving.Add(1)
	go statistics.ReceivingResults(concurrency, ch, &wgReceiving)

	if request.Keepalive {
		httplongclinet.CreateLangHttpClient(request)
	}

	wg.Add(int(concurrency))
	for i := uint64(0); i < concurrency; i++ {
		switch request.Form {
		case model.FormTypeHTTP:
			go golink.HTTP(i, ch, totalNumber, &wg, request)
		case model.FormTypeWebSocket:
			switch conf.ConnectionMode {
			case 1:
				// 连接以后再启动协程
				ws := client.NewWebSocket(request.URL)
				err := ws.GetConn()
				if err != nil {
					fmt.Println("连接失败:", i, err)
					continue
				}
				go golink.WebSocket(i, ch, totalNumber, &wg, request, ws)
			case 2:
				// 并发建立长链接
				go func(i uint64) {
					// 连接以后再启动协程
					ws := client.NewWebSocket(request.URL)
					err := ws.GetConn()
					if err != nil {
						fmt.Println("连接失败:", i, err)
						return
					}
					golink.WebSocket(i, ch, totalNumber, &wg, request, ws)
				}(i)
			default:
				data := fmt.Sprintf("不支持的类型:%d", conf.ConnectionMode)
				panic(data)
			}
		case model.FormTypeGRPC:
			// 连接以后再启动协程
			ws := client.NewGrpcSocket(request.URL)
			err := ws.Link()
			if err != nil {
				fmt.Println("连接失败:", i, err)
				continue
			}
			go golink.Grpc(i, ch, totalNumber, &wg, request, ws)
		default:
			// 类型不支持
			wg.Done()
		}
	}
	// 等待所有的数据都发送完成
	wg.Wait()
	// 延时1毫秒 确保数据都处理完成了
	time.Sleep(1 * time.Millisecond)
	close(ch)
	// 数据全部处理完成了
	wgReceiving.Wait()
	return
}
