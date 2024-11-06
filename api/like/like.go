// like.go - Versão Refatorada com Remoção de Likes
package like

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	//"edsb/models"
)

// Cria a tabela de likes para posts e comentários
func CreateLikesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS likes (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id) ON DELETE CASCADE,
		post_id INT REFERENCES posts(id) ON DELETE CASCADE,
		comment_id INT REFERENCES comments(id) ON DELETE CASCADE,
		UNIQUE(user_id, post_id),
		UNIQUE(user_id, comment_id)
	);`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	log.Println("Tabela likes criada com sucesso (se não existia).")
	return nil
}

// Handler para adicionar um like a um post
func AddLikeToPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Erro ao processar o formulário"}`, http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil || userID == 0 {
			http.Error(w, `{"error": "ID do usuário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || postID == 0 {
			http.Error(w, `{"error": "ID do post é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		query := "INSERT INTO likes (user_id, post_id) VALUES ($1, $2)"
		_, err = db.Exec(query, userID, postID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao adicionar like"}`, http.StatusInternalServerError)
			return
		}

		// Incrementar a contagem de likes do post
		updateQuery := "UPDATE posts SET likes = likes + 1 WHERE id = $1"
		_, err = db.Exec(updateQuery, postID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao incrementar o like no post"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Like adicionado com sucesso"})
	}
}

// Handler para adicionar um like a um comentário
func AddLikeToComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Erro ao processar o formulário"}`, http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil || userID == 0 {
			http.Error(w, `{"error": "ID do usuário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil || commentID == 0 {
			http.Error(w, `{"error": "ID do comentário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		query := "INSERT INTO likes (user_id, comment_id) VALUES ($1, $2)"
		_, err = db.Exec(query, userID, commentID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao adicionar like"}`, http.StatusInternalServerError)
			return
		}

		// Incrementar a contagem de likes do comentário
		updateQuery := "UPDATE comments SET likes = likes + 1 WHERE id = $1"
		_, err = db.Exec(updateQuery, commentID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao incrementar o like no comentário"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Like adicionado com sucesso"})
	}
}

// Handler para remover um like de um post
func RemoveLikeFromPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Erro ao processar o formulário"}`, http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil || userID == 0 {
			http.Error(w, `{"error": "ID do usuário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || postID == 0 {
			http.Error(w, `{"error": "ID do post é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		query := "DELETE FROM likes WHERE user_id = $1 AND post_id = $2"
		_, err = db.Exec(query, userID, postID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao remover like"}`, http.StatusInternalServerError)
			return
		}

		// Decrementar a contagem de likes do post
		updateQuery := "UPDATE posts SET likes = likes - 1 WHERE id = $1"
		_, err = db.Exec(updateQuery, postID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao decrementar o like no post"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Like removido com sucesso"})
	}
}

// Handler para remover um like de um comentário
func RemoveLikeFromComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			http.Error(w, `{"error": "Erro ao processar o formulário"}`, http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(r.FormValue("user_id"))
		if err != nil || userID == 0 {
			http.Error(w, `{"error": "ID do usuário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		commentID, err := strconv.Atoi(r.FormValue("comment_id"))
		if err != nil || commentID == 0 {
			http.Error(w, `{"error": "ID do comentário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		query := "DELETE FROM likes WHERE user_id = $1 AND comment_id = $2"
		_, err = db.Exec(query, userID, commentID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao remover like"}`, http.StatusInternalServerError)
			return
		}

		// Decrementar a contagem de likes do comentário
		updateQuery := "UPDATE comments SET likes = likes - 1 WHERE id = $1"
		_, err = db.Exec(updateQuery, commentID)
		if err != nil {
			http.Error(w, `{"error": "Erro ao decrementar o like no comentário"}`, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Like removido com sucesso"})
	}
}

// Handler para contar likes de um post
func CountLikesForPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		postID, err := strconv.Atoi(r.URL.Query().Get("post_id"))
		if err != nil || postID == 0 {
			http.Error(w, `{"error": "ID do post é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		var count int
		query := `SELECT likes FROM posts WHERE id = $1`
		err = db.QueryRow(query, postID).Scan(&count)
		if err != nil {
			http.Error(w, `{"error": "Erro ao contar likes"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]int{"likes": count})
	}
}

// Handler para contar likes de um comentário
func CountLikesForComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		commentID, err := strconv.Atoi(r.URL.Query().Get("comment_id"))
		if err != nil || commentID == 0 {
			http.Error(w, `{"error": "ID do comentário é obrigatório e deve ser um número"}`, http.StatusBadRequest)
			return
		}

		var count int
		query := `SELECT likes FROM comments WHERE id = $1`
		err = db.QueryRow(query, commentID).Scan(&count)
		if err != nil {
			http.Error(w, `{"error": "Erro ao contar likes"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]int{"likes": count})
	}
}
