package service

import (
	"context"
	"fmt"

	pb "github.com/my_tdm/user-service-shu/genproto"
	l "github.com/my_tdm/user-service-shu/pkg/logger"
	messagebroker "github.com/my_tdm/user-service-shu/pkg/messagebroker"
	cl "github.com/my_tdm/user-service-shu/service/grpc_client"
	"github.com/my_tdm/user-service-shu/storage"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage   storage.IStorage
	logger    l.Logger
	client    cl.GrpcClientI
	publisher map[string]messagebroker.Publisher
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI, publisher map[string]messagebroker.Publisher) *UserService {
	return &UserService{
		storage:   storage.NewStoragePg(db),
		logger:    log,
		client:    client,
		publisher: publisher,
	}
}

func (s *UserService) publishUserMessage(rawUser pb.User) error {
	data, err := rawUser.Marshal()
	if err != nil {
		return err
	}

	logProduct := rawUser.String()

	err = s.publisher["user"].Publish([]byte("user"), data, logProduct)
	if err != nil {
		return err
	}

	return nil
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

	err = s.publishUserMessage(*user)
	if err != nil {
		s.logger.Error("failed while publishing user info", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while publishing user info")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, req *pb.User) (*pb.User, error) {
	updatedUser, err := s.storage.User().Update(req)
	if err != nil {
		s.logger.Error("failed while updating user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while updating user")
	}

	return updatedUser, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error("failed while getting user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user")
	}

	posts, err := s.client.PostService().GetUserPosts(ctx, &pb.GetByUserIdRequest{UserId: req.Id})
	if err != nil {
		s.logger.Error("failed while getting user posts", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user posts")
	}
	user.Posts = posts.Posts

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, req *pb.ListUserReq) (*pb.ListUserResponse, error) {
	users, count, err := s.storage.User().GetUserList(req.Limit, req.Page)
	if err != nil {
		s.logger.Error("failed while getting all users", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting all users")
	}
	for _, user := range users {

		post, err := s.client.PostService().GetUserPosts(ctx, &pb.GetByUserIdRequest{UserId: user.Id})
		if err != nil {
			s.logger.Error("failed while getting users posts", l.Error(err))
			return nil, status.Error(codes.Internal, "failed while getting users posts")
		}

		user.Posts = post.Posts

	}
	return &pb.ListUserResponse{
		Users: users,
		Count: count,
	}, nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	err := s.storage.User().Delete(req.Id)
	if err != nil {
		s.logger.Error("failed delete user by id", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete user by id")
	}
	return &pb.Empty{}, nil
}

func (s *UserService) CheckUniquess(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqResp, error) {
	exists, err := s.storage.User().CheckUniquess(req.Filed, req.Value)

	if err != nil {
		s.logger.Error("failed check uniquess of user data", l.Error(err))
		return nil, status.Error(codes.Internal, "failed to check user uniquess data")
	}

	return &pb.CheckUniqResp{
		IsExist: exists,
	}, nil
}

func (s UserService) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.storage.User().LoginUser(req)
	if err != nil {
		s.logger.Error("failed while logging in user ", l.Error(err))
		return nil, status.Error(codes.Internal, "your password is wrong,please check and retype")
	}
	return user, nil
}

func (s *UserService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	check, err := s.storage.User().CheckField(req.Field, req.Value)
	if err != nil {
		s.logger.Error("failed while getting user", l.Error(err))
		return nil, status.Error(codes.Internal, "failed while getting user")

	}
	return &pb.CheckFieldResponse{
		Check: check,
	}, err
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.CreateUserAuthReqBody) (*pb.CreateUserAuthResBody, error) {
	user, err := s.storage.User().RegisterUser(req)
	if err != nil {
		s.logger.Error("failed while register user", l.Error(err))
		return nil, err
	}

	return user, nil
}

// func (s UserService) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse,error){
// 	user,err:=s.storage.User().LoginUser(req)
// 	if err != nil {
// 		s.logger.Error("failed while logging in user ", l.Error(err))
// 		return nil, status.Error(codes.Internal, "your password is wrong,please check and retype")
// 	}
// 	return user,nil
// }
