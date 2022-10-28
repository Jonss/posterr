package db

import (
	"context"
	"fmt"
	"time"
)

const fetchPosts = `
	SELECT
		p1.id, p1.content, p1.user_id,
		p1.original_post_id,
		p1.created_at,
		u1.username,
		p2.id, p2.content, p2.user_id,
		p2.original_post_id,
		p2.created_at,
		u2.username
	FROM posts p1
	LEFT JOIN posts p2
	ON p1.original_post_id = p2.id
	LEFT JOIN users u1
	ON p1.user_id = u1.id
	LEFT JOIN users u2
	ON p2.user_id = u2.id
`

type FetchPostsParams struct {
	UserID int64
	IsOnlyMyPosts bool
	Size int
	Page int
	StartDate *time.Time
	EndDate *time.Time
}

type DetailedPost struct {
	Post
	Username string
}

type FetchPost struct {
	Post DetailedPost
	OriginalPost *DetailedPost
}

type FetchPosts struct {
	Posts []FetchPost
	HasNext bool
	HasPrev bool
}

func (q *Queries) FetchPosts(ctx context.Context, arg FetchPostsParams) (FetchPosts, error) {
	query := fetchPosts
	count := 1

	values := make([]interface{}, 0)

	if arg.IsOnlyMyPosts {
		query = fmt.Sprintf("%s %s p1.user_id = $%d",query, andOrWhere(len(values)), count)
		count++
		values = append(values, arg.UserID)
	}

	if arg.EndDate != nil {
		query = fmt.Sprintf("%s %s p1.created_at <= $%d", query, andOrWhere(len(values)), count)
		values = append(values, arg.EndDate)
		count++
	}

	if arg.StartDate != nil {
		query = fmt.Sprintf("%s %s p1.created_at >= $%d", query, andOrWhere(len(values)), count)
		values = append(values, arg.StartDate)
		count++
	}

	query = fmt.Sprintf("%s ORDER BY p1.id DESC LIMIT %d OFFSET %d", query, arg.Size, arg.Page)

	rows, err := q.db.QueryContext(ctx, query, values...)
	if err != nil {
		return FetchPosts{}, err
	}
	posts := make([]FetchPost, 0)
	for rows.Next() {
		var p1 DetailedPost
		var p2 DetailedPost

		rows.Scan(
			&p1.ID,
			&p1.Content,
			&p1.OriginalPostID,
			&p1.CreatedAt,
			&p1.Username,
			&p2.ID,
			&p2.Content,
			&p2.OriginalPostID,
			&p2.CreatedAt,
			&p2.Username,
		)

		f := FetchPost{Post: p1, OriginalPost: &p2 }
		posts = append(posts, f)
	}

	fetchPosts := FetchPosts{
		Posts: posts,
		HasNext: len(posts) > 0 && len(posts) == arg.Size,
		HasPrev: arg.Page > 0,
	}

	return fetchPosts, nil
}

func andOrWhere(length int) string {
	if length > 0 {
		return "AND"
	}
	return "WHERE"
}