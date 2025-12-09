package permissions

import (
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
)

type Action string

const (
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type Resource string

const (
	ResourceUser     Resource = "user"
	ResourceEvent    Resource = "event"
	ResourceDocument Resource = "document"
)

// Can checks if the actor can perform the action on the resource.
// This is a simplified policy engine.
func Can(actor models.User, action Action, resource Resource, target interface{}) bool {
	switch resource {
	case ResourceUser:
		return canUser(actor, action, target)
	case ResourceEvent:
		return canEvent(actor, action, target)
	default:
		return false
	}
}

func canUser(actor models.User, action Action, target interface{}) bool {
	targetUser, ok := target.(models.User)
	if !ok {
		// If target is nil or not a user, maybe we are checking general permission
		if target == nil {
			// e.g. Can create user? Only admin
			if action == ActionCreate {
				return actor.Role == models.RoleAdmin || actor.Role == models.RoleSuperAdmin
			}
		}
		return false
	}

	// Superadmin can do anything
	if actor.Role == models.RoleSuperAdmin {
		return true
	}

	// Admin can manage users
	if actor.Role == models.RoleAdmin {
		return true
	}

	// Users can read/update themselves
	if actor.ID == targetUser.ID {
		return true
	}

	// Advisors can read their students (assuming we had advisor relation logic here)
	// For now, let's say Advisors can read any student (simplified)
	if actor.Role == models.RoleAdvisor && targetUser.Role == models.RoleStudent {
		return action == ActionRead
	}

	return false
}

func canEvent(actor models.User, action Action, target interface{}) bool {
	targetEvent, ok := target.(models.Event)
	if !ok {
		if target == nil {
			// Can create event?
			return true // Everyone can create events (maybe restricted later)
		}
		return false
	}

	if actor.Role == models.RoleSuperAdmin || actor.Role == models.RoleAdmin {
		return true
	}

	// Creator can do anything
	if actor.ID == targetEvent.CreatorID {
		return true
	}

	// Attendees can read
	if action == ActionRead {
		// We need to check attendance list, but that's not in Event struct usually (it's separate)
		// Assuming if they have access to the event object, they can read it?
		// Or we rely on service layer to filter.
		return true
	}

	return false
}
