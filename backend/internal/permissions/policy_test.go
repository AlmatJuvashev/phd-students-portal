package permissions

import (
	"testing"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCan_ResourceUser(t *testing.T) {
	superAdmin := models.User{ID: "super", Role: models.RoleSuperAdmin}
	admin := models.User{ID: "admin", Role: models.RoleAdmin}
	student := models.User{ID: "student1", Role: models.RoleStudent}
	otherStudent := models.User{ID: "student2", Role: models.RoleStudent}
	advisor := models.User{ID: "advisor", Role: models.RoleAdvisor}

	tests := []struct {
		name     string
		actor    models.User
		action   Action
		target   interface{}
		expected bool
	}{
		// Superadmin
		{"Superadmin can create user", superAdmin, ActionCreate, nil, true},
		{"Superadmin can update anyone", superAdmin, ActionUpdate, student, true},
		
		// Admin
		{"Admin can create user", admin, ActionCreate, nil, true},
		{"Admin can update student", admin, ActionUpdate, student, true},
		
		// Student
		{"Student cannot create user", student, ActionCreate, nil, false},
		{"Student can read self", student, ActionRead, student, true},
		{"Student can update self", student, ActionUpdate, student, true},
		{"Student cannot update other", student, ActionUpdate, otherStudent, false},
		
		// Advisor
		{"Advisor can read student", advisor, ActionRead, student, true},
		{"Advisor cannot update student", advisor, ActionUpdate, student, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Can(tt.actor, tt.action, ResourceUser, tt.target))
		})
	}
}

func TestCan_ResourceEvent(t *testing.T) {
	admin := models.User{ID: "admin", Role: models.RoleAdmin}
	creator := models.User{ID: "creator", Role: models.RoleStudent}
	other := models.User{ID: "other", Role: models.RoleStudent}

	event := models.Event{ID: "evt1", CreatorID: "creator"}

	tests := []struct {
		name     string
		actor    models.User
		action   Action
		target   interface{}
		expected bool
	}{
		// Create
		{"Anyone can create event", other, ActionCreate, nil, true},
		
		// Admin
		{"Admin can update event", admin, ActionUpdate, event, true},
		
		// Creator
		{"Creator can update own event", creator, ActionUpdate, event, true},
		
		// Other
		{"Other can read event", other, ActionRead, event, true},
		{"Other cannot update event", other, ActionUpdate, event, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Can(tt.actor, tt.action, ResourceEvent, tt.target))
		})
	}
}

func TestCan_UnknownResource(t *testing.T) {
	actor := models.User{ID: "u1", Role: models.RoleAdmin}
	assert.False(t, Can(actor, ActionRead, "unknown_resource", nil))
}
