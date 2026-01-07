package services_test

import (
	"context"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) BatchUpsertAttendance(ctx context.Context, sessionID string, records []models.ClassAttendance, recordedBy string) error {
	args := m.Called(ctx, sessionID, records, recordedBy)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetSessionAttendance(ctx context.Context, sessionID string) ([]models.ClassAttendance, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

func (m *MockAttendanceRepo) GetStudentAttendance(ctx context.Context, studentID string) ([]models.ClassAttendance, error) {
    args := m.Called(ctx, studentID)
    return args.Get(0).([]models.ClassAttendance), args.Error(1)
}

func (m *MockAttendanceRepo) RecordAttendance(ctx context.Context, sessionID string, record models.ClassAttendance) error {
	return nil
}

// -- Tests --

func TestAttendanceService_BatchRecordAttendance(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	svc := services.NewAttendanceService(mockRepo)
	ctx := context.Background()

	sessionID := "session-123"
	teacherID := "teacher-456"
	
	updates := []models.ClassAttendance{
		{StudentID: "s1", Status: "PRESENT", Notes: ""},
		{StudentID: "s2", Status: "ABSENT", Notes: "Sick"},
	}

	// Expectation
	mockRepo.On("BatchUpsertAttendance", ctx, sessionID, updates, teacherID).Return(nil)

	// Act
	err := svc.BatchRecordAttendance(ctx, sessionID, updates, teacherID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
func TestAttendanceService_GetSessionAttendance(t *testing.T) {
	mockRepo := new(MockAttendanceRepo)
	svc := services.NewAttendanceService(mockRepo)
	ctx := context.Background()

	sessionID := "session-123"
	expected := []models.ClassAttendance{{StudentID: "s1", Status: "PRESENT"}}

	mockRepo.On("GetSessionAttendance", ctx, sessionID).Return(expected, nil)

	res, err := svc.GetSessionAttendance(ctx, sessionID)

	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}
