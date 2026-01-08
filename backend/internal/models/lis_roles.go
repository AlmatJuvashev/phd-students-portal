package models

// IMS Global LIS Role Vocabulary
// Reference: https://www.imsglobal.org/spec/lti/v1p3/#role-vocabularies

const (
	// System roles
	LISAdministrator = "http://purl.imsglobal.org/vocab/lis/v2/system/person#Administrator"
	LISSysAdmin      = "http://purl.imsglobal.org/vocab/lis/v2/system/person#SysAdmin"

	// Institution roles
	LISFaculty = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Faculty"
	LISStaff   = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Staff"
	LISStudent = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student"
	LISAdvisor = "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Advisor"

	// Context (course) roles
	LISInstructor = "http://purl.imsglobal.org/vocab/lis/v2/membership#Instructor"
	LISLearner    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner"
	LISMentor     = "http://purl.imsglobal.org/vocab/lis/v2/membership#Mentor"
	LISContentDev = "http://purl.imsglobal.org/vocab/lis/v2/membership#ContentDeveloper"
	LISManager    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Manager"
	LISMember     = "http://purl.imsglobal.org/vocab/lis/v2/membership#Member"
	LISOfficer    = "http://purl.imsglobal.org/vocab/lis/v2/membership#Officer"
)

// InternalToLIS maps internal roles to IMS LIS URIs
var InternalToLIS = map[Role][]string{
	RoleSuperAdmin:      {LISAdministrator, LISSysAdmin},
	RoleAdmin:           {LISAdministrator},
	RoleHRAdmin:         {LISAdministrator, LISStaff},
	RoleDean:            {LISAdministrator, LISManager},
	RoleChair:           {LISManager, LISOfficer},
	RoleInstructor:      {LISInstructor, LISFaculty},
	RoleAdvisor:         {LISAdvisor, LISMentor, LISFaculty},
	RoleStudent:         {LISLearner, LISStudent},
	RoleExternal:        {LISMember},
	RoleRegistrar:       {LISStaff, LISManager},
	RoleContentMgr:      {LISContentDev},
	RoleFacilityManager: {LISStaff},
	RoleSchedulerAdmin:  {LISStaff, LISManager},
}

// LISToInternal maps LIS roles to internal roles (for LTI launches)
var LISToInternal = map[string]Role{
	LISAdministrator: RoleAdmin,
	LISInstructor:    RoleInstructor,
	LISLearner:       RoleStudent,
	LISMentor:        RoleAdvisor,
	LISContentDev:    RoleContentMgr,
	LISFaculty:       RoleInstructor,
	LISStudent:       RoleStudent,
	LISAdvisor:       RoleAdvisor,
}

// MapLISRolesToInternal converts LTI role claims to internal roles
func MapLISRolesToInternal(lisRoles []string) []Role {
	roleSet := make(map[Role]bool)
	for _, lis := range lisRoles {
		if internal, ok := LISToInternal[lis]; ok {
			roleSet[internal] = true
		}
	}

	roles := make([]Role, 0, len(roleSet))
	for role := range roleSet {
		roles = append(roles, role)
	}
	return roles
}
