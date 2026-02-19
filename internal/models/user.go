package models

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
