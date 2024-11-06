// models/user.go
package models

import "time"

// User representa um usuário do sistema
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // omitido em respostas JSON // Utilizado somente para registro, pois não pode ser retornado em GET
	CreatedAt time.Time `json:"created_at"`
}
