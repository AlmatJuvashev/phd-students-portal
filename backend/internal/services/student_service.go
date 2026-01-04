package services

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/dto"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
)

type studentUserRepo interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type studentJourneyRepo interface {
	GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error)
}

type studentLMSRepo interface {
	GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error)
}

type studentSchedulerRepo interface {
	GetOffering(ctx context.Context, id string) (*models.CourseOffering, error)
	ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error)
	ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error)
}

type studentCurriculumRepo interface {
	GetCourse(ctx context.Context, id string) (*models.Course, error)
}

type studentGradingRepo interface {
	ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error)
}

type StudentService struct {
	userRepo      studentUserRepo
	journeyRepo   studentJourneyRepo
	lmsRepo       studentLMSRepo
	schedulerRepo studentSchedulerRepo
	currRepo      studentCurriculumRepo
	gradingRepo   studentGradingRepo
	pb            *pb.Manager
}

func NewStudentService(
	userRepo studentUserRepo,
	journeyRepo studentJourneyRepo,
	lmsRepo studentLMSRepo,
	schedulerRepo studentSchedulerRepo,
	currRepo studentCurriculumRepo,
	gradingRepo studentGradingRepo,
	pbm *pb.Manager,
) *StudentService {
	return &StudentService{
		userRepo:      userRepo,
		journeyRepo:   journeyRepo,
		lmsRepo:       lmsRepo,
		schedulerRepo: schedulerRepo,
		currRepo:      currRepo,
		gradingRepo:   gradingRepo,
		pb:            pbm,
	}
}

func (s *StudentService) GetDashboard(ctx context.Context, tenantID, studentID string) (*dto.StudentDashboard, error) {
	user, err := s.userRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	state, err := s.journeyRepo.GetJourneyState(ctx, studentID, tenantID)
	if err != nil {
		return nil, err
	}

	progress := s.computeProgramProgress(user, state)
	upcoming := s.suggestUpcomingJourneyDeadlines(state, 6)

	grades, err := s.ListGrades(ctx, tenantID, studentID)
	if err != nil {
		return nil, err
	}
	if len(grades) > 5 {
		grades = grades[:5]
	}

	return &dto.StudentDashboard{
		Program:           progress,
		UpcomingDeadlines: upcoming,
		RecentGrades:      grades,
		Announcements:     []dto.StudentAnnouncement{},
	}, nil
}

func (s *StudentService) ListCourses(ctx context.Context, tenantID, studentID string) ([]dto.StudentCourse, error) {
	enrollments, err := s.lmsRepo.GetStudentEnrollments(ctx, studentID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	windowEnd := now.Add(14 * 24 * time.Hour)

	var out []dto.StudentCourse
	for _, e := range enrollments {
		offering, err := s.schedulerRepo.GetOffering(ctx, e.CourseOfferingID)
		if err != nil || offering == nil {
			continue
		}

		course, err := s.currRepo.GetCourse(ctx, offering.CourseID)
		if err != nil || course == nil {
			continue
		}

		instructorName := s.resolvePrimaryInstructorName(ctx, offering.ID)

		var nextSession *dto.StudentCourseNextSession
		sessions, err := s.schedulerRepo.ListSessions(ctx, offering.ID, now, windowEnd)
		if err == nil && len(sessions) > 0 {
			for _, sess := range sessions {
				if sess.IsCancelled {
					continue
				}
				if sess.Date.Before(now.Add(-24 * time.Hour)) {
					continue
				}
				nextSession = &dto.StudentCourseNextSession{
					ID:         sess.ID,
					Date:       sess.Date.Format("2006-01-02"),
					StartTime:  sess.StartTime,
					EndTime:    sess.EndTime,
					RoomID:     sess.RoomID,
					MeetingURL: sess.MeetingURL,
					Type:       sess.Type,
				}
				break
			}
		}

		out = append(out, dto.StudentCourse{
			EnrollmentID:     e.ID,
			CourseOfferingID: offering.ID,
			Status:           e.Status,
			CourseID:         offering.CourseID,
			Code:             course.Code,
			Title:            pickLocalizedTitle(course.Title, "en"),
			Section:          offering.Section,
			TermID:           offering.TermID,
			DeliveryFormat:   offering.DeliveryFormat,
			InstructorName:   instructorName,
			ProgressPercent:  0,
			NextSession:      nextSession,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Code) < strings.ToLower(out[j].Code)
	})

	return out, nil
}

func (s *StudentService) ListAssignments(ctx context.Context, tenantID, studentID string) ([]dto.StudentAssignment, error) {
	state, err := s.journeyRepo.GetJourneyState(ctx, studentID, tenantID)
	if err != nil {
		return nil, err
	}

	deadlines := s.suggestUpcomingJourneyDeadlines(state, 20)
	assignments := make([]dto.StudentAssignment, 0, len(deadlines))
	for _, d := range deadlines {
		assignments = append(assignments, dto.StudentAssignment{
			ID:       d.ID,
			Title:    d.Title,
			Source:   d.Source,
			Status:   d.Status,
			DueAt:    d.DueAt,
			Link:     d.Link,
			Severity: d.Severity,
		})
	}
	return assignments, nil
}

func (s *StudentService) ListGrades(ctx context.Context, tenantID, studentID string) ([]dto.StudentGradeEntry, error) {
	entries, err := s.gradingRepo.ListStudentEntries(ctx, studentID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.StudentGradeEntry, 0, len(entries))
	for _, e := range entries {
		dtoEntry := dto.StudentGradeEntry{
			ID:               e.ID,
			CourseOfferingID: e.CourseOfferingID,
			ActivityID:       e.ActivityID,
			StudentID:        e.StudentID,
			Score:            e.Score,
			MaxScore:         e.MaxScore,
			Grade:            e.Grade,
			Feedback:         e.Feedback,
			GradedByID:       e.GradedByID,
			GradedAt:         e.GradedAt.Format(time.RFC3339),
		}

		offering, err := s.schedulerRepo.GetOffering(ctx, e.CourseOfferingID)
		if err == nil && offering != nil {
			dtoEntry.CourseID = offering.CourseID
			course, err := s.currRepo.GetCourse(ctx, offering.CourseID)
			if err == nil && course != nil {
				dtoEntry.CourseCode = &course.Code
				title := pickLocalizedTitle(course.Title, "en")
				dtoEntry.CourseTitle = &title
			}
		}

		out = append(out, dtoEntry)
	}

	return out, nil
}

func (s *StudentService) computeProgramProgress(user *models.User, state map[string]string) dto.ProgramProgress {
	title := "Program"
	if user != nil && strings.TrimSpace(user.Program) != "" {
		title = user.Program
	}
	total := 0
	done := 0
	if s.pb != nil {
		total = len(s.pb.Nodes)
		for nodeID := range s.pb.Nodes {
			if state[nodeID] == "done" {
				done++
			}
		}
	}
	var pct float64
	if total > 0 {
		pct = float64(done) * 100.0 / float64(total)
	}
	return dto.ProgramProgress{
		Title:           title,
		ProgressPercent: pct,
		CompletedNodes:  done,
		TotalNodes:      total,
		OverdueCount:    0,
	}
}

func (s *StudentService) suggestUpcomingJourneyDeadlines(state map[string]string, limit int) []dto.StudentDeadline {
	if s.pb == nil || limit <= 0 {
		return []dto.StudentDeadline{}
	}

	type candidate struct {
		id       string
		title    string
		world    string
		worldOrd int
		severity string
		status   string
	}

	done := func(nodeID string) bool { return state[nodeID] == "done" }
	prereqsMet := func(n pb.Node) bool {
		for _, p := range n.Prerequisites {
			if !done(p) {
				return false
			}
		}
		return true
	}

	var list []candidate
	for _, n := range s.pb.Nodes {
		status := state[n.ID]
		if status == "done" {
			continue
		}
		if !prereqsMet(n) {
			continue
		}

		world := s.pb.NodeWorldID(n.ID)
		ord := parseWorldOrder(world)

		severity := "normal"
		if status == "needs_fixes" || status == "submitted" {
			severity = "urgent"
		}

		list = append(list, candidate{
			id:       n.ID,
			title:    pickTitleFromMap(n.Title, s.pb.DefaultLocale),
			world:    world,
			worldOrd: ord,
			severity: severity,
			status:   status,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].severity != list[j].severity {
			return list[i].severity == "urgent"
		}
		if list[i].worldOrd != list[j].worldOrd {
			return list[i].worldOrd < list[j].worldOrd
		}
		return list[i].id < list[j].id
	})

	if len(list) > limit {
		list = list[:limit]
	}

	out := make([]dto.StudentDeadline, 0, len(list))
	for _, c := range list {
		link := "/journey"
		out = append(out, dto.StudentDeadline{
			ID:       c.id,
			Title:    c.title,
			DueAt:    nil,
			Source:   "journey",
			Status:   c.status,
			Severity: c.severity,
			Link:     &link,
		})
	}
	return out
}

func parseWorldOrder(world string) int {
	// "W1" -> 1, "W10" -> 10; unknowns go to end.
	if len(world) < 2 {
		return 999
	}
	if world[0] != 'W' {
		return 999
	}
	n, err := strconv.Atoi(world[1:])
	if err != nil {
		return 999
	}
	return n
}

func pickTitleFromMap(title map[string]string, locale string) string {
	if len(title) == 0 {
		return "Untitled"
	}
	if locale != "" {
		if v := strings.TrimSpace(title[locale]); v != "" {
			return v
		}
	}
	if v := strings.TrimSpace(title["en"]); v != "" {
		return v
	}
	if v := strings.TrimSpace(title["ru"]); v != "" {
		return v
	}
	if v := strings.TrimSpace(title["kz"]); v != "" {
		return v
	}
	for _, v := range title {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return "Untitled"
}

func pickLocalizedTitle(titleJSON string, locale string) string {
	// Course.Title is stored as JSONB string in DB and comes through as raw json string.
	// For now, we treat it as an opaque string if it doesn't look like JSON.
	trimmed := strings.TrimSpace(titleJSON)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "{") {
		// Best-effort parse without failing the request.
		m := map[string]string{}
		_ = json.Unmarshal([]byte(trimmed), &m)
		return pickTitleFromMap(m, locale)
	}
	return trimmed
}

func (s *StudentService) resolvePrimaryInstructorName(ctx context.Context, offeringID string) *string {
	staff, err := s.schedulerRepo.ListStaff(ctx, offeringID)
	if err != nil || len(staff) == 0 {
		return nil
	}

	var primary *models.CourseStaff
	for _, st := range staff {
		if strings.EqualFold(st.Role, "INSTRUCTOR") && st.IsPrimary {
			val := st
			primary = &val
			break
		}
	}
	if primary == nil {
		for _, st := range staff {
			if strings.EqualFold(st.Role, "INSTRUCTOR") {
				val := st
				primary = &val
				break
			}
		}
	}
	if primary == nil {
		return nil
	}

	u, err := s.userRepo.GetByID(ctx, primary.UserID)
	if err != nil || u == nil {
		return nil
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	if name == "" {
		name = u.Email
	}
	if name == "" {
		return nil
	}
	return &name
}
