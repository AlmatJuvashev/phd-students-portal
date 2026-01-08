package models_test

import (
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMapLISRolesToInternal(t *testing.T) {
	// Define LIS constants locally for test until they are implemented
	// or rely on them being present (strict TDD says fail compilation)
	
	tests := []struct {
		name     string
		lisRoles []string
		expected []models.Role
	}{
		{
			name:     "Single Known Role - Instructor",
			lisRoles: []string{models.LISInstructor}, // Expected strict failure: undefined
			expected: []models.Role{models.RoleInstructor}, // Expected strict failure: undefined
		},
		{
			name:     "check LISLearner maps to Student",
			lisRoles: []string{models.LISLearner},
			expected: []models.Role{models.RoleStudent},
		},
		{
			name:     "Multiple Known Roles",
			lisRoles: []string{models.LISInstructor, models.LISAdvisor},
			expected: []models.Role{models.RoleInstructor, models.RoleAdvisor},
		},
		{
			name:     "Role with multiple LIS mappings (Instructor matches Faculty too)",
			lisRoles: []string{models.LISFaculty},
			expected: []models.Role{models.RoleInstructor},
		},
		{
			name:     "Unknown Role - Ignored",
			lisRoles: []string{"http://purl.imsglobal.org/vocab/lis/v2/membership#Unknown"},
			expected: []models.Role{},
		},
		{
			name:     "Duplicate Result Roles (Deduplication)",
			lisRoles: []string{models.LISInstructor, models.LISFaculty}, // Both map to Instructor
			expected: []models.Role{models.RoleInstructor},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := models.MapLISRolesToInternal(tt.lisRoles) // Expected strict failure: undefined
			assert.ElementsMatch(t, tt.expected, got)
		})
	}
}
