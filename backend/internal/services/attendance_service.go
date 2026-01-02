package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type AttendanceService struct {
	attendanceRepo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) *AttendanceService {
	return &AttendanceService{attendanceRepo: repo}
}

// BatchRecordAttendance updates attendance for multiple students in a session
func (s *AttendanceService) BatchRecordAttendance(ctx context.Context, sessionID string, updates []models.ClassAttendance, recordedBy string) error {
	// Simple pass-through for now, but could add validation logic here (e.g. check if session exists)
	return s.attendanceRepo.BatchUpsertAttendance(ctx, sessionID, updates, recordedBy)
}
