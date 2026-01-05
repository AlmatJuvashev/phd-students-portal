package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/dto"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type TeacherService struct {
	schedulerRepo repository.SchedulerRepository
	lmsRepo       repository.LMSRepository
	gradingRepo   repository.GradingRepository
	contentRepo   repository.CourseContentRepository
}

func NewTeacherService(schedulerRepo repository.SchedulerRepository, lmsRepo repository.LMSRepository, gradingRepo repository.GradingRepository, contentRepo repository.CourseContentRepository) *TeacherService {
	return &TeacherService{
		schedulerRepo: schedulerRepo,
		lmsRepo:       lmsRepo,
		gradingRepo:   gradingRepo,
		contentRepo:   contentRepo,
	}
}

// TeacherDashboardStats aggregates validation counts for the dashboard
type TeacherDashboardStats struct {
	NextClass         *models.ClassSession `json:"next_class,omitempty"`
	ActiveCourses     int                  `json:"active_courses"`
	PendingGrading    int                  `json:"pending_grading"`
	TodayClassesCount int                  `json:"today_classes_count"`
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

// --- Student Tracker ---

func (s *TeacherService) GetCourseStudents(ctx context.Context, offeringID string) ([]dto.StudentRiskProfile, error) {
	roster, err := s.lmsRepo.GetCourseRoster(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	assignmentsTotal, err := s.countOfferingWorkItems(ctx, offeringID)
	if err != nil {
		// Degrade gracefully: tracker still works from submissions/grades.
		assignmentsTotal = 0
	}

	submissions, err := s.lmsRepo.ListSubmissions(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	studentSubmissionCounts := map[string]int{}
	studentLastActivity := map[string]time.Time{}
	seenActivityPerStudent := map[string]map[string]bool{}
	distinctActivitiesInOffering := map[string]bool{}

	for _, sub := range submissions {
		distinctActivitiesInOffering[sub.ActivityID] = true
		if _, ok := seenActivityPerStudent[sub.StudentID]; !ok {
			seenActivityPerStudent[sub.StudentID] = map[string]bool{}
		}
		if sub.Status == "SUBMITTED" || sub.Status == "GRADED" {
			// Count unique completed activities.
			if !seenActivityPerStudent[sub.StudentID][sub.ActivityID] {
				seenActivityPerStudent[sub.StudentID][sub.ActivityID] = true
				studentSubmissionCounts[sub.StudentID]++
			}
		}
		if sub.SubmittedAt.After(studentLastActivity[sub.StudentID]) {
			studentLastActivity[sub.StudentID] = sub.SubmittedAt
		}
	}

	if assignmentsTotal == 0 {
		assignmentsTotal = len(distinctActivitiesInOffering)
	}

	grades, err := s.gradingRepo.ListEntries(ctx, offeringID)
	if err != nil {
		return nil, err
	}

	type gradeAgg struct {
		sum   float64
		count int
	}
	studentGrades := map[string]gradeAgg{}
	for _, g := range grades {
		if g.MaxScore <= 0 {
			continue
		}
		agg := studentGrades[g.StudentID]
		agg.sum += (g.Score / g.MaxScore) * 100.0
		agg.count++
		studentGrades[g.StudentID] = agg
	}

	now := time.Now()
	out := make([]dto.StudentRiskProfile, 0, len(roster))
	for _, enr := range roster {
		studentName := strings.TrimSpace(enr.StudentName)
		if studentName == "" {
			studentName = strings.TrimSpace(enr.StudentEmail)
		}

		assignmentsCompleted := studentSubmissionCounts[enr.StudentID]
		progress := 0.0
		if assignmentsTotal > 0 {
			progress = (float64(assignmentsCompleted) / float64(assignmentsTotal)) * 100.0
		}

		last := studentLastActivity[enr.StudentID]
		lastISO := ""
		daysInactive := 999
		if !last.IsZero() {
			lastISO = last.UTC().Format(time.RFC3339)
			daysInactive = int(now.Sub(last).Hours() / 24)
		}

		avgGrade := 0.0
		if agg, ok := studentGrades[enr.StudentID]; ok && agg.count > 0 {
			avgGrade = agg.sum / float64(agg.count)
		}

		riskLevel, riskFactors, actions := evaluateStudentRisk(progress, daysInactive, avgGrade, assignmentsCompleted, assignmentsTotal)

		out = append(out, dto.StudentRiskProfile{
			StudentID:            enr.StudentID,
			StudentName:          studentName,
			OverallProgress:      progress,
			AssignmentsCompleted: assignmentsCompleted,
			AssignmentsTotal:     assignmentsTotal,
			AssignmentsOverdue:   0,
			LastActivity:         lastISO,
			DaysInactive:         daysInactive,
			AverageGrade:         avgGrade,
			RiskLevel:            riskLevel,
			RiskFactors:          riskFactors,
			SuggestedActions:     actions,
		})
	}

	return out, nil
}

func (s *TeacherService) GetCourseAtRisk(ctx context.Context, offeringID string) ([]dto.StudentRiskProfile, error) {
	students, err := s.GetCourseStudents(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	var out []dto.StudentRiskProfile
	for _, st := range students {
		if st.RiskLevel == "high" || st.RiskLevel == "critical" {
			out = append(out, st)
		}
	}
	return out, nil
}

func (s *TeacherService) GetStudentActivity(ctx context.Context, studentID string, offeringID string, limit int) ([]dto.TeacherStudentActivityEvent, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var events []dto.TeacherStudentActivityEvent

	submissions, err := s.lmsRepo.ListSubmissions(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	for _, sub := range submissions {
		if sub.StudentID != studentID {
			continue
		}
		title := strings.TrimSpace(sub.ActivityTitle)
		if title == "" {
			title = fmt.Sprintf("Activity %s", sub.ActivityID)
		}
		activityID := sub.ActivityID
		subID := sub.ID
		status := sub.Status
		events = append(events, dto.TeacherStudentActivityEvent{
			Kind:         "submission",
			OccurredAt:   sub.SubmittedAt.UTC().Format(time.RFC3339),
			Title:        title,
			Status:       &status,
			ActivityID:   &activityID,
			SubmissionID: &subID,
		})
	}

	grades, err := s.gradingRepo.ListEntries(ctx, offeringID)
	if err != nil {
		return nil, err
	}
	for _, g := range grades {
		if g.StudentID != studentID {
			continue
		}
		activityID := g.ActivityID
		score := g.Score
		maxScore := g.MaxScore
		grade := g.Grade
		events = append(events, dto.TeacherStudentActivityEvent{
			Kind:       "grade",
			OccurredAt: g.GradedAt.UTC().Format(time.RFC3339),
			Title:      "Grade recorded",
			ActivityID: &activityID,
			Score:      &score,
			MaxScore:   &maxScore,
			Grade:      &grade,
		})
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].OccurredAt > events[j].OccurredAt
	})
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

func (s *TeacherService) countOfferingWorkItems(ctx context.Context, offeringID string) (int, error) {
	if s.contentRepo == nil {
		return 0, nil
	}
	offering, err := s.schedulerRepo.GetOffering(ctx, offeringID)
	if err != nil {
		return 0, err
	}
	modules, err := s.contentRepo.ListModules(ctx, offering.CourseID)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, m := range modules {
		lessons, err := s.contentRepo.ListLessons(ctx, m.ID)
		if err != nil {
			return 0, err
		}
		for _, l := range lessons {
			acts, err := s.contentRepo.ListActivities(ctx, l.ID)
			if err != nil {
				return 0, err
			}
			for _, a := range acts {
				if a.IsOptional {
					continue
				}
				if a.Type == "assignment" || a.Type == "quiz" {
					total++
				}
			}
		}
	}
	return total, nil
}

func evaluateStudentRisk(progress float64, daysInactive int, avgGrade float64, completed int, total int) (string, []string, []string) {
	riskLevel := "low"
	var factors []string
	var actions []string

	upgrade := func(target string) {
		order := map[string]int{"low": 1, "medium": 2, "high": 3, "critical": 4}
		if order[target] > order[riskLevel] {
			riskLevel = target
		}
	}

	if daysInactive >= 30 {
		upgrade("critical")
		factors = append(factors, fmt.Sprintf("%d days inactive", daysInactive))
		actions = append(actions, "Send urgent reminder", "Schedule mandatory check-in")
	} else if daysInactive >= 14 {
		upgrade("high")
		factors = append(factors, fmt.Sprintf("%d days inactive", daysInactive))
		actions = append(actions, "Send reminder", "Schedule check-in")
	} else if daysInactive >= 7 {
		upgrade("medium")
		factors = append(factors, fmt.Sprintf("%d days inactive", daysInactive))
		actions = append(actions, "Send reminder")
	}

	if total > 0 && completed < total/2 && progress < 60 {
		if progress < 25 {
			upgrade("high")
			factors = append(factors, "Very low progress")
		} else if progress < 40 {
			upgrade("medium")
			factors = append(factors, "Low progress")
		}
		actions = append(actions, "Review workload", "Offer support session")
	}

	if avgGrade > 0 && avgGrade < 70 {
		if avgGrade < 60 {
			upgrade("high")
			factors = append(factors, "Low average grade (<60%)")
		} else {
			upgrade("medium")
			factors = append(factors, "Average grade below target")
		}
		actions = append(actions, "Provide feedback", "Recommend extra practice")
	}

	// De-duplicate actions
	seen := map[string]bool{}
	var uniqActions []string
	for _, a := range actions {
		if !seen[a] {
			seen[a] = true
			uniqActions = append(uniqActions, a)
		}
	}

	return riskLevel, factors, uniqActions
}
