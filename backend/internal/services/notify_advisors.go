package services

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// NotifyAdvisorsOnSubmission creates admin_notifications for all advisors
// assigned to the student when a document is submitted for review.
// This creates a shared notification that all assigned advisors can see.
func NotifyAdvisorsOnSubmission(db *sqlx.DB, studentID, nodeID, nodeInstanceID, message string) error {
	// Get student name
	var studentName string
	err := db.Get(&studentName, `SELECT COALESCE(first_name || ' ' || last_name, email, username) 
		FROM users WHERE id=$1`, studentID)
	if err != nil {
		log.Printf("[NotifyAdvisors] Failed to get student name: %v", err)
		studentName = "A student"
	}

	// Get all advisors assigned to this student
	type advisor struct {
		ID string `db:"advisor_id"`
	}
	var advisors []advisor
	err = db.Select(&advisors, `SELECT advisor_id FROM student_advisors WHERE student_id=$1`, studentID)
	if err != nil {
		log.Printf("[NotifyAdvisors] Failed to get advisors: %v", err)
		return err
	}

	if len(advisors) == 0 {
		log.Printf("[NotifyAdvisors] No advisors assigned to student %s, skipping notification", studentID)
		return nil
	}

	// Build notification message
	if message == "" {
		message = studentName + " submitted a document for review"
	}

	// Insert notification into admin_notifications
	// All advisors will see this notification through the list endpoint
	// The notification is associated with the student, so advisors filtering by their students will see it
	_, err = db.Exec(`INSERT INTO admin_notifications 
		(student_id, node_id, node_instance_id, event_type, message, metadata)
		VALUES ($1, $2, $3, 'document_submitted', $4, '{}')`,
		studentID, nodeID, nodeInstanceID, message)
	if err != nil {
		log.Printf("[NotifyAdvisors] Failed to insert notification: %v", err)
		return err
	}

	log.Printf("[NotifyAdvisors] Created submission notification for student %s, node %s, %d advisors assigned",
		studentID, nodeID, len(advisors))
	return nil
}

// GetAdvisorsForStudent returns all advisor IDs for a given student
func GetAdvisorsForStudent(db *sqlx.DB, studentID string) ([]string, error) {
	var advisorIDs []string
	err := db.Select(&advisorIDs, `SELECT advisor_id FROM student_advisors WHERE student_id=$1`, studentID)
	return advisorIDs, err
}

// HasAdvisors checks if a student has any advisors assigned
func HasAdvisors(db *sqlx.DB, studentID string) (bool, error) {
	var count int
	err := db.Get(&count, `SELECT COUNT(*) FROM student_advisors WHERE student_id=$1`, studentID)
	return count > 0, err
}
