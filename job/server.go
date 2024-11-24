package job

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"

	"google.golang.org/grpc"
	"opensnail.com/snail-job/snail-job-go/dto"
	"opensnail.com/snail-job/snail-job-go/rpc"
	"opensnail.com/snail-job/snail-job-go/util"
)

type Server struct {
	rpc.UnimplementedUnaryRequestServer
	endpoint *Dispatcher
	logger   *logrus.Entry
}

// UnaryRequest implements snailjob.UnaryRequestServer
func (s *Server) UnaryRequest(_ context.Context, in *rpc.GrpcSnailJobRequest) (*rpc.GrpcResult, error) {
	s.logger.Debug("Received: %v", in)
	metadata := in.Metadata
	var result dto.Result
	if metadata.Uri == "/job/dispatch/v1" {
		var request []dto.DispatchJobRequest
		util.ToObj([]byte(in.Body), &request)
		result = s.endpoint.DispatchJob(request[0])
	} else if metadata.Uri == "/job/stop/v1" {
		var request []dto.StopJob
		util.ToObj([]byte(in.Body), &request)
		result = s.endpoint.Stop(request[0])
	} else {
		return &rpc.GrpcResult{ReqId: in.ReqId, Status: 0, Message: "uri is not supports, uri=" + metadata.Uri, Data: ""}, errors.New("uri is not supports. uri=" + metadata.Uri)
	}

	return &rpc.GrpcResult{ReqId: in.ReqId, Status: result.Status, Message: result.Message, Data: string(util.ToByteArr(result.Data))}, nil
}

func RunServer(opts *dto.Options, client SnailJobClient, executors map[string]NewJobExecutor, factory LoggerFactory) {

	logger := factory.GetLocalLogger("grpc-server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", opts.HostPort))
	if err != nil {
		logger.Infof("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	rpc.RegisterUnaryRequestServer(s, &Server{endpoint: Init(client, executors, factory), logger: logger})
	logger.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Errorf("failed to serve: %v", err)
	}
}
