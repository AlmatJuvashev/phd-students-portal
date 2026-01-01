package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type CourseContentService struct {
	repo repository.CourseContentRepository
}

func NewCourseContentService(repo repository.CourseContentRepository) *CourseContentService {
	return &CourseContentService{repo: repo}
}

// --- Modules ---

func (s *CourseContentService) CreateModule(ctx context.Context, m *models.CourseModule) error {
	if m.CourseID == "" {
		return errors.New("course_id is required")
	}
	if m.Title == "" {
		return errors.New("title is required")
	}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return s.repo.CreateModule(ctx, m)
}

func (s *CourseContentService) ListModules(ctx context.Context, courseID string) ([]models.CourseModule, error) {
	return s.repo.ListModules(ctx, courseID)
}

func (s *CourseContentService) UpdateModule(ctx context.Context, m *models.CourseModule) error {
	m.UpdatedAt = time.Now()
	return s.repo.UpdateModule(ctx, m)
}

func (s *CourseContentService) DeleteModule(ctx context.Context, id string) error {
	return s.repo.DeleteModule(ctx, id)
}

// --- Lessons ---

func (s *CourseContentService) CreateLesson(ctx context.Context, l *models.CourseLesson) error {
	if l.ModuleID == "" {
		return errors.New("module_id is required")
	}
	if l.Title == "" {
		return errors.New("title is required")
	}
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	return s.repo.CreateLesson(ctx, l)
}

func (s *CourseContentService) ListLessons(ctx context.Context, moduleID string) ([]models.CourseLesson, error) {
	return s.repo.ListLessons(ctx, moduleID)
}

func (s *CourseContentService) UpdateLesson(ctx context.Context, l *models.CourseLesson) error {
	l.UpdatedAt = time.Now()
	return s.repo.UpdateLesson(ctx, l)
}

func (s *CourseContentService) DeleteLesson(ctx context.Context, id string) error {
	return s.repo.DeleteLesson(ctx, id)
}

// --- Activities ---

func (s *CourseContentService) CreateActivity(ctx context.Context, a *models.CourseActivity) error {
	if a.LessonID == "" {
		return errors.New("lesson_id is required")
	}
	if a.Title == "" {
		return errors.New("title is required")
	}
	if a.Type == "" {
		return errors.New("type is required")
	}
	if err := s.validateActivityContent(a); err != nil {
		return err
	}
	if a.Content == "" {
		a.Content = "{}"
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return s.repo.CreateActivity(ctx, a)
}

func (s *CourseContentService) ListActivities(ctx context.Context, lessonID string) ([]models.CourseActivity, error) {
	return s.repo.ListActivities(ctx, lessonID)
}

func (s *CourseContentService) UpdateActivity(ctx context.Context, a *models.CourseActivity) error {
	if err := s.validateActivityContent(a); err != nil {
		return err
	}
	a.UpdatedAt = time.Now()
	return s.repo.UpdateActivity(ctx, a)
}

func (s *CourseContentService) DeleteActivity(ctx context.Context, id string) error {
	return s.repo.DeleteActivity(ctx, id)
}

func (s *CourseContentService) validateActivityContent(a *models.CourseActivity) error {
	if a.Content == "" || a.Content == "{}" {
		return nil // Allow empty content for drafts
	}
	
	switch a.Type {
	case "quiz":
		var qc models.QuizConfig
		if err := json.Unmarshal([]byte(a.Content), &qc); err != nil {
			return errors.New("invalid quiz config json")
		}
		// Basic validation: must have time limit if questions exist? 
		// For now, just ensuring it is valid JSON structure for Quiz
		
	case "survey":
		var sc models.SurveyConfig
		if err := json.Unmarshal([]byte(a.Content), &sc); err != nil {
			return errors.New("invalid survey config json")
		}
	}
	return nil
}
