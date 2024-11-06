// post.go
package post

import (
	"database/sql"
	"edsb/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Cria a tabela de posts com referência ao usuário
func CreatePostsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		likes_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	log.Println("Tabela posts criada com sucesso (se não existia).")
	return nil
}

// Handler para obter todos os posts
func GetPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, user_id, title, content, created_at FROM posts")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// Handler para obter um post específico
func GetPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var post models.Post
		query := "SELECT id, user_id, title, content, created_at FROM posts WHERE id = $1"
		if err := db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

// Handler para criar um novo post
func CreatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, created_at"
		err := db.QueryRow(query, post.UserID, post.Title, post.Content).Scan(&post.ID, &post.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

// Handler para atualizar um post existente
func UpdatePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "UPDATE posts SET title = $1, content = $2 WHERE id = $3"
		if _, err := db.Exec(query, post.Title, post.Content, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler para deletar um post
func DeletePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		query := "DELETE FROM posts WHERE id = $1"
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
