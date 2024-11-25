package job

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

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
	s.logger.Debugf("Received: %v", in)
	metadata := in.Metadata
	var result dto.Result
	if metadata.Uri == "/job/dispatch/v1" {
		var request []dto.DispatchJobRequest
		err := util.ToObj([]byte(in.Body), &request)
		if err != nil {
			return nil, err
		}
		result = s.endpoint.DispatchJob(request[0])
	} else if metadata.Uri == "/job/stop/v1" {
		var request []dto.StopJob
		err := util.ToObj([]byte(in.Body), &request)
		if err != nil {
			return nil, err

		}
		result = s.endpoint.Stop(request[0])
	} else {
		return &rpc.GrpcResult{ReqId: in.ReqId, Status: 0, Message: "uri is not supports, uri=" + metadata.Uri, Data: ""}, errors.New("uri is not supports. uri=" + metadata.Uri)
	}

	arr, _ := util.ToByteArr(result.Data)
	return &rpc.GrpcResult{ReqId: in.ReqId, Status: result.Status, Message: result.Message, Data: string(arr)}, nil
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
