package postgres

import (


	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/my_tdm/post-service/genproto"
)

type postRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *pb.Post) (*pb.Post, error) {
	var (
		rPost = pb.Post{}
	)

	postid, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow("INSERT INTO posts (id, name, description, user_id) VALUES($1, $2, $3, $4) RETURNING id, name, description, user_id", postid, post.Name, post.Description, post.UserId).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId,
	)
	if err != nil {
		return &pb.Post{}, err

	}

	for _, media := range post.Medias {
		id, err := uuid.NewV4()
		var mediaa pb.Media
		if err != nil {
			return nil, err
		}
		err = r.db.QueryRow("INSERT INTO post_medias (id, type, link, post_id) VALUES($1, $2, $3, $4) RETURNING id,type,link,post_id", id, media.Type, media.Link, postid).Scan(
			&mediaa.Id,
			&mediaa.Type,
			&mediaa.Link,
			&mediaa.PostId,
		)
		if err != nil {
			return &pb.Post{}, err
		}
		rPost.Medias = append(rPost.Medias, &mediaa)
	}


	return &rPost, nil
}

func (r *postRepo) GetUserPosts(userID string) ([]*pb.Post, error) {

	var (
		posts []*pb.Post
	)
	query := `SELECT id, name, description FROM posts WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post pb.Post

		err = rows.Scan(
			&post.Id,
			&post.Name,
			&post.Description,
		)
		if err != nil {
			return nil, err
		}

		mQuery := `SELECT id, link, type FROM post_medias WHERE post_id = $1`
		row, err := r.db.Query(mQuery, post.Id)
		if err != nil {
			return nil, err
		}

		for row.Next() {
			var media pb.Media
			err = row.Scan(
				&media.Id,
				&media.Link,
				&media.Type,
			)
			if err != nil {
				return nil, err
			}
			post.Medias = append(post.Medias, &media)
		}

		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *postRepo) CreatePostUser(user *pb.User) error {

	err := r.db.QueryRow("INSERT INTO post_users (id, first_name, last_name) VALUES($1, $2, $3) ", user.Id, user.FirstName, user.LastName).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return err
	}

	return nil
}
