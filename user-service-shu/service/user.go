package service

import (
	"context"
	"fmt"

	pb "github.com/abduboriykhalid/my_tdm/user-service-shu/genproto"
	l "github.com/abduboriykhalid/my_tdm/user-service-shu/pkg/logger"
	cl "github.com/abduboriykhalid/my_tdm/user-service-shu/service/grpc_client"
	"github.com/abduboriykhalid/my_tdm/user-service-shu/storage"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed while generating uuid for new user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while generating uuid")
	}
	req.Id = id.String()
	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error("failed while inserting user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while inserting user")
	}

	for _, post := range req.Posts {
		post.UserId = req.Id
		createdPosts, err := s.client.PostService().CreatePost(ctx, post)
		if err != nil {
			s.logger.Error("failed while inserting user post", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while inserting user post")
		}
		fmt.Println(createdPosts)
	}

	return user, nil
}
