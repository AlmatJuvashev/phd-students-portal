package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
)

type UserCreator interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, string, error)
}

type BulkEnrollmentService struct {
	userService UserCreator
}

func NewBulkEnrollmentService(userService UserCreator) *BulkEnrollmentService {
	return &BulkEnrollmentService{userService: userService}
}

type BulkStudentRow struct {
	FirstName string
	LastName  string
	Email     string
	Role      string // e.g. "student", "undergrad"
	TenantID  string
}

// ImportStudents processes a CSV and creates users.
// Expected CSV Format: first_name, last_name, email, role
func (s *BulkEnrollmentService) ImportStudents(ctx context.Context, r io.Reader, tenantID string) (int, []error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1 // Allow variable fields, we validate manually
	
	createdCount := 0
	var errorList []error
	
	// Skip Header
	_, err := reader.Read()
	if err != nil {
		return 0, []error{err}
	}

	rowIndex := 1 // 0-based index of data rows (header is -1 relative to loop?) 
	// Actually CSV line numbers: Header=1. First Data=2.
	
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowIndex++
		
		if err != nil {
			errorList = append(errorList, fmt.Errorf("row %d: %v", rowIndex, err))
			continue
		}

		if len(row) < 3 {
			errorList = append(errorList, fmt.Errorf("row %d: insufficient columns (need first,last,email)", rowIndex))
			continue
		}

		firstName := strings.TrimSpace(row[0])
		lastName := strings.TrimSpace(row[1])
		email := strings.TrimSpace(row[2])
		role := "student"
		if len(row) > 3 && row[3] != "" {
			role = strings.TrimSpace(row[3])
		}

		// Prepare CreateUserRequest (reusing existing struct if public, or adapting)
		createUserReq := CreateUserRequest{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Role:      role,
			TenantID:  tenantID,
			// Password? Auto-generate or set default. UserService typically handles this or validation.
			// Ideally UserService.CreateUser generates a password reset token/email.
		}

		// Call UserService
		// We need to check UserService signature.
		// Usually (ctx, req).
		var createErr error
		_, _, createErr = s.userService.CreateUser(ctx, createUserReq)
		if createErr != nil {
			errorList = append(errorList, fmt.Errorf("row %d (%s): %v", rowIndex, email, createErr))
			continue
		}
		createdCount++
	}

	return createdCount, errorList
}
