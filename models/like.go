// models/like.go
package models

// Like representa a estrutura de um like
type Like struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	PostID    int `json:"post_id,omitempty"`
	CommentID int `json:"comment_id,omitempty"`
}
