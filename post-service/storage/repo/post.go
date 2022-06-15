package repo

import (
	pb "github.com/abduboriykhalid/my_tdm/post-service/genproto"
)

//PostStorageI ...
type PostStorageI interface {
	CreatePost(*pb.Post) (*pb.Post, error)
	GetUserPosts(userId string) ([]*pb.Post, error)
}
