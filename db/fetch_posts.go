package db

import (
	"context"
	"time"
)


const fetchPosts = `
	SELECT p1.id, p1.content, p1.user_id, p1.original_post_id, p1.created_at
	FROM posts p1
	LEFT JOIN posts p2
	ON p1.original_post_id = p2.id
	WHERE p1.user_id = $1
`

type FetchPostsParams struct {
	UserID int64
	IsOnlyMyPosts bool
	Size int
	Page int
	StartDate *time.Time
	EndDate *time.Time
}

func (q *Queries) FetchPosts(ctx context.Context, arg FetchPostsParams) ([]Post, error){
	rows, err := q.db.QueryContext(ctx, fetchPosts, arg.UserID)
	if err != nil {
		return []Post{}, err
	}
	posts := make([]Post, 0)
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID,&p.Content,&p.OriginalPostID, &p.CreatedAt)
		posts = append(posts, p)
	}
	return posts, nil	
}