package job

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/rpc"
	"opensnail.com/snail-job/snail-job-go/util"
)

type Server struct {
	rpc.UnimplementedUnaryRequestServer
	endpoint *Dispatcher
	logger   Logger
}

// UnaryRequest implements snailjob.UnaryRequestServer
func (s *Server) UnaryRequest(_ context.Context, in *rpc.GrpcSnailJobRequest) (*rpc.GrpcResult, error) {
	s.logger.Debug("Received: %v", in)
	var request []dto.DispatchJobRequest
	util.ToObj([]byte(in.Body), &request)
	result := s.endpoint.DispatchJob(request[0])
	return &rpc.GrpcResult{ReqId: in.ReqId, Status: result.Status, Message: result.Message, Data: string(util.ToByteArr(result.Data))}, nil
}

func RunServer(opts *dto.Options, client SnailJobClient, executors map[string]NewJobExecutor, factory LoggerFactory) {

	logger := factory.GetLogger("grpc-server", nil)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", opts.HostPort))
	if err != nil {
		logger.Info("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterUnaryRequestServer(s, &Server{endpoint: Init(client, executors, factory), logger: logger})
	logger.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve: %v", err)
	}
}
