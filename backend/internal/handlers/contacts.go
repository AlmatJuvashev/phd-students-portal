package handlers

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ContactsHandler struct {
	db *sqlx.DB
}

type Contact struct {
	ID        string       `db:"id" json:"id"`
	TenantID  string       `db:"tenant_id" json:"tenant_id"` // Added for multitenancy
	Name      LocalizedMap `db:"name" json:"name"`
	Title     LocalizedMap `db:"title" json:"title,omitempty"`
	Email     *string      `db:"email" json:"email,omitempty"`
	Phone     *string      `db:"phone" json:"phone,omitempty"`
	SortOrder int          `db:"sort_order" json:"sort_order"`
	IsActive  bool         `db:"is_active" json:"is_active"`
	CreatedAt string       `db:"created_at" json:"created_at"`
	UpdatedAt string       `db:"updated_at" json:"updated_at"`
}

type contactPayload struct {
	Name      map[string]string `json:"name"`
	Title     map[string]string `json:"title"`
	Email     *string           `json:"email"`
	Phone     *string           `json:"phone"`
	SortOrder *int              `json:"sort_order"`
	IsActive  *bool             `json:"is_active"`
}

func NewContactsHandler(db *sqlx.DB) *ContactsHandler {
	return &ContactsHandler{db: db}
}

func (h *ContactsHandler) PublicList(c *gin.Context) {
	var contacts []Contact
	err := h.db.Select(
		&contacts,
		`SELECT id, name, title, email, phone, sort_order, is_active,
		        to_char(created_at,'YYYY-MM-DD\"T\"HH24:MI:SSZ') AS created_at,
		        to_char(updated_at,'YYYY-MM-DD\"T\"HH24:MI:SSZ') AS updated_at
		   FROM contacts
		  WHERE is_active = true
		  ORDER BY sort_order, created_at`,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if contacts == nil {
		contacts = []Contact{}
	}
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactsHandler) AdminList(c *gin.Context) {
	includeInactive := c.Query("all") == "true"
	query := `SELECT id, name, title, email, phone, sort_order, is_active,
		        to_char(created_at,'YYYY-MM-DD\"T\"HH24:MI:SSZ') AS created_at,
		        to_char(updated_at,'YYYY-MM-DD\"T\"HH24:MI:SSZ') AS updated_at
		   FROM contacts`
	if !includeInactive {
		query += ` WHERE is_active = true`
	}
	query += ` ORDER BY sort_order, created_at`

	var contacts []Contact
	if err := h.db.Select(&contacts, query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if contacts == nil {
		contacts = []Contact{}
	}
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactsHandler) Create(c *gin.Context) {
	var req contactPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	tenantID := c.GetString("tenant_id") // Get tenant from context
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant context required"})
		return
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var id string
	err := h.db.QueryRow(
		`INSERT INTO contacts (tenant_id, name, title, email, phone, sort_order, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		tenantID,
		toJSON(req.Name),
		toJSON(req.Title),
		contactNullablePtr(req.Email),
		contactNullablePtr(req.Phone),
		sortOrder,
		isActive,
	).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "contacts_name_key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Contact already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ContactsHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req contactPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setParts := []string{"updated_at = now()"}
	args := []interface{}{}

	if len(req.Name) > 0 {
		setParts = append(setParts, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, toJSON(req.Name))
	}
	if req.Title != nil {
		setParts = append(setParts, "title = $"+strconv.Itoa(len(args)+1))
		args = append(args, toJSON(req.Title))
	}
	if req.Email != nil {
		setParts = append(setParts, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, contactNullablePtr(req.Email))
	}
	if req.Phone != nil {
		setParts = append(setParts, "phone = $"+strconv.Itoa(len(args)+1))
		args = append(args, contactNullablePtr(req.Phone))
	}
	if req.SortOrder != nil {
		setParts = append(setParts, "sort_order = $"+strconv.Itoa(len(args)+1))
		args = append(args, *req.SortOrder)
	}
	if req.IsActive != nil {
		setParts = append(setParts, "is_active = $"+strconv.Itoa(len(args)+1))
		args = append(args, *req.IsActive)
	}

	if len(setParts) == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	args = append(args, id)
	query := "UPDATE contacts SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(len(args))

	if _, err := h.db.Exec(query, args...); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *ContactsHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := h.db.Exec(`UPDATE contacts SET is_active = false, updated_at = now() WHERE id = $1`, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

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

func toJSON(m map[string]string) interface{} {
	if len(m) == 0 {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return b
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
