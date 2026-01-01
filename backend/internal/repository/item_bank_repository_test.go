package repository

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestSQLItemBankRepository_CreateBank(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLItemBankRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		bank := &models.QuestionBank{
			TenantID:    "t-1",
			Name:        "Bank 1",
			Description: "Desc",
			IsActive:    true,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("b-1", time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO question_banks").
			WithArgs(bank.TenantID, bank.Name, bank.Description, bank.IsActive).
			WillReturnRows(rows)

		err := repo.CreateBank(ctx, bank)
		assert.NoError(t, err)
		assert.Equal(t, "b-1", bank.ID)
	})
}

func TestSQLItemBankRepository_CreateItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLItemBankRepository(sqlxDB)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		item := &models.QuestionItem{
			BankID:     "b-1",
			Type:       "essay",
			Content:    types.JSONText("{}"),
			Difficulty: 3,
			Tags:       pq.StringArray{"tag1"},
			IsActive:   true,
		}

		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("i-1", time.Now(), time.Now())

		// Note: pq.StringArray formatting in mock can be tricky, typically specific driver handling.
		// sqlmock matches args by value.
		mock.ExpectQuery("INSERT INTO question_items").
			WithArgs(item.BankID, item.Type, item.Content, item.Difficulty, item.Tags, item.IsActive).
			WillReturnRows(rows)

		err := repo.CreateItem(ctx, item)
		assert.NoError(t, err)
		assert.Equal(t, "i-1", item.ID)
	})
}
