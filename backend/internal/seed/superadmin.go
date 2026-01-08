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
		VALUES ($1, 'superadmin', 'Superadmin Tenant', 'university', true, ARRAY['chat'], 'Superadmin Universe')
		ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, slug=EXCLUDED.slug
	`, PlatformTenantID)
	if err != nil {
		return "", err
	}

	// 2. Strip is_superadmin from all users first to ensure only one exists
	_, err = db.Exec(`UPDATE users SET is_superadmin = false WHERE is_superadmin = true`)
	if err != nil {
		return "", err
	}

	// 3. Remove existing superadmin memberships in other tenants
	_, err = db.Exec(`DELETE FROM user_tenant_memberships WHERE role = 'superadmin' OR 'superadmin' = ANY(roles)`)
	if err != nil {
		return "", err
	}

	pw := cfg.AdminPassword
	if pw == "" {
		pw = auth.GeneratePass()
		generated = pw
	}
	hash, _ := auth.HashPassword(pw)

	// 4. Upsert superadmin user by email
	var id string
	err = db.QueryRowx(`SELECT id FROM users WHERE email=$1`, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// insert new superadmin
			username := strings.Split(email, "@")[0]
			first := "System"
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
		// update status and password
		_, err = db.Exec(`UPDATE users SET is_superadmin=true, role='superadmin', password_hash=$1 WHERE id=$2`, hash, id)
		if err != nil {
			return generated, err
		}
	}

	// 5. Ensure superadmin has membership in platform tenant
	_, err = db.Exec(`
		INSERT INTO user_tenant_memberships (user_id, tenant_id, role, roles, is_primary)
		VALUES ($1, $2, 'superadmin', ARRAY['superadmin']::text[], true)
		ON CONFLICT (user_id, tenant_id) DO UPDATE SET role='superadmin', roles=ARRAY['superadmin']::text[]
	`, id, PlatformTenantID)

	return generated, err
}
