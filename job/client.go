package job

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"opensnail.com/snail-job/snail-job-go/constant"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/rpc"
)

type SnailJobClient struct {
	opts   *dto.Options
	client rpc.UnaryRequestClient
}

func NewSnailJobClient(opts *dto.Options) SnailJobClient {
	// 创建 gRPC 连接
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", opts.ServerHost, opts.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// todo
	//defer conn.Close()
	return SnailJobClient{
		opts:   opts,
		client: rpc.NewUnaryRequestClient(conn),
	}
}

func GenerateHostId(length int) string {
	result := "go-"
	for i := 0; i < length-3; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

var snailHostID = GenerateHostId(20)

func GenerateReqID() int64 {
	return time.Now().UnixMilli()
}

func (receiver *SnailJobClient) SendToServer(uri string, payload interface{}, jobName string) constant.StatusEnum {

	c := receiver.opts

	// 构建 Metadata 和请求体
	// TODO: 提取
	headers := map[string]string{
		"host-id":      snailHostID,
		"host-ip":      c.HostIP,
		"version":      "1.2.0",
		"host-port":    c.HostPort,
		"namespace":    c.Namespace,
		"group-name":   c.GroupName,
		"token":        c.Token,
		"content-type": "application/json",
	}
	reqId := GenerateReqID()
	body, _ := json.Marshal([]interface{}{payload})
	req := &rpc.GrpcSnailJobRequest{
		ReqId: reqId,
		Metadata: &rpc.Metadata{
			Uri:     uri,
			Headers: headers,
		},
		Body: string(body),
	}

	// 发送请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := receiver.client.UnaryRequest(ctx, req)
	if err != nil {
		log.Printf("%s失败: 无法连接服务器: %v", jobName, err)
		return constant.NO
	}

	// 检查响应
	if response.ReqId != reqId {
		log.Fatalf("reqId 不一致!")
	}

	if response.Status == int32(constant.YES) {
		log.Printf("%s成功: reqId=%d", jobName, reqId)
		data, err := json.Marshal(payload)
		if err == nil {
			log.Printf("data=%s", string(data))
		} else {
			log.Printf("data=%v", payload)
		}
		return constant.YES
	}

	log.Printf("%s失败: %s", jobName, response.Message)
	return constant.NO
}

func (receiver *SnailJobClient) SendBatchLogReport(payload []*dto.JobLogTask) {
	URI := "/batch/server/report/log"
	receiver.SendToServer(URI, payload, "日志批量上报")
}

func (receiver *SnailJobClient) SendDispatchResult(payload interface{}) {
	URI := "/report/dispatch/result"
	receiver.SendToServer(URI, payload, "结果上报")
}

func (receiver *SnailJobClient) SendHeartbeat() {
	for {
		receiver.SendToServer("/beat", "PING", "发送心跳")
		time.Sleep(time.Second * 30)
	}
}
