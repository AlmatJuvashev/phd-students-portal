package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSQLResourceRepository_RoomAttributes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLResourceRepository(sqlxDB)
	ctx := context.Background()

	t.Run("SetRoomAttribute_Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO room_attributes").
			WithArgs("room-1", "Projector", "true").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SetRoomAttribute(ctx, &models.RoomAttribute{RoomID: "room-1", Key: "Projector", Value: "true"})
		assert.NoError(t, err)
	})

	t.Run("SetRoomAttribute_Error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO room_attributes").
			WithArgs("room-1", "Projector", "true").
			WillReturnError(fmt.Errorf("db error"))

		err := repo.SetRoomAttribute(ctx, &models.RoomAttribute{RoomID: "room-1", Key: "Projector", Value: "true"})
		assert.Error(t, err)
	})

	t.Run("GetRoomAttributes_Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"room_id", "key", "value"}).
			AddRow("room-1", "Projector", "true").
			AddRow("room-1", "Seats", "50")

		mock.ExpectQuery(`SELECT (.+) FROM room_attributes WHERE room_id=\$1`).
			WithArgs("room-1").
			WillReturnRows(rows)

		attrs, err := repo.GetRoomAttributes(ctx, "room-1")
		assert.NoError(t, err)
		assert.Len(t, attrs, 2)
		assert.Equal(t, "Projector", attrs[0].Key)
		assert.Equal(t, "true", attrs[0].Value)
	})
}
