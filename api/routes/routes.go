// routes.go
package routes

import (
	"database/sql"
	"edsb/api/comment"
	"edsb/api/like"
	"edsb/api/post"
	"edsb/api/user"

	"github.com/gorilla/mux"
)

// ConfigureRoutes define todas as rotas da aplicação
func ConfigureRoutes(r *mux.Router, db *sql.DB) {

	// Rotas para usuários
	r.HandleFunc("/users", user.GetUsers(db)).Methods("GET")
	r.HandleFunc("/users/{id}", user.GetUser(db)).Methods("GET")
	r.HandleFunc("/users/register", user.CreateUser(db)).Methods("POST")
	r.HandleFunc("/users/login", user.LoginUser(db)).Methods("POST")
	r.HandleFunc("/users/{id}", user.UpdateUser(db)).Methods("PUT")
	r.HandleFunc("/users/{id}", user.DeleteUser(db)).Methods("DELETE")

	// Rotas para posts
	r.HandleFunc("/posts", post.GetPosts(db)).Methods("GET")
	r.HandleFunc("/posts/{id}", post.GetPost(db)).Methods("GET")
	r.HandleFunc("/posts", post.CreatePost(db)).Methods("POST")
	r.HandleFunc("/posts/{id}", post.UpdatePost(db)).Methods("PUT")
	r.HandleFunc("/posts/{id}", post.DeletePost(db)).Methods("DELETE")

	// Rotas para comentários
	r.HandleFunc("/comments", comment.GetComments(db)).Methods("GET")
	r.HandleFunc("/comments/{id}", comment.GetComment(db)).Methods("GET")
	r.HandleFunc("/comments", comment.CreateComment(db)).Methods("POST")
	r.HandleFunc("/comments/{id}", comment.UpdateComment(db)).Methods("PUT")
	r.HandleFunc("/comments/{id}", comment.DeleteComment(db)).Methods("DELETE")

	// Rotas para likes em posts e comentários
	r.HandleFunc("/posts/{id}/like", like.AddLikeToPost(db)).Methods("POST")
	r.HandleFunc("/posts/{id}/likes/count", like.CountLikesForPost(db)).Methods("GET")
	r.HandleFunc("/comments/{id}/like", like.AddLikeToComment(db)).Methods("POST")
	r.HandleFunc("/comments/{id}/likes/count", like.CountLikesForComment(db)).Methods("GET")

}
