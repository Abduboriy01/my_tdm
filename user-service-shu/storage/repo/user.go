package repo

import (
	pb "github.com/abduboriykhalid/my_tdm/user-service-shu/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
}
