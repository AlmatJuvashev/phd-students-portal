package services_test

import (
	"context"
	"strings"
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// -- Mocks --

type MockUserCreator struct {
	mock.Mock
}

func (m *MockUserCreator) CreateUser(ctx context.Context, req services.CreateUserRequest) (*models.User, string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

// -- Tests --

func TestBulkEnrollmentService_ImportStudents(t *testing.T) {
	mockUC := new(MockUserCreator)
	svc := services.NewBulkEnrollmentService(mockUC)
	ctx := context.Background()

	// CSV Data: Header + 2 valid rows + 1 invalid row (missing email)
	csvData := `first_name,last_name,email,role
John,Doe,john@example.com,student
Jane,Smith,jane@example.com,undergrad
Invalid,Row`

	reader := strings.NewReader(csvData)
	tenantID := "tenant-1"

	// Mock Expectations
	mockUC.On("CreateUser", ctx, mock.MatchedBy(func(req services.CreateUserRequest) bool {
		return req.Email == "john@example.com" && req.FirstName == "John" && req.TenantID == "tenant-1"
	})).Return(&models.User{ID: "u1"}, "tempPass1", nil)

	mockUC.On("CreateUser", ctx, mock.MatchedBy(func(req services.CreateUserRequest) bool {
		return req.Email == "jane@example.com" && req.Role == "undergrad"
	})).Return(&models.User{ID: "u2"}, "tempPass2", nil)

	// Act
	count, errors := svc.ImportStudents(ctx, reader, tenantID)

	// Assert
	assert.Equal(t, 2, count)
	
	// Should adhere to error for the 3rd row (index 2 in 0-based iteration over body, or line 3)
	// Actually CSV reader: ReadAll returns all records including header if not stripped.
	// Service logic skips row 0 (header).
	// Row 1: John (Valid)
	// Row 2: Jane (Valid)
	// Row 3: Invalid (Insufficient columns) -> Error
	
	assert.Len(t, errors, 1)
	assert.Contains(t, errors[0].Error(), "insufficient columns")

	mockUC.AssertExpectations(t)
}
