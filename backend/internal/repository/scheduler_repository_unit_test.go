package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newMockSchedulerRepo(t *testing.T) (*SQLSchedulerRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewSQLSchedulerRepository(sqlxDB)

	return repo, mock, func() {
		db.Close()
	}
}



// --- Terms Tests ---

func TestSQLSchedulerRepository_CreateTerm(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	term := &models.AcademicTerm{
		TenantID:  "tenant-1",
		Name:      "Fall 2024",
		Code:      "2024-FALL",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("term-1", time.Now(), time.Now())

		// Regex to match "INSERT INTO academic_terms" and verify arguments
		mock.ExpectQuery(`INSERT INTO academic_terms`).
			WithArgs(term.TenantID, term.Name, term.Code, term.StartDate, term.EndDate, term.IsActive).
			WillReturnRows(rows)

		err := repo.CreateTerm(ctx, term)
		assert.NoError(t, err)
		assert.Equal(t, "term-1", term.ID)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO academic_terms`).
			WillReturnError(fmt.Errorf("db error"))

		err := repo.CreateTerm(ctx, term)
		assert.Error(t, err)
	})
}

func TestSQLSchedulerRepository_GetTerm(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow("term-1", "Fall 2024")
		mock.ExpectQuery(`SELECT \* FROM academic_terms WHERE id = \$1`).
			WithArgs("term-1").
			WillReturnRows(rows)

		term, err := repo.GetTerm(ctx, "term-1")
		assert.NoError(t, err)
		assert.Equal(t, "Fall 2024", term.Name)
	})
}

func TestSQLSchedulerRepository_ListTerms(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow("term-1", "Fall 2024").
			AddRow("term-2", "Spring 2025")

		mock.ExpectQuery(`SELECT \* FROM academic_terms WHERE tenant_id = \$1 ORDER BY start_date DESC`).
			WithArgs("tenant-1").
			WillReturnRows(rows)

		terms, err := repo.ListTerms(ctx, "tenant-1")
		assert.NoError(t, err)
		assert.Len(t, terms, 2)
	})
}

func TestSQLSchedulerRepository_UpdateTerm(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	term := &models.AcademicTerm{
		ID:        "term-1",
		Name:      "Fall 2024 Updated",
		Code:      "2024-FALL",
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE academic_terms`).
			WithArgs(term.Name, term.Code, term.StartDate, term.EndDate, term.IsActive, sqlmock.AnyArg(), term.ID). // AnyArg for updated_at
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateTerm(ctx, term)
		assert.NoError(t, err)
	})
}

func TestSQLSchedulerRepository_DeleteTerm(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM academic_terms WHERE id = \$1`).
			WithArgs("term-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteTerm(ctx, "term-1")
		assert.NoError(t, err)
	})
}

// --- Offerings Tests ---

func TestSQLSchedulerRepository_CreateOffering(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	offering := &models.CourseOffering{
		CourseID:        "course-1",
		TermID:          "term-1",
		TenantID:        "tenant-1",
		Section:         "A",
		DeliveryFormat:  "in_person",
		MaxCapacity:     30,
		VirtualCapacity: toPtr(0),
		MeetingURL:      toPtr(""),
		Status:          "active",
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("offering-1", time.Now(), time.Now())

		mock.ExpectQuery(`INSERT INTO course_offerings`).
			WithArgs(offering.CourseID, offering.TermID, offering.TenantID, offering.Section, offering.DeliveryFormat, offering.MaxCapacity, offering.VirtualCapacity, offering.MeetingURL, offering.Status).
			WillReturnRows(rows)

		err := repo.CreateOffering(ctx, offering)
		assert.NoError(t, err)
		assert.Equal(t, "offering-1", offering.ID)
	})
}

func TestSQLSchedulerRepository_ListOfferingsByInstructor(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "section"}).AddRow("offering-1", "A")

		mock.ExpectQuery(`SELECT co\.\* FROM course_offerings co JOIN course_staff cs ON`).
			WithArgs("inst-1", "term-1").
			WillReturnRows(rows)

		offerings, err := repo.ListOfferingsByInstructor(ctx, "inst-1", "term-1")
		assert.NoError(t, err)
		assert.Len(t, offerings, 1)
	})
}

func TestSQLSchedulerRepository_GetOffering(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "section"}).AddRow("offering-1", "A")
		mock.ExpectQuery(`SELECT \* FROM course_offerings WHERE id = \$1`).
			WithArgs("offering-1").
			WillReturnRows(rows)

		o, err := repo.GetOffering(ctx, "offering-1")
		assert.NoError(t, err)
		assert.Equal(t, "A", o.Section)
	})
}

func TestSQLSchedulerRepository_ListOfferings(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("All Terms", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "section"}).AddRow("offering-1", "A")
		mock.ExpectQuery(`SELECT \* FROM course_offerings WHERE tenant_id = \$1 ORDER BY course_id, section`).
			WithArgs("tenant-1").
			WillReturnRows(rows)

		list, err := repo.ListOfferings(ctx, "tenant-1", "")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("Specific Term", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "section"}).AddRow("offering-1", "A")
		mock.ExpectQuery(`SELECT \* FROM course_offerings WHERE tenant_id = \$1 AND term_id = \$2 ORDER BY course_id, section`).
			WithArgs("tenant-1", "term-1").
			WillReturnRows(rows)

		list, err := repo.ListOfferings(ctx, "tenant-1", "term-1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_UpdateOffering(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	offering := &models.CourseOffering{
		ID:              "offering-1",
		Section:         "B",
		DeliveryFormat:  "online",
		MaxCapacity:     100,
		VirtualCapacity: toPtr(100),
		MeetingURL:      toPtr("http://zoom.us"),
		Status:          "cancelled",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`UPDATE course_offerings`).
			WithArgs(offering.Section, offering.DeliveryFormat, offering.MaxCapacity, offering.VirtualCapacity, offering.MeetingURL, offering.Status, sqlmock.AnyArg(), offering.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateOffering(ctx, offering)
		assert.NoError(t, err)
	})
}

// --- Staff Tests ---

func TestSQLSchedulerRepository_AddStaff(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	staff := &models.CourseStaff{
		CourseOfferingID: "offering-1",
		UserID:           "user-1",
		Role:             "INSTRUCTOR",
		IsPrimary:        true,
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at"}).
			AddRow("staff-1", time.Now())

		mock.ExpectQuery(`INSERT INTO course_staff`).
			WithArgs(staff.CourseOfferingID, staff.UserID, staff.Role, staff.IsPrimary).
			WillReturnRows(rows)

		err := repo.AddStaff(ctx, staff)
		assert.NoError(t, err)
		assert.Equal(t, "staff-1", staff.ID)
	})
}

func TestSQLSchedulerRepository_ListStaff(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow("staff-1", "user-1")
		mock.ExpectQuery(`SELECT \* FROM course_staff WHERE course_offering_id = \$1`).
			WithArgs("offering-1").
			WillReturnRows(rows)

		list, err := repo.ListStaff(ctx, "offering-1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_RemoveStaff(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM course_staff WHERE id = \$1`).
			WithArgs("staff-1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.RemoveStaff(ctx, "staff-1")
		assert.NoError(t, err)
	})
}

// --- Sessions Tests ---

func TestSQLSchedulerRepository_CreateSession(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	// Using nil for optional fields to simplify args
	session := &models.ClassSession{
		CourseOfferingID: "offering-1",
		Title:            "Lecture 1",
		Date:             time.Now(),
		StartTime:        "10:00",
		EndTime:          "12:00",
		Type:             "lecture",
		SessionFormat:    toPtr("in_person"),
		MeetingURL:       nil,
		RoomID:           nil,
		InstructorID:     nil,
	}

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow("session-1", time.Now(), time.Now())

		mock.ExpectQuery(`INSERT INTO class_sessions`).
			WithArgs(session.CourseOfferingID, session.Title, session.Date, session.StartTime, session.EndTime, session.RoomID, session.InstructorID, session.Type, session.SessionFormat, session.MeetingURL).
			WillReturnRows(rows)

		err := repo.CreateSession(ctx, session)
		assert.NoError(t, err)
		assert.Equal(t, "session-1", session.ID)
	})
}

func TestSQLSchedulerRepository_ListSessions(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("s1", "Lec 1")
		mock.ExpectQuery(`SELECT \* FROM class_sessions WHERE course_offering_id = \$1`).
			WithArgs("offering-1", start, end).
			WillReturnRows(rows)

		list, err := repo.ListSessions(ctx, "offering-1", start, end)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_ListSessionsByRoom(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "room_id"}).AddRow("s1", "room-1")
		mock.ExpectQuery(`SELECT \* FROM class_sessions WHERE room_id = \$1`).
			WithArgs("room-1", start, end).
			WillReturnRows(rows)

		list, err := repo.ListSessionsByRoom(ctx, "room-1", start, end)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_ListSessionsForTerm(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("s1", "Lec 1")
		mock.ExpectQuery(`SELECT s\.\* FROM class_sessions s JOIN course_offerings o ON`).
			WithArgs("term-1").
			WillReturnRows(rows)

		list, err := repo.ListSessionsForTerm(ctx, "term-1")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_UpdateSession(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	roomID := "room-2"
	session := &models.ClassSession{
		ID:           "s1",
		RoomID:       &roomID,
		Title:        "Updated Lec",
		Date:         time.Now(),
		StartTime:    "11:00",
		EndTime:      "13:00",
		Type:         "lab",
		IsCancelled:  false,
		InstructorID: nil, // If nil, it won't be in the UPDATE args based on repository logic
	}

	t.Run("Success", func(t *testing.T) {
		// Verify logic inside UpdateSession:
		// setParts includes "room_id" if not nil. "instructor_id" if not nil.
		// Always: title, date, start_time, end_time, type, is_cancelled.
		mock.ExpectExec(`UPDATE class_sessions SET updated_at = \?, room_id = \?, title=\?, date=\?, start_time=\?, end_time=\?, type=\?, is_cancelled=\? WHERE id=\?`).
			WithArgs(sqlmock.AnyArg(), *session.RoomID, session.Title, session.Date, session.StartTime, session.EndTime, session.Type, session.IsCancelled, session.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateSession(ctx, session)
		assert.NoError(t, err)
	})
}

// --- Cohorts Tests ---

func TestSQLSchedulerRepository_AddCohortToOffering(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO course_offering_cohorts").
			WithArgs("offering-1", "cohort-1").
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.AddCohortToOffering(ctx, "offering-1", "cohort-1")
		assert.NoError(t, err)
	})
}

func TestSQLSchedulerRepository_GetOfferingCohorts(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"cohort_id"}).AddRow("c1").AddRow("c2")
		mock.ExpectQuery(`SELECT cohort_id FROM course_offering_cohorts WHERE course_offering_id = \$1`).
			WithArgs("offering-1").
			WillReturnRows(rows)
		cohorts, err := repo.GetOfferingCohorts(ctx, "offering-1")
		assert.NoError(t, err)
		assert.Len(t, cohorts, 2)
	})
}

func TestSQLSchedulerRepository_ListSessionsForCohorts(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()
	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)
	
	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("s1", "Lec 1")
		mock.ExpectQuery(`SELECT s\.\* FROM class_sessions s`).
			WithArgs("c1", "c2", start, end).
			WillReturnRows(rows)
		
		sessions, err := repo.ListSessionsForCohorts(ctx, []string{"c1", "c2"}, start, end)
		assert.NoError(t, err)
		assert.Len(t, sessions, 1)
	})
}

func TestSQLSchedulerRepository_ListSessionsByInstructor(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title"}).AddRow("s1", "Lec 1")
		mock.ExpectQuery(`SELECT \* FROM class_sessions WHERE instructor_id = \$1`).
			WithArgs("inst-1", start, end).
			WillReturnRows(rows)

		list, err := repo.ListSessionsByInstructor(ctx, "inst-1", start, end)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestSQLSchedulerRepository_DeleteSession(t *testing.T) {
	repo, mock, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM class_sessions WHERE id = \$1`).
			WithArgs("s1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteSession(ctx, "s1")
		assert.NoError(t, err)
	})
}

func TestSQLSchedulerRepository_ListSessionsForCohorts_Empty(t *testing.T) {
	repo, _, teardown := newMockSchedulerRepo(t)
	defer teardown()

	ctx := context.Background()
	start := time.Now()
	end := start.Add(24 * time.Hour)

	t.Run("Empty Cohorts", func(t *testing.T) {
		list, err := repo.ListSessionsForCohorts(ctx, []string{}, start, end)
		assert.NoError(t, err)
		assert.Empty(t, list)
	})
}
