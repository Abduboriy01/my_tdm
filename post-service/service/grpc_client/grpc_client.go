package grpcClient

import (
	"fmt"

	"github.com/my_tdm/post-service/config"
	pb "github.com/my_tdm/post-service/genproto"

	"google.golang.org/grpc"
)

//GrpcClientI ...
type GrpcClientI interface {
	UserService() pb.UserServiceClient
}

//GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	userService pb.UserServiceClient
}

//New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port: %d",
			cfg.UserServiceHost, cfg.UserServicePort)
	}

	return &GrpcClient{
		cfg:         cfg,
		userService: pb.NewUserServiceClient(connUser),
	}, nil
}

func (s *GrpcClient) UserService() pb.UserServiceClient {
	return s.userService
}
