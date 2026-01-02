package services

import (
	"context"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type TeacherService struct {
	schedulerRepo repository.SchedulerRepository
	lmsRepo       repository.LMSRepository
	gradingRepo   repository.GradingRepository
}

func NewTeacherService(schedulerRepo repository.SchedulerRepository, lmsRepo repository.LMSRepository, gradingRepo repository.GradingRepository) *TeacherService {
	return &TeacherService{
		schedulerRepo: schedulerRepo,
		lmsRepo:       lmsRepo,
		gradingRepo:   gradingRepo,
	}
}

// TeacherDashboardStats aggregates validation counts for the dashboard
type TeacherDashboardStats struct {
	NextClass        *models.ClassSession `json:"next_class,omitempty"`
	ActiveCourses    int                  `json:"active_courses"`
	PendingGrading   int                  `json:"pending_grading"`
	TodayClassesCount int                 `json:"today_classes_count"`
}

func (s *TeacherService) GetDashboardStats(ctx context.Context, instructorID string) (*TeacherDashboardStats, error) {
	stats := &TeacherDashboardStats{}

	// 1. Get Today's Classes
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	sessions, err := s.schedulerRepo.ListSessionsByInstructor(ctx, instructorID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	stats.TodayClassesCount = len(sessions)

	// Determine next class (first session after now)
	for _, sess := range sessions {
		sessTime, _ := time.Parse("15:04", sess.StartTime) // simplified
		sessDateTime := time.Date(sess.Date.Year(), sess.Date.Month(), sess.Date.Day(), sessTime.Hour(), sessTime.Minute(), 0, 0, now.Location())
		if sessDateTime.After(now) {
			val := sess // Active copy
			stats.NextClass = &val
			break
		}
	}

	// 2. Count Active Courses
	// Assume empty termID gets all, filtered by current date/active in repo if needed.
	// For now, list all.
	courses, err := s.schedulerRepo.ListOfferingsByInstructor(ctx, instructorID, "")
	if err != nil {
		return nil, err
	}
	stats.ActiveCourses = len(courses)

	// 3. Pending Grading
	// Iterate courses -> list submissions -> count 'SUBMITTED' status
	// Optimization: This N+1 query pattern is bad for scale, but OK for MVP with few courses.
	// Ideal: `lmsRepo.CountPendingSubmissions(instructorID)`
	pendingCount := 0
	for _, course := range courses {
		subs, err := s.lmsRepo.ListSubmissions(ctx, course.ID)
		if err == nil {
			for _, sub := range subs {
				if sub.Status == "SUBMITTED" {
					pendingCount++
				}
			}
		}
	}
	stats.PendingGrading = pendingCount

	return stats, nil
}

func (s *TeacherService) GetMyCourses(ctx context.Context, instructorID string) ([]models.CourseOffering, error) {
	return s.schedulerRepo.ListOfferingsByInstructor(ctx, instructorID, "")
}

func (s *TeacherService) GetCourseRoster(ctx context.Context, offeringID string) ([]models.CourseEnrollment, error) {
	return s.lmsRepo.GetCourseRoster(ctx, offeringID)
}

func (s *TeacherService) GetGradebook(ctx context.Context, offeringID string) ([]models.GradebookEntry, error) {
	return s.gradingRepo.ListEntries(ctx, offeringID)
}

func (s *TeacherService) GetSubmissions(ctx context.Context, instructorID string) ([]models.ActivitySubmission, error) {
	// Filter by instructor's courses
	courses, err := s.schedulerRepo.ListOfferingsByInstructor(ctx, instructorID, "")
	if err != nil {
		return nil, err
	}

	var allSubs []models.ActivitySubmission
	for _, c := range courses {
		subs, err := s.lmsRepo.ListSubmissions(ctx, c.ID)
		if err == nil {
			allSubs = append(allSubs, subs...)
		}
	}
	return allSubs, nil
}

// --- Annotations ---

func (s *TeacherService) AddAnnotation(ctx context.Context, ann models.SubmissionAnnotation) (*models.SubmissionAnnotation, error) {
	return s.lmsRepo.CreateAnnotation(ctx, ann)
}

func (s *TeacherService) GetAnnotationsForSubmission(ctx context.Context, submissionID string) ([]models.SubmissionAnnotation, error) {
	return s.lmsRepo.ListAnnotations(ctx, submissionID)
}

func (s *TeacherService) RemoveAnnotation(ctx context.Context, id string) error {
	return s.lmsRepo.DeleteAnnotation(ctx, id)
}
