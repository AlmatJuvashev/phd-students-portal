package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Local Mocks --

type MockTranscriptRepo struct {
	mock.Mock
}

func (m *MockTranscriptRepo) GetStudentGrades(ctx context.Context, studentID string) ([]models.TermGrade, error) {
	args := m.Called(ctx, studentID)
	return args.Get(0).([]models.TermGrade), args.Error(1)
}

type MockSchedulerRepo struct {
	mock.Mock
}

func (m *MockSchedulerRepo) GetTerm(ctx context.Context, id string) (*models.AcademicTerm, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AcademicTerm), args.Error(1)
}

// Stubs for unused methods to satisfy interface
// Stubs converted to Mocks
func (m *MockSchedulerRepo) CreateTerm(ctx context.Context, term *models.AcademicTerm) error {
	args := m.Called(ctx, term)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListTerms(ctx context.Context, tenantID string) ([]models.AcademicTerm, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.AcademicTerm), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateTerm(ctx context.Context, term *models.AcademicTerm) error {
	args := m.Called(ctx, term)
	return args.Error(0)
}
func (m *MockSchedulerRepo) DeleteTerm(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) CreateOffering(ctx context.Context, offering *models.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}
func (m *MockSchedulerRepo) GetOffering(ctx context.Context, id string) (*models.CourseOffering, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*models.CourseOffering), args.Error(1)
}
func (m *MockSchedulerRepo) ListOfferings(ctx context.Context, tenantID string, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, tenantID, termID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *MockSchedulerRepo) ListOfferingsByInstructor(ctx context.Context, instructorID string, termID string) ([]models.CourseOffering, error) {
	args := m.Called(ctx, instructorID, termID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseOffering), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateOffering(ctx context.Context, offering *models.CourseOffering) error {
	args := m.Called(ctx, offering)
	return args.Error(0)
}
func (m *MockSchedulerRepo) AddStaff(ctx context.Context, staff *models.CourseStaff) error {
	args := m.Called(ctx, staff)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.CourseStaff), args.Error(1)
}
func (m *MockSchedulerRepo) RemoveStaff(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) CreateSession(ctx context.Context, session *models.ClassSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}
func (m *MockSchedulerRepo) ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, offeringID, startDate, endDate)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsByRoom(ctx context.Context, roomID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, roomID, startDate, endDate)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsByInstructor(ctx context.Context, instructorID string, startDate, endDate time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, instructorID, startDate, endDate)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsForTerm(ctx context.Context, termID string) ([]models.ClassSession, error) {
	args := m.Called(ctx, termID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}
func (m *MockSchedulerRepo) UpdateSession(ctx context.Context, session *models.ClassSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}
func (m *MockSchedulerRepo) DeleteSession(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockSchedulerRepo) AddCohortToOffering(ctx context.Context, offeringID, cohortID string) error {
	args := m.Called(ctx, offeringID, cohortID)
	return args.Error(0)
}
func (m *MockSchedulerRepo) GetOfferingCohorts(ctx context.Context, offeringID string) ([]string, error) {
	args := m.Called(ctx, offeringID)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockSchedulerRepo) ListSessionsForCohorts(ctx context.Context, cohortIDs []string, startTime, endTime time.Time) ([]models.ClassSession, error) {
	args := m.Called(ctx, cohortIDs, startTime, endTime)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).([]models.ClassSession), args.Error(1)
}

// -- Tests --

func TestTranscriptService_GetTranscript(t *testing.T) {
	mockTR := new(MockTranscriptRepo)
	mockSR := new(MockSchedulerRepo)
	svc := services.NewTranscriptService(mockTR, mockSR)
	ctx := context.Background()

	// Term 1: Fall 2025
	term1ID := "term-1"
	term1 := &models.AcademicTerm{ID: term1ID, Name: "Fall 2025"}
	
	// Term 2: Spring 2026
	term2ID := "term-2"
	term2 := &models.AcademicTerm{ID: term2ID, Name: "Spring 2026"}

	// Grades
	// Term 1: 2 Courses. 4.0 (3 credits), 3.0 (3 credits). Term GPA = (12+9)/6 = 3.5
	g1 := models.TermGrade{TermID: term1ID, CourseCode: "CS101", Credits: 3, GradePoints: 4.0, Grade: "A"}
	g2 := models.TermGrade{TermID: term1ID, CourseCode: "CS102", Credits: 3, GradePoints: 3.0, Grade: "B"}
	
	// Term 2: 1 Course. 2.0 (4 credits). Term GPA = 8/4 = 2.0
	g3 := models.TermGrade{TermID: term2ID, CourseCode: "CS201", Credits: 4, GradePoints: 2.0, Grade: "C"}

	mockTR.On("GetStudentGrades", ctx, "student-123").Return([]models.TermGrade{g1, g2, g3}, nil)
	
	// Mock expects calls for distinct terms
	mockSR.On("GetTerm", ctx, term1ID).Return(term1, nil)
	mockSR.On("GetTerm", ctx, term2ID).Return(term2, nil)

	// Act
	transcript, err := svc.GetTranscript(ctx, "student-123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "student-123", transcript.StudentID)
	assert.Equal(t, 2, len(transcript.Terms))

	// Verify Cumulative
	// Total Points: (4*3) + (3*3) + (2*4) = 12 + 9 + 8 = 29
	// Total Credits: 3 + 3 + 4 = 10
	// Cumulative GPA: 29 / 10 = 2.9
	assert.Equal(t, float64(10), transcript.TotalCredits)
	assert.Equal(t, float64(29), transcript.TotalPoints)
	assert.Equal(t, float32(2.9), transcript.CumulativeGPA)

	// Verify Terms Order implicitly (repo returns in order, service preserves)
	// We mocked return array {g1, g2, g3}. g1/g2 are term1. g3 is term2.
	// Map iteration order is random, but service logic specifically preserves order based on first appearance in grades list.
	// Since g1 comes first, term1 should be first.
	assert.Equal(t, "Fall 2025", transcript.Terms[0].TermName)
	assert.Equal(t, float32(3.5), transcript.Terms[0].TermGPA)
	
	assert.Equal(t, "Spring 2026", transcript.Terms[1].TermName)
	assert.Equal(t, float32(2.0), transcript.Terms[1].TermGPA)
}
