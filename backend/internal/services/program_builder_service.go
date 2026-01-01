package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ProgramBuilderService struct {
	repo repository.CurriculumRepository
}

func NewProgramBuilderService(repo repository.CurriculumRepository) *ProgramBuilderService {
	return &ProgramBuilderService{repo: repo}
}

// EnsureDraftMap ensures a JourneyMap exists for the program.
// In a real builder, this might handle "Draft vs Published" versions.
// For now, it gets or creates the active map.
func (s *ProgramBuilderService) EnsureDraftMap(ctx context.Context, programID string) (*models.JourneyMap, error) {
	jm, err := s.repo.GetJourneyMapByProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	if jm != nil {
		return jm, nil
	}

	// Create new
	newMap := &models.JourneyMap{
		ProgramID: programID,
		Title:     `{"en": "Default Journey"}`,
		Version:   "0.0.1",
		IsActive:  true,
	}
	if err := s.repo.CreateJourneyMap(ctx, newMap); err != nil {
		return nil, err
	}
	return newMap, nil
}

// UpdateNodeConfig validates and updates a node's configuration.
// It accepts the strict ProgramNodeConfig, validates it against the node type,
// and saves it as JSONB.
func (s *ProgramBuilderService) UpdateNodeConfig(ctx context.Context, nodeID string, config models.ProgramNodeConfig) error {
	node, err := s.repo.GetNodeDefinition(ctx, nodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	// Validation Logic
	if err := s.validateConfig(node.Type, config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// Marshal to JSONB
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	node.Config = string(configBytes)

	return s.repo.UpdateNodeDefinition(ctx, node)
}

// CreateNode creates a new node in the journey map.
func (s *ProgramBuilderService) CreateNode(ctx context.Context, journeyMapID string, nodeDef models.JourneyNodeDefinition, config models.ProgramNodeConfig) (*models.JourneyNodeDefinition, error) {
	// Validate type and config
	if err := s.validateConfig(nodeDef.Type, config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	nodeDef.Config = string(configBytes)
	nodeDef.JourneyMapID = journeyMapID

	if err := s.repo.CreateNodeDefinition(ctx, &nodeDef); err != nil {
		return nil, err
	}
	return &nodeDef, nil
}

func (s *ProgramBuilderService) validateConfig(nodeType string, config models.ProgramNodeConfig) error {
	switch nodeType {
	case "formEntry":
		if len(config.Fields) == 0 {
			return errors.New("formEntry must have at least one field")
		}
		for _, f := range config.Fields {
			if f.Key == "" {
				return errors.New("field key is required")
			}
			if f.Type == "" {
				return errors.New("field type is required")
			}
			// Validate known types
			switch f.Type {
			case "text", "textarea", "boolean", "date", "file", "note", "select":
			default:
				return fmt.Errorf("unknown field type: %s", f.Type)
			}
		}
	case "checklist":
		if len(config.Fields) == 0 {
			return errors.New("checklist must have at least one item")
		}
	case "cards":
		if len(config.Slides) == 0 {
			return errors.New("cards must have at least one slide")
		}
	case "info":
		// Info nodes might just use Title/Description, config optional
	default:
		return fmt.Errorf("unknown node type: %s", nodeType)
	}
	return nil
}

func (s *ProgramBuilderService) GetNodes(ctx context.Context, programID string) ([]models.JourneyNodeDefinition, error) {
	jm, err := s.repo.GetJourneyMapByProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	if jm == nil {
		return []models.JourneyNodeDefinition{}, nil
	}
	return s.repo.GetNodeDefinitions(ctx, jm.ID)
}

func (s *ProgramBuilderService) GetNode(ctx context.Context, nodeID string) (*models.JourneyNodeDefinition, error) {
	return s.repo.GetNodeDefinition(ctx, nodeID)
}
