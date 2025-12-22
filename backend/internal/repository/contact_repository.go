package repository

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ContactRepository interface {
	ListPublic(ctx context.Context, tenantID string) ([]models.Contact, error)
	ListAdmin(ctx context.Context, tenantID string, includeInactive bool) ([]models.Contact, error)
	Create(ctx context.Context, tenantID string, contact models.Contact) (string, error)
	Update(ctx context.Context, tenantID string, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, tenantID string, id string) error
}

type SQLContactRepository struct {
	db *sqlx.DB
}

func NewSQLContactRepository(db *sqlx.DB) *SQLContactRepository {
	return &SQLContactRepository{db: db}
}

func (r *SQLContactRepository) ListPublic(ctx context.Context, tenantID string) ([]models.Contact, error) {
	var contacts []models.Contact
	err := r.db.SelectContext(ctx,
		&contacts,
		`SELECT id, tenant_id, name, title, email, phone, sort_order, is_active,
		        to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') AS created_at,
		        to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') AS updated_at
		   FROM contacts
		  WHERE tenant_id = $1 AND is_active = true
		  ORDER BY sort_order, created_at`,
		tenantID,
	)
	if err != nil {
		return nil, err
	}
	if contacts == nil {
		contacts = []models.Contact{}
	}
	return contacts, nil
}

func (r *SQLContactRepository) ListAdmin(ctx context.Context, tenantID string, includeInactive bool) ([]models.Contact, error) {
	query := `SELECT id, tenant_id, name, title, email, phone, sort_order, is_active,
		        to_char(created_at,'YYYY-MM-DD"T"HH24:MI:SSZ') AS created_at,
		        to_char(updated_at,'YYYY-MM-DD"T"HH24:MI:SSZ') AS updated_at
		   FROM contacts
		  WHERE tenant_id = $1`
	if !includeInactive {
		query += ` AND is_active = true`
	}
	query += ` ORDER BY sort_order, created_at`

	var contacts []models.Contact
	if err := r.db.SelectContext(ctx, &contacts, query, tenantID); err != nil {
		return nil, err
	}
	if contacts == nil {
		contacts = []models.Contact{}
	}
	return contacts, nil
}

func (r *SQLContactRepository) Create(ctx context.Context, tenantID string, contact models.Contact) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO contacts (tenant_id, name, title, email, phone, sort_order, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		tenantID,
		contact.Name, // LocalizedMap handles Value()
		contact.Title,
		contactNullablePtr(contact.Email),
		contactNullablePtr(contact.Phone),
		contact.SortOrder,
		contact.IsActive,
	).Scan(&id)
	return id, err
}

func (r *SQLContactRepository) Update(ctx context.Context, tenantID string, id string, updates map[string]interface{}) error {
	setParts := []string{"updated_at = now()"}
	args := []interface{}{}
	
	// Sort keys for deterministic order
	keys := make([]string, 0, len(updates))
	for k := range updates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := updates[k]
		// Handle custom nullable logic for email/phone if passed as nil?
		// Value should be prepared by Service layer.
		// If Service passes `nil` for email/phone, it means set to NULL?
		// Or Service passes `contactNullablePtr` result?
		// Repo shouldn't assume too much logic. It should take values.
		// But `contactNullablePtr` logic is specific.
		
		// Let's assume input map values are ready for DB driver (e.g. sql.NullString or nil).
		// EXCEPT for email/phone special empty string handling.
		
		if k == "email" || k == "phone" {
			// Check if value is string ptr or string?
			// If pointer to string, apply helper?
			if vPtr, ok := val.(*string); ok {
				val = contactNullablePtr(vPtr)
			} else if vStr, ok := val.(string); ok {
				val = contactNullableString(vStr)
			}
		}
		
		setParts = append(setParts, fmt.Sprintf("%s = $%d", k, len(args)+1))
		args = append(args, val)
	}
	
	if len(setParts) == 1 {
		return nil // Nothing to update
	}

	args = append(args, id, tenantID)
	query := "UPDATE contacts SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(len(args)-1) + " AND tenant_id = $" + strconv.Itoa(len(args))

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *SQLContactRepository) Delete(ctx context.Context, tenantID string, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE contacts SET is_active = false, updated_at = now() WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	return err
}

// Helpers
func contactNullableString(v string) interface{} {
	if strings.TrimSpace(v) == "" {
		return nil
	}
	return v
}

func contactNullablePtr(v *string) interface{} {
	if v == nil {
		return nil
	}
	return contactNullableString(*v)
}
