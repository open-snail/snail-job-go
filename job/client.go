package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/open-snail/snail-job-go/constant"
	"github.com/open-snail/snail-job-go/dto"
	"github.com/open-snail/snail-job-go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SnailJobClient struct {
	opts   *dto.Options
	client rpc.UnaryRequestClient
	log    *logrus.Entry
}

func NewSnailJobClient(opts *dto.Options, factory LoggerFactory) SnailJobClient {
	// 创建 gRPC 连接
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", opts.ServerHost, opts.ServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// todo
	//defer conn.Close()
	return SnailJobClient{
		opts:   opts,
		client: rpc.NewUnaryRequestClient(conn),
		log:    factory.GetLocalLogger("grpc-client"),
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
	l := receiver.log
	c := receiver.opts

	// 构建 Metadata 和请求体
	headers := map[string]string{
		"host-id":      snailHostID,
		"host-ip":      c.HostIP,
		"version":      constant.VERSION,
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
		l.Errorf("%s失败: 无法连接服务器: %s", jobName, err)
		return constant.NO
	}

	// 检查响应
	if response.ReqId != reqId {
		l.Errorf("reqId 不一致!")
	}

	if response.Status == int32(constant.YES) {
		l.Debugf("%s成功: reqId=%d", jobName, reqId)
		data, err := json.Marshal(payload)
		if err == nil {
			l.Debugf("data=%s", string(data))
		} else {
			l.Errorf("data=%v", payload)
		}
		return constant.YES
	}

	l.Errorf("%s失败: %s", jobName, response.Message)
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

func (receiver *SnailJobClient) SendBatchReportMapTask(req dto.MapTaskRequest) constant.StatusEnum {
	return receiver.SendToServer("/batch/report/job/map/task/v1", req, "请求分片")
}
