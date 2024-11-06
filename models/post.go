// models/post.go
package models

import "time"

// Post representa a estrutura de um post
type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"` // ID do usuário que criou o post
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
