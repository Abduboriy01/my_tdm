package postgres

import (
	"github.com/jmoiron/sqlx"

	pb "github.com/abduboriykhalid/my_tdm/user-service-shu/genproto"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		rUser = pb.User{}
	)

	err := r.db.QueryRow("INSERT INTO users (id, first_name, last_name) VALUES($1, $2, $3) RETURNING id, first_name, last_name", user.Id, user.FirstName, user.LastName).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.LastName,
	)
	if err != nil {
		return &pb.User{}, err
	}

	return &rUser, nil
}
