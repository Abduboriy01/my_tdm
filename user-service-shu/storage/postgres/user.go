package postgres

import (
	//	"os/user"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	pb "github.com/my_tdm/user-service-shu/genproto"
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
	err := r.db.QueryRow("INSERT INTO users (id, first_name, last_name, email, password, phone_number) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, first_name, last_name, email, password, phone_number", user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.LastName,
		&rUser.Email,
		&rUser.Password,
		&rUser.PhoneNumber,
	)
	if err != nil {
		return &pb.User{}, err
	}

	return &rUser, nil
}

func (r *userRepo) Update(user *pb.User) (*pb.User, error) {
	updateUser := pb.User{}
	query := `UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3 RETURNING id, first_name, last_name`
	err := r.db.QueryRow(query, user.FirstName, user.LastName, user.Id).Scan(
		&updateUser.Id,
		&updateUser.FirstName,
		&updateUser.LastName,
	)
	if err != nil {
		return nil, err
	}

	return &updateUser, nil
}

func (r *userRepo) GetUserById(userID string) (*pb.User, error) {
	query := `SELECT id, first_name, last_name FROM users WHERE id = $1`
	var ID, firstName, lastName string
	err := r.db.QueryRow(query, userID).Scan(
		&ID,
		&firstName,
		&lastName,
	)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:        ID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (r *userRepo) GetUserList(limit, page int64) ([]*pb.User, int64, error) {
	var (
		users []*pb.User
		count int64
	)
	offset := (page - 1) * limit

	query := `SELECT id, first_name, last_name FROM users order by first_name OFFSET $1 LIMIT $2`

	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	countQuery := `SELECT count(*) FROM users`
	err = r.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (r *userRepo) GetFilteredUsers(ID string) (*pb.User, error) {
	var user pb.User
	sql := sqlbuilder.Select("id", "first_name", "last_name").From("users").
		Where("id=$1").String()
	err := r.db.QueryRow(sql, ID).Scan(&user.Id, &user.FirstName, &user.LastName)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *userRepo) Delete(Id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) CheckUniquess(filed, value string) (bool, error) {
	var exists int64
	query := `SELECT count(1) FROM users WHERE $1 = $2`
	err := r.db.QueryRow(query, filed, value).Scan(
		&exists,
	)

	if err != nil {
		return false, err
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil
}

func (r *userRepo) LoginUser(req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var rUser = pb.LoginResponse{}

	err := r.db.QueryRow(`SELECT id,first_name,username,email,password,phone_number from users WHERE email=$1`, req.Email).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.Username,
		&rUser.Email,
		&rUser.Password,
		&rUser.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(rUser.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}
	return &rUser, err
}

func (r *userRepo) ListUser(req *pb.ListUserReq) (*pb.ListUserResponse, error) {
	rUser := pb.GetAllUser{}

	offset := (req.Page - 1) * req.Limit

	rows, err := r.db.Query("select id,first_name,last_name from users  OFFSET $1 LIMIT $2", offset, req.Limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var data pb.User
		err = rows.Scan(&data.Id, &data.FirstName, &data.LastName)
		if err != nil {
			return nil, err
		}
		rUser.Users = append(rUser.Users, &data)
	}
	count := 0
	countQuery := `SELECT count(*)from users`
	err = r.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		return nil, err
	}
	return &pb.ListUserResponse{
		Users: rUser.Users,
		Count: int64(count),
	}, nil
}

func (r *userRepo) CheckField(field, value string) (bool, error) {
	var existClient int64
	if field == "username" {
		row := r.db.QueryRow(`SELECT count(1) FROM users where username=$1`, value)

		if err := row.Scan(&existClient); err != nil {
			return false, err
		}
	} else if field == "email" {
		row := r.db.QueryRow(`SELECT count(1) FROM users where email=$1`, value)
		if err := row.Scan(&existClient); err != nil {
			return false, err
		}
	} else {
		return false, nil
	}
	if existClient == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *userRepo) RegisterUser(user *pb.CreateUserAuthReqBody) (*pb.CreateUserAuthResBody, error) {
	var rUser = pb.CreateUserAuthResBody{}
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow("INSERT INTO users (id,first_name,username,email,password,phone_number) values($1,$2,$3,$4,$5,$6)RETURNING id,first_name,username,email,password,phone_number", id, user.FirstName, user.Username, user.Email, user.Password, user.PhoneNumber).Scan(
		&rUser.Id,
		&rUser.FirstName,
		&rUser.Username,
		&rUser.Email,
		&rUser.Password,
		&rUser.PhoneNumber,
	)
	if err != nil {
		return &pb.CreateUserAuthResBody{}, err
	}
	return &rUser, err
}
