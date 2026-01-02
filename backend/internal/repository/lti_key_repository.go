package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
)
func (r *SQLLTIRepository) CreateKey(ctx context.Context, k models.LTIKey) error {
	query := `INSERT INTO lti_keys (kid, private_key, public_key, algorithm, use) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, k.KID, k.PrivateKey, k.PublicKey, k.Algorithm, k.Use)
	return err
}

func (r *SQLLTIRepository) GetActiveKey(ctx context.Context) (*models.LTIKey, error) {
	// Get the most recent active key
	query := `SELECT id, kid, private_key, public_key, algorithm, use, is_active, created_at, expires_at FROM lti_keys WHERE is_active=true ORDER BY created_at DESC LIMIT 1`
	return scanKey(r.db.QueryRowContext(ctx, query))
}

func (r *SQLLTIRepository) ListActiveKeys(ctx context.Context) ([]models.LTIKey, error) {
	query := `SELECT id, kid, private_key, public_key, algorithm, use, is_active, created_at, expires_at FROM lti_keys WHERE is_active=true ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.LTIKey
	for rows.Next() {
		k, err := scanKeyRow(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, *k)
	}
	return keys, nil
}

func scanKey(row *sql.Row) (*models.LTIKey, error) {
	var k models.LTIKey
	err := row.Scan(
		&k.ID, &k.KID, &k.PrivateKey, &k.PublicKey, &k.Algorithm, &k.Use, &k.IsActive, &k.CreatedAt, &k.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &k, nil
}

func scanKeyRow(rows *sql.Rows) (*models.LTIKey, error) {
	var k models.LTIKey
	err := rows.Scan(
		&k.ID, &k.KID, &k.PrivateKey, &k.PublicKey, &k.Algorithm, &k.Use, &k.IsActive, &k.CreatedAt, &k.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return &k, nil
}
