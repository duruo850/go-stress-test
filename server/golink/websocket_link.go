// Package golink 连接
package golink

import (
	"fmt"
	"go-stress-test/conf"
	"go-stress-test/helper"
	"sync"
	"time"

	"go-stress-test/model"
	"go-stress-test/server/client"
)

// WebSocket webSocket go link
func WebSocket(chanID uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup,
	request *model.Request, ws *client.WebSocket) {
	defer func() {
		wg.Done()
	}()
	defer func() {
		_ = ws.Close()
	}()

	var (
		i uint64
	)
	for {
		select {
		default:
			// 请求
			webSocketRequest(chanID, ch, i, request, ws)
			// 结束条件
			i = i + 1
			if i >= totalNumber {
				goto end
			}
		}
	}
end:

	if request.Keepalive == true {
		// 保持连接
		chWaitFor := make(chan int, 0)
		<-chWaitFor
	}
	return
}

func encodeLength(length int) []byte {
	var encLength []byte
	for {
		digit := byte(length % 128)
		length /= 128
		if length > 0 {
			digit |= 0x80
		}
		encLength = append(encLength, digit)
		if length == 0 {
			break
		}
	}
	return encLength
}

func boolToByte(b bool) byte {
	switch b {
	case true:
		return 1
	default:
		return 0
	}
}

// webSocketRequest 请求
func webSocketRequest(chanID uint64, ch chan<- *model.RequestResults, i uint64, request *model.Request,
	ws *client.WebSocket) {
	var (
		startTime = time.Now()
		isSucceed = false
		errCode   = model.HTTPOk
		msg       []byte
	)
	// 需要发送的数据
	seq := fmt.Sprintf("%d_%d", chanID, i)
	err := ws.Write([]byte(conf.WriteData))

	if err != nil {
		errCode = model.RequestErr // 请求错误
	} else {
		msg, err = ws.Read()
		if err != nil {
			errCode = model.ParseError
			fmt.Println("读取数据 失败~")
		} else {
			errCode, isSucceed = request.GetVerifyWebSocket()(request, seq, msg)
		}
	}
	requestTime := uint64(helper.DiffNano(startTime))
	requestResults := &model.RequestResults{
		Time:      requestTime,
		IsSucceed: isSucceed,
		ErrCode:   errCode,
	}
	requestResults.SetID(chanID, i)
	ch <- requestResults
}
