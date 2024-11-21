package job

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/rpc"
	"opensnail.com/snail-job/snail-job-go/util"
)

type Server struct {
	rpc.UnimplementedUnaryRequestServer
	endpoint *Dispatcher
}

// UnaryRequest implements snailjob.UnaryRequestServer
func (s *Server) UnaryRequest(_ context.Context, in *rpc.GrpcSnailJobRequest) (*rpc.GrpcResult, error) {
	log.Printf("Received: %v", in)
	var request []dto.DispatchJobRequest
	util.ToObj([]byte(in.Body), &request)
	result := s.endpoint.DispatchJob(request[0])
	return &rpc.GrpcResult{ReqId: in.ReqId, Status: result.Status, Message: result.Message, Data: string(util.ToByteArr(result.Data))}, nil
}

func RunServer(opts *dto.Options, client SnailJobClient, executors map[string]IJobExecutor) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", opts.HostPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterUnaryRequestServer(s, &Server{endpoint: Init(client, executors)})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
