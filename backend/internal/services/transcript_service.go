package services

import (
	"context"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type TranscriptService struct {
	transcriptRepo repository.TranscriptRepository
	schedulerRepo  repository.SchedulerRepository // To get Term details (Name)
}

func NewTranscriptService(tr repository.TranscriptRepository, sr repository.SchedulerRepository) *TranscriptService {
	return &TranscriptService{
		transcriptRepo: tr,
		schedulerRepo:  sr,
	}
}

func (s *TranscriptService) GetTranscript(ctx context.Context, studentID string) (*models.Transcript, error) {
	// 1. Fetch all grades
	grades, err := s.transcriptRepo.GetStudentGrades(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch all terms (or just unique ones from grades, but fetching all is safer for ordering)
	// Actually, we can just fetch terms that appear in grades.
	// But we need Term Names.
	// Let's assume SchedulerRepo has GetTerm(id).
	// Optimization: We could join in Repo, but let's stick to simple composition for now.
	
	// Map to hold term info
	termMap := make(map[string]*models.AcademicTerm)
	
	// 3. Aggregate
	transcript := &models.Transcript{
		StudentID: studentID,
		Terms:     []models.TranscriptTerm{},
	}
	
	var totalPoints float64
	var totalCredits float64
	
	// Group by Term
	gradesByTerm := make(map[string][]models.TermGrade)
	termOrder := []string{} // To preserve order from Repo (which ordered by start_date)
	
	for _, g := range grades {
		if _, exists := gradesByTerm[g.TermID]; !exists {
			termOrder = append(termOrder, g.TermID)
			// Fetch Term Info if missing
			if _, haveTerm := termMap[g.TermID]; !haveTerm {
				term, err := s.schedulerRepo.GetTerm(ctx, g.TermID)
				if err == nil {
					termMap[g.TermID] = term
				}
			}
		}
		gradesByTerm[g.TermID] = append(gradesByTerm[g.TermID], g)
		
		// Cumulative stats
		// Only count passed courses or all courses? GPA usually includes Fs (0 points).
		// We assume `GradePoints` is correctly populated (e.g. F=0).
		totalPoints += (g.GradePoints * g.Credits)
		// Usually GPA is calculated on "Attempted Credits", passing checks is for "Earned Credits".
		// For simplicity MVP: TotalCredits = Attempted.
		totalCredits += g.Credits
	}
	
	// Build TranscriptTerms
	for _, termID := range termOrder {
		termGrades := gradesByTerm[termID]
		
		var termPoints float64
		var termCredits float64
		
		for _, g := range termGrades {
			// Don't double count if same course? No, simplified logic.
			termPoints += (g.GradePoints * g.Credits)
			termCredits += g.Credits
		}
		
		termGPA := float32(0)
		if termCredits > 0 {
			termGPA = float32(termPoints / termCredits)
		}
		
		termName := "Unknown Term"
		if t, ok := termMap[termID]; ok {
			termName = t.Name
		}
		
		transcript.Terms = append(transcript.Terms, models.TranscriptTerm{
			TermID:      termID,
			TermName:    termName,
			TermGPA:     termGPA,
			TermCredits: termCredits,
			Grades:      termGrades,
		})
	}
	
	transcript.TotalCredits = totalCredits
	transcript.TotalPoints = totalPoints
	if totalCredits > 0 {
		transcript.CumulativeGPA = float32(totalPoints / totalCredits)
	}
	
	return transcript, nil
}
