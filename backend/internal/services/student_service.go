package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/dto"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	pb "github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx/types"
)

type studentUserRepo interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type studentJourneyRepo interface {
	GetJourneyState(ctx context.Context, userID, tenantID string) (map[string]string, error)
}

type studentLMSRepo interface {
	GetStudentEnrollments(ctx context.Context, studentID string) ([]models.CourseEnrollment, error)
	CreateSubmission(ctx context.Context, sub *models.ActivitySubmission) error
	GetSubmissionByStudent(ctx context.Context, activityID, studentID string) (*models.ActivitySubmission, error)
}

type studentSchedulerRepo interface {
	GetOffering(ctx context.Context, id string) (*models.CourseOffering, error)
	ListStaff(ctx context.Context, offeringID string) ([]models.CourseStaff, error)
	ListSessions(ctx context.Context, offeringID string, startDate, endDate time.Time) ([]models.ClassSession, error)
}

type studentCurriculumRepo interface {
	GetCourse(ctx context.Context, id string) (*models.Course, error)
	ListCourses(ctx context.Context, tenantID string, programID *string) ([]models.Course, error)
}

type studentGradingRepo interface {
	ListStudentEntries(ctx context.Context, studentID string) ([]models.GradebookEntry, error)
}

type studentCourseContentRepo interface {
	// Read-only operations needed for the student portal.
	GetModule(ctx context.Context, id string) (*models.CourseModule, error)
	GetLesson(ctx context.Context, id string) (*models.CourseLesson, error)
	GetActivity(ctx context.Context, id string) (*models.CourseActivity, error)
	ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error)
	ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error)
	ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error)
}

type studentForumRepo interface {
	ListForums(ctx context.Context, courseOfferingID string) ([]models.Forum, error)
	ListTopics(ctx context.Context, forumID string, limit, offset int) ([]models.Topic, error)
}

type studentAttendanceRepo interface {
	RecordAttendance(ctx context.Context, sessionID string, record models.ClassAttendance) error
}

type StudentService struct {
	userRepo      studentUserRepo
	journeyRepo   studentJourneyRepo
	lmsRepo       studentLMSRepo
	schedulerRepo studentSchedulerRepo
	currRepo      studentCurriculumRepo
	gradingRepo   studentGradingRepo
	contentRepo   studentCourseContentRepo
	forumRepo     studentForumRepo
	attendanceRepo studentAttendanceRepo
	pb            *pb.Manager
}

func NewStudentService(
	userRepo studentUserRepo,
	journeyRepo studentJourneyRepo,
	lmsRepo studentLMSRepo,
	schedulerRepo studentSchedulerRepo,
	currRepo studentCurriculumRepo,
	gradingRepo studentGradingRepo,
	contentRepo studentCourseContentRepo,
	forumRepo studentForumRepo,
	attendanceRepo studentAttendanceRepo,
	pbm *pb.Manager,
) *StudentService {
	return &StudentService{
		userRepo:      userRepo,
		journeyRepo:   journeyRepo,
		lmsRepo:       lmsRepo,
		schedulerRepo: schedulerRepo,
		currRepo:      currRepo,
		gradingRepo:   gradingRepo,
		contentRepo:   contentRepo,
		forumRepo:     forumRepo,
		attendanceRepo: attendanceRepo,
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

func (s *StudentService) GetCourseDetail(ctx context.Context, tenantID, studentID, courseOfferingID string) (*dto.StudentCourseDetail, error) {
	if courseOfferingID == "" {
		return nil, errors.New("course_offering_id is required")
	}
	if err := s.ensureStudentEnrolled(ctx, studentID, courseOfferingID); err != nil {
		return nil, err
	}

	offering, err := s.schedulerRepo.GetOffering(ctx, courseOfferingID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, errors.New("course offering not found")
	}

	course, err := s.currRepo.GetCourse(ctx, offering.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	instructorName := s.resolvePrimaryInstructorName(ctx, offering.ID)

	now := time.Now()
	windowEnd := now.Add(60 * 24 * time.Hour)
	sessions, err := s.schedulerRepo.ListSessions(ctx, offering.ID, now.Add(-24*time.Hour), windowEnd)
	if err != nil {
		return nil, err
	}

	outSessions := make([]dto.StudentCourseNextSession, 0, len(sessions))
	for _, sess := range sessions {
		if sess.IsCancelled {
			continue
		}
		outSessions = append(outSessions, dto.StudentCourseNextSession{
			ID:         sess.ID,
			Date:       sess.Date.Format("2006-01-02"),
			StartTime:  sess.StartTime,
			EndTime:    sess.EndTime,
			RoomID:     sess.RoomID,
			MeetingURL: sess.MeetingURL,
			Type:       sess.Type,
		})
	}

	c := dto.StudentCourse{
		EnrollmentID:     "",
		CourseOfferingID: offering.ID,
		Status:           models.EnrollmentStatusEnrolled,
		CourseID:         offering.CourseID,
		Code:             course.Code,
		Title:            pickLocalizedTitle(course.Title, "en"),
		Section:          offering.Section,
		TermID:           offering.TermID,
		DeliveryFormat:   offering.DeliveryFormat,
		InstructorName:   instructorName,
		ProgressPercent:  0,
		NextSession:      nil,
	}
	if len(outSessions) > 0 {
		next := outSessions[0]
		c.NextSession = &next
	}

	return &dto.StudentCourseDetail{
		Course:   c,
		Sessions: outSessions,
	}, nil
}

func (s *StudentService) GetCourseModules(ctx context.Context, tenantID, studentID, courseOfferingID string) ([]models.CourseModule, error) {
	if courseOfferingID == "" {
		return nil, errors.New("course_offering_id is required")
	}
	if s.contentRepo == nil {
		return nil, errors.New("course content repository not configured")
	}
	if err := s.ensureStudentEnrolled(ctx, studentID, courseOfferingID); err != nil {
		return nil, err
	}

	offering, err := s.schedulerRepo.GetOffering(ctx, courseOfferingID)
	if err != nil {
		return nil, err
	}
	if offering == nil {
		return nil, errors.New("course offering not found")
	}

	modules, err := s.contentRepo.ListModules(ctx, offering.CourseID)
	if err != nil {
		return nil, err
	}

	for mi := range modules {
		lessons, err := s.contentRepo.ListLessons(ctx, modules[mi].ID)
		if err != nil {
			return nil, err
		}

		for li := range lessons {
			acts, err := s.contentRepo.ListActivities(ctx, lessons[li].ID)
			if err != nil {
				return nil, err
			}
			lessons[li].Activities = acts
		}

		modules[mi].Lessons = lessons
	}

	return modules, nil
}

func (s *StudentService) ListCourseAnnouncements(ctx context.Context, tenantID, studentID, courseOfferingID string) ([]dto.StudentAnnouncement, error) {
	if courseOfferingID == "" {
		return nil, errors.New("course_offering_id is required")
	}
	if s.forumRepo == nil {
		return nil, errors.New("forum repository not configured")
	}
	if err := s.ensureStudentEnrolled(ctx, studentID, courseOfferingID); err != nil {
		return nil, err
	}

	forums, err := s.forumRepo.ListForums(ctx, courseOfferingID)
	if err != nil {
		return nil, err
	}

	var annForumID string
	for _, f := range forums {
		if f.Type == models.ForumTypeAnnouncement {
			annForumID = f.ID
			break
		}
	}
	if annForumID == "" {
		return []dto.StudentAnnouncement{}, nil
	}

	topics, err := s.forumRepo.ListTopics(ctx, annForumID, 20, 0)
	if err != nil {
		return nil, err
	}

	out := make([]dto.StudentAnnouncement, 0, len(topics))
	for _, t := range topics {
		out = append(out, dto.StudentAnnouncement{
			ID:      t.ID,
			Title:   t.Title,
			Body:    t.Content,
			Created: t.CreatedAt.Format(time.RFC3339),
			Link:    nil,
		})
	}
	return out, nil
}

func (s *StudentService) ListCourseResources(ctx context.Context, tenantID, studentID, courseOfferingID string) ([]models.CourseActivity, error) {
	modules, err := s.GetCourseModules(ctx, tenantID, studentID, courseOfferingID)
	if err != nil {
		return nil, err
	}

	var resources []models.CourseActivity
	for _, m := range modules {
		for _, l := range m.Lessons {
			for _, a := range l.Activities {
				if strings.EqualFold(a.Type, "resource") {
					resources = append(resources, a)
				}
			}
		}
	}
	return resources, nil
}

func (s *StudentService) GetAssignmentDetail(ctx context.Context, tenantID, studentID, activityID, courseOfferingID string) (*models.CourseActivity, *models.ActivitySubmission, string, error) {
	if activityID == "" {
		return nil, nil, "", errors.New("activity_id is required")
	}
	if s.contentRepo == nil {
		return nil, nil, "", errors.New("course content repository not configured")
	}

	resolvedOfferingID, err := s.resolveStudentOfferingForActivity(ctx, studentID, activityID, courseOfferingID)
	if err != nil {
		return nil, nil, "", err
	}

	activity, err := s.contentRepo.GetActivity(ctx, activityID)
	if err != nil {
		return nil, nil, "", err
	}
	if activity == nil {
		return nil, nil, "", errors.New("activity not found")
	}

	var sub *models.ActivitySubmission
	if s.lmsRepo != nil {
		found, err := s.lmsRepo.GetSubmissionByStudent(ctx, activityID, studentID)
		if err == nil {
			sub = found
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, nil, "", err
		}
	}

	return activity, sub, resolvedOfferingID, nil
}

func (s *StudentService) SubmitAssignment(ctx context.Context, tenantID, studentID, activityID, courseOfferingID string, content json.RawMessage, status string) (*models.ActivitySubmission, error) {
	if activityID == "" {
		return nil, errors.New("activity_id is required")
	}
	if status == "" {
		status = "SUBMITTED"
	}

	resolvedOfferingID, err := s.resolveStudentOfferingForActivity(ctx, studentID, activityID, courseOfferingID)
	if err != nil {
		return nil, err
	}

	sub := &models.ActivitySubmission{
		ActivityID:       activityID,
		StudentID:        studentID,
		CourseOfferingID: resolvedOfferingID,
		Content:          types.JSONText(content),
		Status:           status,
	}
	if err := s.lmsRepo.CreateSubmission(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *StudentService) ensureStudentEnrolled(ctx context.Context, studentID, courseOfferingID string) error {
	enrollments, err := s.lmsRepo.GetStudentEnrollments(ctx, studentID)
	if err != nil {
		return err
	}
	for _, e := range enrollments {
		if e.CourseOfferingID == courseOfferingID {
			return nil
		}
	}
	return ErrForbidden
}

func (s *StudentService) resolveStudentOfferingForActivity(ctx context.Context, studentID, activityID, requestedOfferingID string) (string, error) {
	if requestedOfferingID != "" {
		if err := s.ensureStudentEnrolled(ctx, studentID, requestedOfferingID); err != nil {
			return "", err
		}
		if err := s.ensureActivityBelongsToOffering(ctx, activityID, requestedOfferingID); err != nil {
			return "", err
		}
		return requestedOfferingID, nil
	}

	activity, err := s.contentRepo.GetActivity(ctx, activityID)
	if err != nil {
		return "", err
	}
	if activity == nil {
		return "", errors.New("activity not found")
	}
	lesson, err := s.contentRepo.GetLesson(ctx, activity.LessonID)
	if err != nil {
		return "", err
	}
	if lesson == nil {
		return "", errors.New("lesson not found")
	}
	module, err := s.contentRepo.GetModule(ctx, lesson.ModuleID)
	if err != nil {
		return "", err
	}
	if module == nil {
		return "", errors.New("module not found")
	}

	enrollments, err := s.lmsRepo.GetStudentEnrollments(ctx, studentID)
	if err != nil {
		return "", err
	}

	var matches []string
	for _, e := range enrollments {
		offering, err := s.schedulerRepo.GetOffering(ctx, e.CourseOfferingID)
		if err != nil || offering == nil {
			continue
		}
		if offering.CourseID == module.CourseID {
			matches = append(matches, offering.ID)
		}
	}

	if len(matches) == 0 {
		return "", ErrForbidden
	}
	if len(matches) > 1 {
		return "", errors.New("multiple course offerings match this activity; course_offering_id is required")
	}
	return matches[0], nil
}

func (s *StudentService) ensureActivityBelongsToOffering(ctx context.Context, activityID, offeringID string) error {
	offering, err := s.schedulerRepo.GetOffering(ctx, offeringID)
	if err != nil {
		return err
	}
	if offering == nil {
		return errors.New("course offering not found")
	}

	activity, err := s.contentRepo.GetActivity(ctx, activityID)
	if err != nil {
		return err
	}
	if activity == nil {
		return errors.New("activity not found")
	}
	lesson, err := s.contentRepo.GetLesson(ctx, activity.LessonID)
	if err != nil {
		return err
	}
	if lesson == nil {
		return errors.New("lesson not found")
	}
	module, err := s.contentRepo.GetModule(ctx, lesson.ModuleID)
	if err != nil {
		return err
	}
	if module == nil {
		return errors.New("module not found")
	}
	if module.CourseID != offering.CourseID {
		return errors.New("activity does not belong to this course offering")
	}
	return nil
}

func (s *StudentService) ListAvailableCourses(ctx context.Context, tenantID string) ([]models.Course, error) {
	courses, err := s.currRepo.ListCourses(ctx, tenantID, nil)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

// CheckIn allows a student to self-check-in to a session via code or QR
func (s *StudentService) CheckIn(ctx context.Context, tenantID, userID, sessionID, code string) error {
	// 1. In a real app, verify code matches session code.
	// For demo, we trust the student scanning the QR code (which contains sessionID).
    
    if s.attendanceRepo == nil {
        return errors.New("attendance repo not configured")
    }

	// Double check session exists? Or just record.
	// We'll just record.

    rec := models.ClassAttendance{
        ClassSessionID: sessionID,
        StudentID:      userID,
        Status:         "PRESENT",
        Notes:          "Self check-in via QR",
        RecordedByID:   userID, // Check-in by self
    }

	return s.attendanceRepo.RecordAttendance(ctx, sessionID, rec)
}
