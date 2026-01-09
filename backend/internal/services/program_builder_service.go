package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/models"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
)

type ProgramBuilderService struct {
	repo repository.CurriculumRepository
}

func NewProgramBuilderService(repo repository.CurriculumRepository) *ProgramBuilderService {
	return &ProgramBuilderService{repo: repo}
}

// EnsureDraftMap ensures a Program Version exists for the program.
// "Journey Map" is a UI term; in storage this is a versioned program template.
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
		Title:     `{"en": "Draft Program Version"}`,
		Version:   "0.0.1",
		Config:    `{"phases":[]}`,
		IsActive:  true,
	}
	if err := s.repo.CreateJourneyMap(ctx, newMap); err != nil {
		return nil, err
	}
	return newMap, nil
}

type BuilderNode struct {
	ID              string          `json:"id"`
	ProgramVersionID string         `json:"program_version_id"`
	// JourneyMapID is a deprecated alias kept for backward compatibility with older frontends.
	JourneyMapID    string          `json:"journey_map_id,omitempty"`
	ParentNodeID *string         `json:"parent_node_id,omitempty"`
	Slug         string          `json:"slug"`
	Type         string          `json:"type"`
	Title        json.RawMessage `json:"title"`
	Description  json.RawMessage `json:"description"`
	ModuleKey    string          `json:"module_key"`
	Coordinates  json.RawMessage `json:"coordinates"`
	Config       json.RawMessage `json:"config"`
	Prerequisites []string       `json:"prerequisites"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type UpdateMapInput struct {
	Title    *json.RawMessage `json:"title,omitempty"`
	Version  *string         `json:"version,omitempty"`
	Config   *json.RawMessage `json:"config,omitempty"`
	Phases   *json.RawMessage `json:"phases,omitempty"`
	IsActive *bool           `json:"is_active,omitempty"`
}

// FullJourneyMap represents the complete map structure for the Builder UI
type FullJourneyMap struct {
	ID        string          `json:"id"`
	ProgramID string          `json:"program_id"`
	Title     json.RawMessage `json:"title"`
	Version   string          `json:"version"`
	Phases    []interface{}   `json:"phases"` // from config
	Nodes     []BuilderNode   `json:"nodes"`
	Edges     []interface{}   `json:"edges"` // derived or stored
}

func (s *ProgramBuilderService) GetJourneyMap(ctx context.Context, programID string) (*FullJourneyMap, error) {
	// 1. Get Map
	jm, err := s.repo.GetJourneyMapByProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	if jm == nil {
		// Auto-create draft if missing?
		jm, err = s.EnsureDraftMap(ctx, programID)
		if err != nil {
			return nil, err
		}
	}

	// 2. Get Nodes
	rawNodes, err := s.repo.GetNodeDefinitions(ctx, jm.ID)
	if err != nil {
		return nil, err
	}

	// 3. Parse Config for Phases
	var config map[string]interface{}
	phases := []interface{}{}
	if jm.Config != "" {
		if err := json.Unmarshal([]byte(jm.Config), &config); err == nil {
			if p, ok := config["phases"].([]interface{}); ok {
				phases = p
			}
		}
	}

	nodes := make([]BuilderNode, 0, len(rawNodes))
	for _, n := range rawNodes {
		nodes = append(nodes, toBuilderNode(n))
	}

	// 4. Construct Response
	return &FullJourneyMap{
		ID:        jm.ID,
		ProgramID: jm.ProgramID,
		Title:     normalizeJSONValue(jm.Title),
		Version:   jm.Version,
		Phases:    phases,
		Nodes:     nodes,
		Edges:     []interface{}{}, // Edges are derived from nodes' prerequisites in current model usually, or we can compute them here
	}, nil
}

func (s *ProgramBuilderService) UpdateJourneyMap(ctx context.Context, programID string, in UpdateMapInput) (*FullJourneyMap, error) {
	jm, err := s.repo.GetJourneyMapByProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	if jm == nil {
		jm, err = s.EnsureDraftMap(ctx, programID)
		if err != nil {
			return nil, err
		}
	}

	if in.Title != nil {
		jm.Title = string(*in.Title)
	}
	if in.Version != nil {
		jm.Version = *in.Version
	}
	if in.IsActive != nil {
		jm.IsActive = *in.IsActive
	}

	configChanged := in.Config != nil || in.Phases != nil
	if configChanged {
		cfg := map[string]interface{}{}
		if strings.TrimSpace(jm.Config) != "" {
			_ = json.Unmarshal([]byte(jm.Config), &cfg)
		}

		if in.Config != nil {
			var next map[string]interface{}
			if err := json.Unmarshal(*in.Config, &next); err != nil {
				return nil, err
			}
			cfg = next
		}

		if in.Phases != nil {
			var phases interface{}
			if err := json.Unmarshal(*in.Phases, &phases); err != nil {
				return nil, err
			}
			cfg["phases"] = phases
		}

		b, err := json.Marshal(cfg)
		if err != nil {
			return nil, err
		}
		jm.Config = string(b)
	}

	if err := s.repo.UpdateJourneyMap(ctx, jm); err != nil {
		return nil, err
	}

	return s.GetJourneyMap(ctx, programID)
}

func (s *ProgramBuilderService) GetNodes(ctx context.Context, programID string) ([]BuilderNode, error) {
	jm, err := s.repo.GetJourneyMapByProgram(ctx, programID)
	if err != nil {
		return nil, err
	}
	if jm == nil {
		return []BuilderNode{}, nil
	}
	rawNodes, err := s.repo.GetNodeDefinitions(ctx, jm.ID)
	if err != nil {
		return nil, err
	}
	nodes := make([]BuilderNode, 0, len(rawNodes))
	for _, n := range rawNodes {
		nodes = append(nodes, toBuilderNode(n))
	}
	return nodes, nil
}

func (s *ProgramBuilderService) GetNode(ctx context.Context, nodeID string) (*models.JourneyNodeDefinition, error) {
	return s.repo.GetNodeDefinition(ctx, nodeID)
}

func (s *ProgramBuilderService) GetBuilderNode(ctx context.Context, nodeID string) (*BuilderNode, error) {
	node, err := s.repo.GetNodeDefinition(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	bn := toBuilderNode(*node)
	return &bn, nil
}

func (s *ProgramBuilderService) CreateNode(ctx context.Context, journeyMapID string, nodeDef *models.JourneyNodeDefinition) error {
	if nodeDef == nil {
		return errors.New("node is required")
	}
	nodeDef.JourneyMapID = journeyMapID
	if strings.TrimSpace(nodeDef.Title) == "" {
		titleBytes, _ := json.Marshal("Untitled")
		nodeDef.Title = string(titleBytes)
	}
	if nodeDef.Description == nil || strings.TrimSpace(*nodeDef.Description) == "" {
		s := "null"
		nodeDef.Description = &s
	}
	if strings.TrimSpace(nodeDef.Coordinates) == "" {
		nodeDef.Coordinates = `{"x":0,"y":0}`
	}
	if strings.TrimSpace(nodeDef.Config) == "" {
		nodeDef.Config = "{}"
	}
	return s.repo.CreateNodeDefinition(ctx, nodeDef)
}

func (s *ProgramBuilderService) UpdateNode(ctx context.Context, nodeDef *models.JourneyNodeDefinition) error {
	if nodeDef == nil {
		return errors.New("node is required")
	}
	if strings.TrimSpace(nodeDef.Title) == "" {
		nodeDef.Title = "null"
	}
	if nodeDef.Description == nil || strings.TrimSpace(*nodeDef.Description) == "" {
		s := "null"
		nodeDef.Description = &s
	}
	if strings.TrimSpace(nodeDef.Coordinates) == "" {
		nodeDef.Coordinates = `{"x":0,"y":0}`
	}
	if strings.TrimSpace(nodeDef.Config) == "" {
		nodeDef.Config = "{}"
	}
	return s.repo.UpdateNodeDefinition(ctx, nodeDef)
}

func normalizeJSONValue(value string) json.RawMessage {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	if json.Valid([]byte(trimmed)) {
		return json.RawMessage(trimmed)
	}
	b, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return json.RawMessage(b)
}

func normalizeJSONObject(value string, defaultJSON string) json.RawMessage {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return json.RawMessage(defaultJSON)
	}
	if json.Valid([]byte(trimmed)) && strings.HasPrefix(trimmed, "{") {
		return json.RawMessage(trimmed)
	}
	return json.RawMessage(defaultJSON)
}

func toBuilderNode(n models.JourneyNodeDefinition) BuilderNode {
	desc := ""
	if n.Description != nil {
		desc = *n.Description
	}
	return BuilderNode{
		ID:              n.ID,
		ProgramVersionID: n.JourneyMapID,
		JourneyMapID:     n.JourneyMapID,
		ParentNodeID:  n.ParentNodeID,
		Slug:          n.Slug,
		Type:          n.Type,
		Title:         normalizeJSONValue(n.Title),
		Description:   normalizeJSONValue(desc),
		ModuleKey:     n.ModuleKey,
		Coordinates:   normalizeJSONObject(n.Coordinates, `{"x":0,"y":0}`),
		Config:        normalizeJSONObject(n.Config, `{}`),
		Prerequisites: []string(n.Prerequisites),
		CreatedAt:     n.CreatedAt,
		UpdatedAt:     n.UpdatedAt,
	}
}
