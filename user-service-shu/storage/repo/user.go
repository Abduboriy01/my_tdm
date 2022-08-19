package repo

import (
	pb "github.com/my_tdm/user-service-shu/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	GetUserById(userID string) (*pb.User, error)
	GetUserList(limit, page int64) ([]*pb.User, int64, error)
	Update(*pb.User) (*pb.User, error)
	Delete(id string) error
	CheckUniquess(filed, value string) (bool, error)
	LoginUser(*pb.LoginRequest) (*pb.LoginResponse, error)
	CheckField(field, value string) (bool, error)
	RegisterUser(*pb.CreateUserAuthReqBody) (*pb.CreateUserAuthResBody, error)
}
