// user.go
package user

import (
	"database/sql"
	"edsb/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// Cria tabela de usuário
func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(100) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	log.Println("Tabela users criada com sucesso (se não existia).")
	return nil
}

// Registro de um novo usuario no banco de dados
func RegisterUser(db *sql.DB, username, email, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, username, email, passwordHash)
	if err != nil {
		return err
	}
	log.Println("Usuário registrado com sucesso.")
	return nil
}

// Autenticação do usuário com base no email e senha fornecidos pelo mesmo
func AuthenticateUser(db *sql.DB, email, password string) (bool, error) {
	var storedHash string

	query := `SELECT password_hash FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Usuário não encontrado
		}
		return false, err // Erro ao buscar usuário
	}

	// Verifica se a senha fornecida corresponde ao hash armazenado
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false, nil // Senha incorreta
	}
	return true, nil // Autenticação bem-sucedida
}

// Handler para autenticar um usuário
func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extrai e valida os valores dos comapos do forumlário
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Error(w, `{"error": "Todos os campos são obrigatórios"}`, http.StatusBadRequest)
		}

		isAuthenticated, err := AuthenticateUser(db, email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isAuthenticated {
			http.Error(w, `{"error": "Usuário ou senha inválidos"}`, http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login bem-sucedido"))
		//json.NewEncoder(w).Encode(map[string]string{"message": "Login feito com sucesso"})
	}
}

// Handler para obter todos os usuários
func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT id, username, email, created_at FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		json.NewEncoder(w).Encode(users)
	}
}

// Handler para obter um usuário específico
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]

		var user models.User
		query := "SELECT id, username, email, created_at FROM users WHERE id = $1"
		if err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// Handler para criar um novo usuário
func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Extrai e valida os valores dos campos do formulário
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if username == "" || email == "" || password == "" {
			http.Error(w, `{"error": "Todos os campos são obrigatórios"}`, http.StatusBadRequest)
			return
		}

		// Registra o usuário no banco de dados
		if err := RegisterUser(db, username, email, password); err != nil {
			http.Error(w, `{"error": "Erro ao registrar usuário"}`, http.StatusInternalServerError)
			return
		}

		// Responde com uma mensagem de sucesso em JSON
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Usuário registrado com sucesso"))
		//json.NewEncoder(w).Encode(map[string]string{"message": "Usuário registrado com sucesso"})
	}
}

// Handler para atualizar um usuário existente
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "UPDATE users SET username = $1, email = $2 WHERE id = $3"
		if _, err := db.Exec(query, user.Username, user.Email, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuário atualizado com sucesso"})
	}
}

// Handler para deletar um usuário
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := mux.Vars(r)["id"]
		query := "DELETE FROM users WHERE id = $1"
		if _, err := db.Exec(query, id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuário deletado com sucesso"})
	}
}
