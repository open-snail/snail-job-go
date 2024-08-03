package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	SNAIL_SERVER_HOST = "127.0.0.1"
	SNAIL_SERVER_PORT = "1788"
	SNAIL_HOST_IP     = "127.0.0.1"
	SNAIL_HOST_PORT   = "1789"
	SNAIL_NAMESPACE   = "764d604ec6fc45f68cd92514c40e9e1a"
	SNAIL_GROUP_NAME  = "snail_job_demo_group"

	SNAIL_LOG_LOCAL_FILENAME     = "snail_job.log"
	SNAIL_LOG_REMOTE_BUFFER_SIZE = 10
	SNAIL_LOG_REMOTE_INTERVAL    = 10
)

var HEADERS = map[string]string{
	"Content-Type": "application/json",
	"host-id":      GenerateHostId(20),
	"host-ip":      SNAIL_HOST_IP,
	"version":      "1.0.0",
	"host-port":    SNAIL_HOST_PORT,
	"namespace":    SNAIL_NAMESPACE,
	"group-name":   SNAIL_GROUP_NAME,
	"token":        "SJ_Wyz3dmsdbDOkDujOTSSoBjGQP1BMsVnj",
}

func GenerateHostId(length int) string {
	result := "go-"
	for i := 0; i < length-3; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

func SendToServer(uri string, payload interface{}, jobName string) {
	request := SnailJobRequest{
		ReqID: GenerateReqID(),
		Args:  []interface{}{payload},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		LocalLog.Printf("Failed to marshal request: %v", err)
		return
	}

	url := fmt.Sprintf("http://%s:%s/%s", SNAIL_SERVER_HOST, SNAIL_SERVER_PORT, uri)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	for key, value := range HEADERS {
		req.Header.Set(key, value)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		LocalLog.Printf("%s失败: %v", jobName, err)
		return
	}

	var serverResponse NettyResult
	if err := json.NewDecoder(resp.Body).Decode(&serverResponse); err != nil {
		LocalLog.Printf("Failed to decode server response: %v", err)
		return
	}

	if request.ReqID != serverResponse.ReqID {
		LocalLog.Println("reqId 不一致的!")
		return
	}

	if serverResponse.Status == STATUS_SUCCESS {
		LocalLog.Printf("%s成功: reqId=%d", jobName, request.ReqID)
	} else {
		LocalLog.Printf("%s失败: %s", jobName, serverResponse.Data)
	}
}

func SendDispatchResult(payload interface{}) {
	URI := "report/dispatch/result"
	SendToServer(URI, payload, "结果上报")
}

func SendBatchLogReport(payload []*JobLogTask) {
	URI := "batch/server/report/log"
	SendToServer(URI, payload, "日志批量上报")
}

func SendHeartbeat() {
	URI := "beat"
	for {
		request := SnailJobRequest{
			ReqID: GenerateReqID(),
		}

		requestBody, err := json.Marshal(request)
		if err != nil {
			LocalLog.Printf("Failed to marshal request: %v", err)
			continue
		}

		url := fmt.Sprintf("http://%s:%s/%s", SNAIL_SERVER_HOST, SNAIL_SERVER_PORT, URI)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		for key, value := range HEADERS {
			req.Header.Set(key, value)
		}

		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			LocalLog.Printf("Failed to send heartbeat: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		var serverResponse NettyResult
		if err := json.NewDecoder(resp.Body).Decode(&serverResponse); err != nil {
			LocalLog.Printf("Failed to decode server response: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		if request.ReqID != serverResponse.ReqID {
			LocalLog.Println("reqId 不一致的!")
			time.Sleep(30 * time.Second)
			continue
		}

		if serverResponse.Status == STATUS_SUCCESS {
			LocalLog.Printf("发送心跳成功: reqId=%d", request.ReqID)
		} else {
			LocalLog.Printf("发送心跳失败: %s", serverResponse.Data)
		}

		time.Sleep(30 * time.Second)
	}
}

func GenerateReqID() int64 {
	return time.Now().UnixMilli()
}
