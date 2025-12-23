package seed

import (
	"database/sql"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/auth"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/jmoiron/sqlx"
)

// Platform tenant ID - fixed UUID for the superadmin's dedicated tenant
const PlatformTenantID = "00000000-0000-0000-0000-000000000000"

// EnsureSuperAdmin creates or updates the superadmin user based on env ADMIN_EMAIL/ADMIN_PASSWORD.
// If ADMIN_PASSWORD is empty, a random password is generated and printed to logs by the caller.
// The superadmin is assigned to a dedicated "platform" tenant, separate from user-facing tenants.
func EnsureSuperAdmin(db *sqlx.DB, cfg config.AppConfig) (generated string, err error) {
	email := strings.TrimSpace(cfg.AdminEmail)
	if email == "" {
		return "", nil
	}

	// 1. Ensure platform tenant exists
	_, err = db.Exec(`
		INSERT INTO tenants (id, slug, name, tenant_type, is_active, enabled_services, app_name)
		VALUES ($1, 'platform', 'Platform Administration', 'university', true, ARRAY['chat'], 'Platform Admin')
		ON CONFLICT (id) DO NOTHING
	`, PlatformTenantID)
	if err != nil {
		return "", err
	}

	pw := cfg.AdminPassword
	if pw == "" {
		pw = auth.GeneratePass()
		generated = pw
	}
	hash, _ := auth.HashPassword(pw)

	// 2. Upsert superadmin user by email
	var id string
	err = db.QueryRowx(`SELECT id FROM users WHERE email=$1`, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// insert new superadmin
			username := strings.Split(email, "@")[0]
			first := "Super"
			last := "Admin"
			err = db.QueryRow(`INSERT INTO users (username,email,first_name,last_name,role,password_hash,is_active,is_superadmin)
				VALUES ($1,$2,$3,$4,'superadmin',$5,true,true) RETURNING id`, username, email, first, last, hash).Scan(&id)
			if err != nil {
				return generated, err
			}
		} else {
			return generated, err
		}
	} else {
		// update password if provided
		if cfg.AdminPassword != "" {
			_, err = db.Exec(`UPDATE users SET password_hash=$1 WHERE id=$2`, hash, id)
			if err != nil {
				return generated, err
			}
		}
	}

	// 3. Ensure superadmin has membership in platform tenant
	_, err = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, is_primary)
		VALUES ($1, $2, 'admin', true)
		ON CONFLICT (user_id, tenant_id) DO NOTHING
	`, id, PlatformTenantID)

	return generated, err
}
