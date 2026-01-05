package dto

// StudentCourseDetail is a student-facing view for a single course offering.
type StudentCourseDetail struct {
	Course    StudentCourse            `json:"course"`
	Sessions  []StudentCourseNextSession `json:"sessions"`
}

