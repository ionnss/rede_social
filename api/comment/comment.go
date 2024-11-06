// comment.go
package comment

import (
	"database/sql"
	"edsb/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Cria a tabela de comentários
func CreateCommentsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		post_id INT REFERENCES posts(id) ON DELETE CASCADE,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		likes_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	log.Println("Tabela comments criada com sucesso (se não existia).")
	return nil
}

// Handler para obter todos os comentários
func GetComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []models.Comment
		for rows.Next() {
			var comment models.Comment
			if err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

// Handler para obter um comentário específico
func GetComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var comment models.Comment
		query := "SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1"
		if err := db.QueryRow(query, id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)
	}
}

// Handler para criar um novo comentário
func CreateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at"
		err := db.QueryRow(query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}
}

// Handler para atualizar um comentário existente
func UpdateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "UPDATE comments SET content = $1 WHERE id = $2"
		if _, err := db.Exec(query, comment.Content, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Handler para deletar um comentário
func DeleteComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		query := "DELETE FROM comments WHERE id = $1"
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
