package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	snailjob "opensnail.com/snail-job/snail-job-go"
	"strconv"
	"time"

	"math/rand"

	pb "opensnail.com/snail-job/snail-job-go/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StatusEnum 是响应状态的枚举
type StatusEnum int

const (
	NO  StatusEnum = 0
	YES StatusEnum = 1
)

var (
	addr = flag.String("addr", "localhost:17888", "the address to connect to")
	port = flag.Int("port", 17889, "The server port")
)

// SnailJobRequest 定义 SnailJobRequest 和 Metadata 的数据结构
type SnailJobRequest struct {
	ReqId int64
}

func GenerateHostId(length int) string {
	result := "go-"
	for i := 0; i < length-3; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}

func GenerateReqID() int64 {
	return time.Now().UnixMilli()
}

var snailHostID = GenerateHostId(20)

// BuildRequest 定义构建 SnailJobRequest 的方法
func BuildRequest(args []interface{}) SnailJobRequest {
	return SnailJobRequest{ReqId: GenerateReqID()}
}

// SendToServer 发送请求到程服务器
func SendToServer(uri string, payload interface{}, jobName string) StatusEnum {
	// 构建请求
	request := BuildRequest([]interface{}{payload})

	// 创建 gRPC 连接
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 初始化 gRPC 客户端
	client := pb.NewUnaryRequestClient(conn) // 替换为实际的生成的 gRPC 客户端

	// 构建 Metadata 和请求体
	// TODO: 提取
	headers := map[string]string{
		"host-id":      snailHostID,
		"host-ip":      "localhost",
		"version":      "1.2.0",
		"host-port":    "17889",
		"namespace":    "764d604ec6fc45f68cd92514c40e9e1a",
		"group-name":   "snail_job_demo_group",
		"token":        "SJ_Wyz3dmsdbDOkDujOTSSoBjGQP1BMsVnj",
		"content-type": "application/json",
	}
	body, _ := json.Marshal([]interface{}{payload})
	req := &pb.GrpcSnailJobRequest{
		ReqId: request.ReqId,
		Metadata: &pb.Metadata{
			Uri:     uri,
			Headers: headers,
		},
		Body: string(body),
	}

	// 发送请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := client.UnaryRequest(ctx, req)
	if err != nil {
		log.Printf("%s失败: 无法连接服务器: %v", jobName, err)
		return NO
	}

	// 检查响应
	if response.ReqId != request.ReqId {
		log.Fatalf("reqId 不一致!")
	}

	if response.Status == int32(YES) {
		log.Printf("%s成功: reqId=%d", jobName, request.ReqId)
		data, err := json.Marshal(payload)
		if err == nil {
			log.Printf("data=%s", string(data))
		} else {
			log.Printf("data=%v", payload)
		}
		return YES
	}

	log.Printf("%s失败: %s", jobName, response.Message)
	return NO
}

func SendHeartbeat() {
	for {
		SendToServer("/beat", []string{"PING"}, "发送心跳")
		time.Sleep(time.Second * 30)
	}

}

type server struct {
	pb.UnimplementedUnaryRequestServer
}

// UnaryRequest implements snailjob.UnaryRequestServer
func (s *server) UnaryRequest(_ context.Context, in *pb.GrpcSnailJobRequest) (*pb.GrpcResult, error) {
	log.Printf("Received: %v", in)
	return &pb.GrpcResult{ReqId: in.ReqId, Status: 1, Message: "", Data: "true"}, nil
}

// todo 未完成
type GrpcServer struct {
	config snailjob.SysConfig
}

func (receiver GrpcServer) name() {

}

func (receiver GrpcServer) SendHeartbeat() {
	for {
		SendToServer("/beat", []string{"PING"}, "发送心跳")
		time.Sleep(time.Second * 30)
	}
}

//
//func main() {
//	flag.Parse()
//	go SendHeartbeat()
//	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	s := grpc.NewServer()
//	pb.RegisterUnaryRequestServer(s, &server{})
//	log.Printf("server listening at %v", lis.Addr())
//	if err := s.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
