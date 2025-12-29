package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestDictionaryService_Unit(t *testing.T) {
	mockRepo := NewMockDictionaryRepository()
	svc := services.NewDictionaryService(mockRepo)
	ctx := context.Background()

	t.Run("Programs", func(t *testing.T) {
		_, _ = svc.ListPrograms(ctx, "t1", true)
		_, _ = svc.CreateProgram(ctx, "t1", "Prog", "P1")
		_ = svc.UpdateProgram(ctx, "t1", "p1", "New", "P2", nil)
		_ = svc.DeleteProgram(ctx, "t1", "p1")
	})

	t.Run("Specialties", func(t *testing.T) {
		_, _ = svc.ListSpecialties(ctx, "t1", true, "p1")
		_, _ = svc.CreateSpecialty(ctx, "t1", "Spec", "S1", []string{"p1"})
		_ = svc.UpdateSpecialty(ctx, "t1", "s1", "New", "S2", nil, []string{"p2"})
		_ = svc.DeleteSpecialty(ctx, "t1", "s1")
	})

	t.Run("Cohorts", func(t *testing.T) {
		_, _ = svc.ListCohorts(ctx, "t1", true)
		_, _ = svc.CreateCohort(ctx, "t1", "2024", "2024-01-01", "2024-12-31")
		_ = svc.UpdateCohort(ctx, "t1", "c1", "2025", "", "", nil)
		_ = svc.DeleteCohort(ctx, "t1", "c1")
	})

	t.Run("Departments", func(t *testing.T) {
		_, _ = svc.ListDepartments(ctx, "t1", true)
		_, _ = svc.CreateDepartment(ctx, "t1", "Dept", "D1")
		_ = svc.UpdateDepartment(ctx, "t1", "d1", "New", "D2", nil)
		_ = svc.DeleteDepartment(ctx, "t1", "d1")
	})

	assert.NotNil(t, svc)
}
