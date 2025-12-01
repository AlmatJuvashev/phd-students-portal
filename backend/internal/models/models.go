package models

import "time"

type Role string

const (
	RoleSuperAdmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleStudent    Role = "student"
	RoleAdvisor    Role = "advisor"
	RoleChair      Role = "chair"
)

type User struct {
	ID           string    `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Role         Role      `db:"role" json:"role"`
	PasswordHash string    `db:"password_hash" json:"-"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	Phone        string    `db:"phone" json:"phone"`
	Program      string    `db:"program" json:"program"`
	Specialty    string    `db:"specialty" json:"specialty"`
	Department   string    `db:"department" json:"department"`
	Cohort       string    `db:"cohort" json:"cohort"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
