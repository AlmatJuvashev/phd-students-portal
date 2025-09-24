package seed

import (
	"database/sql"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/jmoiron/sqlx"
)

// EnsureSuperAdmin creates or updates the superadmin user based on env ADMIN_EMAIL/ADMIN_PASSWORD.
// If ADMIN_PASSWORD is empty, a random password is generated and printed to logs by the caller.
func EnsureSuperAdmin(db *sqlx.DB, cfg config.AppConfig) (generated string, err error) {
	email := strings.TrimSpace(cfg.AdminEmail)
	if email == "" {
		return "", nil
	}
	pw := cfg.AdminPassword
	if pw == "" {
		pw = auth.GeneratePass()
		generated = pw
	}
	hash, _ := auth.HashPassword(pw)
	// Upsert by email
	var id string
	err = db.QueryRowx(`SELECT id FROM users WHERE email=$1`, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// insert
			username := strings.Split(email, "@")[0]
			first := "Super"
			last := "Admin"
			_, err = db.Exec(`INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active)
				VALUES ($1,$2,$3,$4,'superadmin',$5,true)`, username, email, first, last, hash)
			return generated, err
		}
		return generated, err
	}
	// update password if provided
	if cfg.AdminPassword != "" {
		_, err = db.Exec(`UPDATE users SET password_hash=$1 WHERE id=$2`, hash, id)
	}
	return generated, err
}
