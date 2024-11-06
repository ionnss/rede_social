// main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"edsb/api/comment"
	"edsb/api/like"
	"edsb/api/post"
	"edsb/api/routes"
	"edsb/api/user"
	"edsb/views"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Função de conexão ao banco de dados
func connectDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	return sql.Open("postgres", connStr)
}

func main() {
	// Inicia conexão com banco de dados
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Cria as tabelas usando as funções nos arquivos de API
	if err := user.CreateUsersTable(db); err != nil {
		log.Fatalf("Erro ao criar tabela users: %v", err)
	}
	if err := post.CreatePostsTable(db); err != nil {
		log.Fatalf("Erro ao criar tabela posts: %v", err)
	}
	if err := comment.CreateCommentsTable(db); err != nil {
		log.Fatalf("Erro ao criar tabela comments: %v", err)
	}
	if err := like.CreateLikesTable(db); err != nil {
		log.Fatalf("Erro ao criar tabela likes: %v", err)
	}
	log.Println("Banco de dados inicializado com sucesso.")

	// Configura o roteador
	r := mux.NewRouter()

	// Configura as rotas da API
	routes.ConfigureRoutes(r, db)

	// Configura as rotas para os templates
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.RenderTemplate(w, "index.html", nil)
	}).Methods("GET")

	r.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		views.RenderTemplate(w, "register.html", nil)
	}).Methods("GET")

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		views.RenderTemplate(w, "login.html", nil)
	}).Methods("GET")

	// Inicia o servidor
	log.Println("Servidor rodando na porta :8081 em http://localhost:8081")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
