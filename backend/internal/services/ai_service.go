package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/config"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	openai "github.com/sashabaranov/go-openai"
)

type AIService struct {
	client *openai.Client
	enabled bool
}

func NewAIService(cfg config.AppConfig) *AIService {
	if cfg.OpenAIKey == "" {
		return &AIService{enabled: false}
	}
	return &AIService{
		client: openai.NewClient(cfg.OpenAIKey),
		enabled: true,
	}
}

// GenerateCourseStructure parses raw syllabus text into a structured CourseModule list
func (s *AIService) GenerateCourseStructure(ctx context.Context, syllabusText string) ([]models.CourseModule, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI service is not configured (missing OPENAI_API_KEY)")
	}

	systemPrompt := `You are an expert educational curriculum designer.
Your task is to parse the provided syllabus text (which might be unstructured) into a structured JSON format matching the schema below.
Extract modules, lessons, and activities.
If specific activities are not mentioned, infer reasonable default activities (e.g., "Reading", "Quiz") based on the lesson topic.

Output JSON Schema:
[
  {
    "title": "Module Title",
    "order": 1,
    "lessons": [
      {
        "title": "Lesson Title",
        "order": 1,
        "activities": [
           {
             "title": "Activity Title",
             "type": "text", // or "quiz", "video", "assignment"
             "points": 10,
             "is_optional": false
           }
        ]
      }
    ]
  }
]

Return ONLY valid JSON. No markdown formatting.
`

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: syllabusText,
				},
			},
			Temperature: 0.1, // Low temperature for deterministic structure
		},
	)

	if err != nil {
		return nil, fmt.Errorf("openai api error: %w", err)
	}

	content := resp.Choices[0].Message.Content
	// Strip markdown code blocks if present
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	
	var modules []models.CourseModule
	if err := json.Unmarshal([]byte(content), &modules); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w. Response: %s", err, content)
	}

	return modules, nil
}

// GenerateQuizConfig creates a quiz based on topic and difficulty
func (s *AIService) GenerateQuizConfig(ctx context.Context, topic string, difficulty string, count int) (*models.QuizConfig, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI service is not configured")
	}

	systemPrompt := fmt.Sprintf(`Generate a JSON quiz configuration for the topic "%s".
Difficulty: %s
Number of Questions: %d
Format: JSON matching schema:
{
  "timeLimit": 30,
  "passingScore": 70,
  "shuffleQuestions": true,
  "showResults": true,
  "questions": [
    {
      "id": "uuid-1",
      "type": "multiple_choice",
      "text": "Question text?",
      "points": 10,
      "options": [
        {"id": "opt-1", "text": "Option A", "isCorrect": true},
        {"id": "opt-2", "text": "Option B", "isCorrect": false}
      ]
    }
  ]
}`, topic, difficulty, count)

	resp, err := s.callAI(ctx, systemPrompt, topic)
	if err != nil {
		return nil, err
	}

	var config models.QuizConfig
	if err := json.Unmarshal([]byte(resp), &config); err != nil {
		return nil, fmt.Errorf("failed to parse quiz config: %w", err)
	}
	return &config, nil
}

// GenerateSurveyConfig creates a survey based on topic
func (s *AIService) GenerateSurveyConfig(ctx context.Context, topic string, count int) (*models.SurveyConfig, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI service is not configured")
	}

	systemPrompt := fmt.Sprintf(`Generate a JSON survey configuration for the topic "%s".
Number of Questions: %d
Format: JSON matching schema:
{
  "anonymous": true,
  "showProgressBar": true,
  "questions": [
    {
      "id": "uuid-1",
      "type": "rating_stars", // or text, choice
      "text": "Question text?",
      "required": true
    }
  ]
}`, topic, count)

	resp, err := s.callAI(ctx, systemPrompt, topic)
	if err != nil {
		return nil, err
	}

	var config models.SurveyConfig
	if err := json.Unmarshal([]byte(resp), &config); err != nil {
		return nil, fmt.Errorf("failed to parse survey config: %w", err)
	}
	return &config, nil
}

type GeneratedAssessmentItem struct {
	Type       string          `json:"type"`
	Difficulty int             `json:"difficulty"`
	Content    json.RawMessage `json:"content"`
	Tags       []string        `json:"tags"`
}

// GenerateAssessmentItems creates raw items for an Item Bank
func (s *AIService) GenerateAssessmentItems(ctx context.Context, topic string, itemType string, count int) ([]GeneratedAssessmentItem, error) {
	if !s.enabled {
		return nil, fmt.Errorf("AI service is not configured")
	}

	systemPrompt := fmt.Sprintf(`Generate %d assessment items for topic "%s".
Type: %s (multiple_choice, true_false, essay)
Format: JSON Array of objects:
[
  {
    "type": "multiple_choice",
    "difficulty": 3,
    "tags": ["tag1"],
    "content": { "text": "...", "options": [...], "answer": "..." }
  }
]`, count, topic, itemType)

	resp, err := s.callAI(ctx, systemPrompt, topic)
	if err != nil {
		return nil, err
	}

	var items []GeneratedAssessmentItem
	if err := json.Unmarshal([]byte(resp), &items); err != nil {
		return nil, fmt.Errorf("failed to parse assessment items: %w", err)
	}
	return items, nil
}

// Helper to handle AI call and basic cleanup
func (s *AIService) callAI(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: userMessage},
			},
			Temperature: 0.2,
		},
	)
	if err != nil {
		return "", err
	}
	content := resp.Choices[0].Message.Content
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return content, nil
}
