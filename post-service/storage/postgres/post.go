package postgres

import (
	// "strings"
	// "testing/quick"

	pb "github.com/abduboriykhalid/my_tdm/post-service/genproto"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
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

	err := r.db.QueryRow("INSERT INTO posts (id, name, description, user_id) VALUES($1, $2, $3, $4) RETURNING id, name, description, user_id", post.Id, post.Name, post.Description, post.UserId).Scan(
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
		if err != nil {
			return nil, err
		}
		_, err = r.db.Exec("INSERT INTO post_medias (id, type, link, post_id) VALUES($1, $2, $3, $4)", id, media.Type, media.Link, post.Id)
		if err != nil {
			return &pb.Post{}, err
		}
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

		mQuery := `SELECT id, link, type, FROM post_medias WHERE post_id = $1`
		rows, err := r.db.Query(mQuery, post.Id)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var media pb.Media
			err = rows.Scan(
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
