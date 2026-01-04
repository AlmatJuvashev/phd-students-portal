package dto

type ProgramProgress struct {
	Title           string  `json:"title"`
	ProgressPercent float64 `json:"progress_percent"`
	CompletedNodes  int     `json:"completed_nodes"`
	TotalNodes      int     `json:"total_nodes"`
	OverdueCount    int     `json:"overdue_count"`
}

type StudentDeadline struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	DueAt    *string `json:"due_at,omitempty"`
	Source   string  `json:"source"` // "journey" | "course"
	Status   string  `json:"status,omitempty"`
	Severity string  `json:"severity,omitempty"` // "urgent" | "normal"
	Link     *string `json:"link,omitempty"`
}

type StudentAnnouncement struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Body    string  `json:"body"`
	Created string  `json:"created_at"`
	Link    *string `json:"link,omitempty"`
}

type StudentGradeEntry struct {
	ID               string  `json:"id"`
	CourseOfferingID string  `json:"course_offering_id"`
	CourseID         string  `json:"course_id,omitempty"`
	CourseCode       *string `json:"course_code,omitempty"`
	CourseTitle      *string `json:"course_title,omitempty"`
	ActivityID       string  `json:"activity_id"`
	StudentID        string  `json:"student_id"`
	Score            float64 `json:"score"`
	MaxScore         float64 `json:"max_score"`
	Grade            string  `json:"grade"`
	Feedback         string  `json:"feedback"`
	GradedByID       string  `json:"graded_by_id"`
	GradedAt         string  `json:"graded_at"`
}

type StudentCourseNextSession struct {
	ID         string  `json:"id"`
	Date       string  `json:"date"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	RoomID     *string `json:"room_id,omitempty"`
	MeetingURL *string `json:"meeting_url,omitempty"`
	Type       string  `json:"type"`
}

type StudentCourse struct {
	EnrollmentID     string                    `json:"enrollment_id"`
	CourseOfferingID string                    `json:"course_offering_id"`
	Status           string                    `json:"status"`
	CourseID         string                    `json:"course_id"`
	Code             string                    `json:"code"`
	Title            string                    `json:"title"`
	Section          string                    `json:"section"`
	TermID           string                    `json:"term_id"`
	DeliveryFormat   string                    `json:"delivery_format"`
	InstructorName   *string                   `json:"instructor_name,omitempty"`
	ProgressPercent  float64                   `json:"progress_percent"`
	NextSession      *StudentCourseNextSession `json:"next_session,omitempty"`
}

type StudentAssignment struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Source   string  `json:"source"` // "journey" | "course"
	Status   string  `json:"status"`
	DueAt    *string `json:"due_at,omitempty"`
	Link     *string `json:"link,omitempty"`
	Severity string  `json:"severity,omitempty"`
}

type StudentDashboard struct {
	Program           ProgramProgress       `json:"program"`
	UpcomingDeadlines []StudentDeadline     `json:"upcoming_deadlines"`
	RecentGrades      []StudentGradeEntry   `json:"recent_grades"`
	Announcements     []StudentAnnouncement `json:"announcements"`
}
