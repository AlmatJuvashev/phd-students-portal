package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Contact struct {
	ID        string       `db:"id" json:"id"`
	TenantID  string       `db:"tenant_id" json:"tenant_id"`
	Name      LocalizedMap `db:"name" json:"name"`
	Title     LocalizedMap `db:"title" json:"title,omitempty"`
	Email     *string      `db:"email" json:"email,omitempty"`
	Phone     *string      `db:"phone" json:"phone,omitempty"`
	SortOrder int          `db:"sort_order" json:"sort_order"`
	IsActive  bool         `db:"is_active" json:"is_active"`
	CreatedAt string       `db:"created_at" json:"created_at"`
	UpdatedAt string       `db:"updated_at" json:"updated_at"`
}

type LocalizedMap map[string]string

func (m *LocalizedMap) Scan(value interface{}) error {
	if value == nil {
		*m = LocalizedMap{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
	if len(b) == 0 {
		*m = LocalizedMap{}
		return nil
	}
	return json.Unmarshal(b, m)
}

func (m LocalizedMap) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
