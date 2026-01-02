package repository

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLCourseContentRepository_Modules(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()
	repo := NewSQLCourseContentRepository(db)

	ctx := context.Background()

	// 1. Create Tenant
	tID := "00000000-0000-0000-0000-000000000001"
	courseID := "00000000-0000-0000-0000-000000000001"
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type) VALUES ($1, $2, 'T Docs', 'university') ON CONFLICT DO NOTHING`, tID, tID)
	require.NoError(t, err)

	// 2. Create Program
	_, err = db.Exec(`INSERT INTO programs (id, code, name, tenant_id, title) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'P1', 'Prog', $1, '{"en": "Prog"}') ON CONFLICT DO NOTHING`, tID)
	require.NoError(t, err)

	// 3. Create Course
	_, err = db.Exec(`INSERT INTO courses (id, program_id, code, title, credits, tenant_id) 
		VALUES ($1, '00000000-0000-0000-0000-000000000001', 'C1', '{"en": "Course 1"}', 3, $2) ON CONFLICT DO NOTHING`, courseID, tID)
	require.NoError(t, err)

	// 2. Create Module
	mod := &models.CourseModule{
		CourseID: courseID,
		Title:    "Module 1",
		Order:    1,
		IsActive: true,
	}
	err = repo.CreateModule(ctx, mod)
	require.NoError(t, err)
	assert.NotEmpty(t, mod.ID)
	assert.NotZero(t, mod.CreatedAt)

	// 3. Get Module
	fetched, err := repo.GetModule(ctx, mod.ID)
	require.NoError(t, err)
	assert.Equal(t, mod.Title, fetched.Title)

	// 4. Update Module
	mod.Title = "Module 1 Updated"
	err = repo.UpdateModule(ctx, mod)
	require.NoError(t, err)

	fetched, err = repo.GetModule(ctx, mod.ID)
	require.NoError(t, err)
	assert.Equal(t, "Module 1 Updated", fetched.Title)

	// 5. List Modules
	list, err := repo.ListModules(ctx, courseID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// 6. Delete Module
	err = repo.DeleteModule(ctx, mod.ID)
	require.NoError(t, err)

	_, err = repo.GetModule(ctx, mod.ID)
	assert.Error(t, err)
}

func TestSQLCourseContentRepository_Lessons(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()
	repo := NewSQLCourseContentRepository(db)
	ctx := context.Background()

	// Seed Course & Module
	courseID := "00000000-0000-0000-0000-000000000001"
	moduleID := "00000000-0000-0000-0000-000000000002"
	tID := "00000000-0000-0000-0000-000000000001"
	
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type) VALUES ($1, $2, 'T Docs', 'university') ON CONFLICT DO NOTHING`, tID, tID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO programs (id, code, name, tenant_id, title) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'P1', 'Prog', $1, '{"en": "Prog"}') ON CONFLICT DO NOTHING`, tID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO courses (id, program_id, code, title, credits, tenant_id) 
		VALUES ($1, '00000000-0000-0000-0000-000000000001', 'C1', '{"en": "Course 1"}', 3, $2) ON CONFLICT DO NOTHING`, courseID, tID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO course_modules (id, course_id, title, sort_order) VALUES ($1, $2, 'M1', 1)`, moduleID, courseID)
	require.NoError(t, err)

	// Test Lessons
	l := &models.CourseLesson{
		ModuleID: moduleID,
		Title:    "Lesson 1",
		Order:    1,
		IsActive: true,
	}
	err = repo.CreateLesson(ctx, l)
	require.NoError(t, err)
	assert.NotEmpty(t, l.ID)

	l.Title = "Lesson 1 Upd"
	err = repo.UpdateLesson(ctx, l)
	require.NoError(t, err)

	list, err := repo.ListLessons(ctx, moduleID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Lesson 1 Upd", list[0].Title)

	err = repo.DeleteLesson(ctx, l.ID)
	require.NoError(t, err)
}

func TestSQLCourseContentRepository_Activities(t *testing.T) {
	db, cleanup := testutils.SetupTestDB()
	defer cleanup()
	repo := NewSQLCourseContentRepository(db)
	ctx := context.Background()

	// Seed Hierarchy
	courseID := "00000000-0000-0000-0000-000000000001"
	moduleID := "00000000-0000-0000-0000-000000000002"
	lessonID := "00000000-0000-0000-0000-000000000003"
	tID := "00000000-0000-0000-0000-000000000001"
	
	_, err := db.Exec(`INSERT INTO tenants (id, slug, name, tenant_type) VALUES ($1, $2, 'T Docs', 'university') ON CONFLICT DO NOTHING`, tID, tID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO programs (id, code, name, tenant_id, title) 
		VALUES ('00000000-0000-0000-0000-000000000001', 'P1', 'Prog', $1, '{"en": "Prog"}') ON CONFLICT DO NOTHING`, tID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO courses (id, program_id, code, title, credits, tenant_id) 
		VALUES ($1, '00000000-0000-0000-0000-000000000001', 'C1', '{"en": "Course 1"}', 3, $2) ON CONFLICT DO NOTHING`, courseID, tID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO course_modules (id, course_id, title, sort_order) VALUES ($1, $2, 'M1', 1) ON CONFLICT DO NOTHING`, moduleID, courseID)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO course_lessons (id, module_id, title, sort_order) VALUES ($1, $2, 'L1', 1)`, lessonID, moduleID)
	require.NoError(t, err)

	// Test Activities
	a := &models.CourseActivity{
		LessonID: lessonID,
		Type:     "video",
		Title:    "Video 1",
		Order:    1,
		Content:  `{"url":"http://test.com"}`,
	}
	err = repo.CreateActivity(ctx, a)
	require.NoError(t, err)
	assert.NotEmpty(t, a.ID)

	a.Title = "Video 1 Upd"
	err = repo.UpdateActivity(ctx, a)
	require.NoError(t, err)

	list, err := repo.ListActivities(ctx, lessonID)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	err = repo.DeleteActivity(ctx, a.ID)
	require.NoError(t, err)
}
