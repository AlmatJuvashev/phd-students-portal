package seed

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
)

type rawContact struct {
	Name  interface{}       `json:"name"`
	Title interface{}       `json:"title"`
	Email string            `json:"email"`
	Phone string            `json:"phone"`
	Extra map[string]string `json:"-"`
}

func normalizeLocalized(v interface{}) map[string]string {
	out := map[string]string{}
	switch val := v.(type) {
	case string:
		if val != "" {
			out["ru"] = val
		}
	case map[string]interface{}:
		for k, vv := range val {
			if s, ok := vv.(string); ok && s != "" {
				out[k] = s
			}
		}
	case map[string]string:
		for k, s := range val {
			if s != "" {
				out[k] = s
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// EnsureContacts seeds contacts from the bundled JSON file if the table is empty.
func EnsureContacts(db *sqlx.DB) error {
	var exists int
	if err := db.Get(&exists, "SELECT COUNT(*) FROM contacts"); err != nil {
		if strings.Contains(err.Error(), "contacts") {
			// Table might not exist yet; skip seeding.
			return nil
		}
		return err
	}
	if exists > 0 {
		return nil
	}

	paths := []string{
		filepath.Join("frontend", "src", "playbooks", "supervisors_contacts.json"),
		filepath.Join("..", "frontend", "src", "playbooks", "supervisors_contacts.json"),
		filepath.Join("internal", "seed", "contacts.json"),
	}

	var data []byte
	for _, p := range paths {
		if b, err := os.ReadFile(p); err == nil {
			data = b
			break
		}
	}
	if data == nil {
		return errors.New("contacts seed file not found")
	}

	var raw []rawContact
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("parse contacts seed: %w", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	order := 1
	for _, rc := range raw {
		name := normalizeLocalized(rc.Name)
		if name == nil {
			continue
		}
		title := normalizeLocalized(rc.Title)
		if _, err := tx.Exec(
			`INSERT INTO contacts (name, title, email, phone, sort_order) VALUES ($1,$2,$3,$4,$5)`,
			toJSON(name), toJSON(title), nullableString(rc.Email), nullableString(rc.Phone), order,
		); err != nil {
			return fmt.Errorf("insert contact seed failed: %w", err)
		}
		order++
	}

	return tx.Commit()
}

func nullableString(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func toJSON(m map[string]string) interface{} {
	if m == nil || len(m) == 0 {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
}
