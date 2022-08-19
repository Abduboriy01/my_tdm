package service

import (
	"context"

	pb "github.com/my_tdm/post-service/genproto"
	l "github.com/my_tdm/post-service/pkg/logger"
	"github.com/my_tdm/post-service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
}

//NewUserService ...
func NewPostService(db *sqlx.DB, log l.Logger) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	// id, err := uuid.NewV4()
	// if err != nil {
	// 	s.logger.Error("failed while generating uuid for new post", l.Error(err))
	// 	return nil, status.Error(codes.Internal, "failed while generating uuid")
	// }
	// req.Id = id.String()

	post, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error("failed while inserting post", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while inserting post")
	}

	return post, nil
}

func (s *PostService) GetUserPosts(ctx context.Context, req *pb.GetByUserIdRequest) (*pb.GetUserPostsResponse, error) {
	posts, err := s.storage.Post().GetUserPosts(req.UserId)
	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}

	return &pb.GetUserPostsResponse{
		Posts: posts,
	}, nil
}
