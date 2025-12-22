package repository

import (
	"context"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type SearchRepository interface {
	SearchUsers(ctx context.Context, query string, role string, userID string, limit int) ([]models.SearchResult, error)
	SearchDocuments(ctx context.Context, query string, role string, userID string, limit int) ([]models.SearchResult, error)
}

type SQLSearchRepository struct {
	db *sqlx.DB
}

func NewSQLSearchRepository(db *sqlx.DB) *SQLSearchRepository {
	return &SQLSearchRepository{db: db}
}

func (r *SQLSearchRepository) SearchUsers(ctx context.Context, query string, role string, userID string, limit int) ([]models.SearchResult, error) {
	// Only admins/staff can search users
	if role != "admin" && role != "superadmin" && role != "chair" && role != "advisor" {
		return []models.SearchResult{}, nil
	}

	sqlQuery := `
		SELECT id, first_name, last_name, email, role 
		FROM users 
		WHERE (first_name ILIKE $1 OR last_name ILIKE $1 OR email ILIKE $1)
	`
	args := []any{"%" + query + "%"}

	if role == "advisor" {
		sqlQuery += ` AND (role != 'student' OR id IN (SELECT student_id FROM student_advisors WHERE advisor_id = $2))`
		args = append(args, userID)
	}

	sqlQuery += ` LIMIT $` + fmt.Sprint(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var u struct {
			ID        string `db:"id"`
			FirstName string `db:"first_name"`
			LastName  string `db:"last_name"`
			Email     string `db:"email"`
			Role      string `db:"role"`
		}
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role); err != nil {
			continue
		}

		link := fmt.Sprintf("/admin/students/%s", u.ID)
		if u.Role != "student" {
			link = "/admin/users"
		}

		results = append(results, models.SearchResult{
			Type:        "student",
			ID:          u.ID,
			Title:       fmt.Sprintf("%s %s", u.FirstName, u.LastName),
			Subtitle:    u.Email,
			Description: u.Role,
			Link:        link,
		})
	}
	return results, nil
}

func (r *SQLSearchRepository) SearchDocuments(ctx context.Context, query string, role string, userID string, limit int) ([]models.SearchResult, error) {
	sqlQuery := `
		SELECT a.id, a.filename, u.first_name, u.last_name, ni.node_id
		FROM node_instance_slot_attachments a
		JOIN node_instance_slots s ON a.slot_id = s.id
		JOIN node_instances ni ON s.node_instance_id = ni.id
		JOIN users u ON ni.user_id = u.id
		WHERE a.filename ILIKE $1 AND a.is_active = true
	`
	args := []any{"%" + query + "%"}

	if role == "student" {
		sqlQuery += ` AND ni.user_id = $2`
		args = append(args, userID)
	} else if role == "advisor" {
		sqlQuery += ` AND ni.user_id IN (SELECT student_id FROM student_advisors WHERE advisor_id = $2)`
		args = append(args, userID)
	}
	// Admin/Chair see all

	sqlQuery += ` LIMIT $` + fmt.Sprint(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var d struct {
			ID        string `db:"id"`
			Filename  string `db:"filename"`
			FirstName string `db:"first_name"`
			LastName  string `db:"last_name"`
			NodeID    string `db:"node_id"`
		}
		if err := rows.Scan(&d.ID, &d.Filename, &d.FirstName, &d.LastName, &d.NodeID); err != nil {
			continue
		}

		results = append(results, models.SearchResult{
			Type:        "document",
			ID:          d.ID,
			Title:       d.Filename,
			Subtitle:    fmt.Sprintf("Owner: %s %s", d.FirstName, d.LastName),
			Description: fmt.Sprintf("Node: %s", d.NodeID),
			Link:        "#", // Placeholder
		})
	}
	return results, nil
}
