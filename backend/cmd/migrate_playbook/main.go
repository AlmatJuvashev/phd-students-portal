package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"
)

// Define structures matching playbook.json
type LocalizedString map[string]string

type Playbook struct {
    PlaybookID    string         `json:"playbook_id"`
    Version       string         `json:"version"`
    Worlds        []WorldDef     `json:"worlds"`
    Roles         []RoleDef      `json:"roles"`
    Conditions    []ConditionDef `json:"conditions"`
}

type WorldDef struct {
    ID    string          `json:"id"`
    Title LocalizedString `json:"title"`
    Order int             `json:"order"`
    Nodes []NodeDef       `json:"nodes"`
}

type RoleDef struct {
    ID    string          `json:"id"`
    Label LocalizedString `json:"label"`
}

type ConditionDef struct {
    ID          string          `json:"id"`
    Expr        string          `json:"expr"`
    Description LocalizedString `json:"description"`
}

type NodeDef struct {
    ID             string                 `json:"id"`
    Title          LocalizedString        `json:"title"`
    Type           string                 `json:"type"`
    WhoCanComplete []string               `json:"who_can_complete"`
    Prerequisites  []string               `json:"prerequisites"`
    Next           []string               `json:"next"`
    Outcomes       []OutcomeDef           `json:"outcomes"`
    Requirements   map[string]interface{} `json:"requirements"`
    Screen         map[string]interface{} `json:"screen"`
    ActionHints    []string               `json:"actionHints"`
    Condition      string                 `json:"condition"`
    Timer          interface{}            `json:"timer"`
}

type OutcomeDef struct {
    Value string          `json:"value"`
    Label LocalizedString `json:"label"`
    Next  []string        `json:"next"`
}

func main() {
    playbookPath := flag.String("playbook", "../frontend/src/playbooks/playbook.json", "Path to playbook.json")
    dsn := flag.String("dsn", "postgres://postgres:postgres@localhost:5435/phd?sslmode=disable", "Database connection string")
    flag.Parse()

    // 1. Read playbook.json
    absPath, _ := filepath.Abs(*playbookPath)
    log.Printf("Reading playbook from: %s", absPath)
    
    data, err := os.ReadFile(*playbookPath)
    if err != nil {
        log.Fatalf("Failed to read playbook file: %v", err)
    }

    var playbook Playbook
    if err := json.Unmarshal(data, &playbook); err != nil {
        log.Fatalf("Failed to unmarshal playbook JSON: %v", err)
    }

    log.Printf("Loaded playbook version %s with %d worlds", playbook.Version, len(playbook.Worlds))

    // 2. Connect to DB
    db, err := sql.Open("pgx", *dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    ctx := context.Background()
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to begin transaction: %v", err)
    }
    defer tx.Rollback()

    // 2.5 Ensure Tenant exists (default)
    var tenantID string
    err = tx.QueryRowContext(ctx, "SELECT id FROM tenants LIMIT 1").Scan(&tenantID)
    if err == sql.ErrNoRows {
        tenantID = uuid.New().String()
        _, err = tx.ExecContext(ctx, "INSERT INTO tenants (id, name, slug, created_at, updated_at) VALUES ($1, 'Default Tenant', 'default', NOW(), NOW())", tenantID)
        if err != nil {
            log.Fatalf("Failed to create default tenant: %v", err)
        }
        log.Printf("Created default tenant: %s", tenantID)
    } else if err != nil {
         log.Fatalf("Failed to query tenants: %v", err)
    }

    // 3. Ensure Program exists
    var programID string
    titleJson := `{"en": "PhD Program", "ru": "PhD Программа"}`
    
    // Check if program exists by code
    err = tx.QueryRowContext(ctx, "SELECT id FROM programs WHERE code = $1", playbook.PlaybookID).Scan(&programID)
    if err == sql.ErrNoRows {
        // Insert
        programID = uuid.New().String()
        _, err = tx.ExecContext(ctx, `
            INSERT INTO programs (id, tenant_id, code, title, name, created_at, updated_at) 
            VALUES ($1, $4, $2, $3, 'PhD Program', NOW(), NOW())
        `, programID, playbook.PlaybookID, titleJson, tenantID)
        if err != nil {
             log.Fatalf("Failed to insert program: %v", err)
        }
        log.Printf("Inserted new Program: %s", programID)
    } else if err != nil {
        log.Fatalf("Failed to query program: %v", err)
    } else {
        // Update Title (optional)
        _, err = tx.ExecContext(ctx, "UPDATE programs SET title = $2 WHERE id = $1", programID, titleJson)
        log.Printf("Found and updated existing Program: %s", programID)
    }

    // 4. Ensure Journey Map (program_version) exists
    // Targeting 'journey_maps' table
    var journeyMapID string
    mapTitleJson := `{"en": "Doctoral Journey Map", "ru": "Карта докторского пути"}`
    
    // Check if map exists
    err = tx.QueryRowContext(ctx, "SELECT id FROM journey_maps WHERE program_id = $1 AND version = $2", programID, playbook.Version).Scan(&journeyMapID)
    if err == sql.ErrNoRows {
        journeyMapID = uuid.New().String()
        _, err = tx.ExecContext(ctx, `
            INSERT INTO journey_maps (id, program_id, title, version, is_active, created_at)
            VALUES ($1, $2, $3, $4, true, NOW())
        `, journeyMapID, programID, mapTitleJson, playbook.Version)
        if err != nil {
             log.Fatalf("Failed to insert journey map: %v", err)
        }
        log.Printf("Inserted new Journey Map: %s", journeyMapID)
    } else if err != nil {
        log.Fatalf("Failed to query journey map: %v", err)
    } else {
        log.Printf("Found existing Journey Map: %s", journeyMapID)
    }
    
    // 5. Upsert Nodes into journey_node_definitions
    // Columns: id, journey_map_id, slug, type, title, description, module_key, coordinates, config, prerequisites
    
    count := 0
    for _, world := range playbook.Worlds {
        for _, node := range world.Nodes {
            // Prepare config JSON
            configData := map[string]interface{}{
                "who_can_complete": node.WhoCanComplete,
                "next":             node.Next,
                "outcomes":         node.Outcomes,
                "requirements":     node.Requirements,
                "screen":           node.Screen,
                "actionHints":      node.ActionHints,
                "condition":        node.Condition,
                "timer":            node.Timer,
            }
            configJSON, _ := json.Marshal(configData)
            titleJSON, _ := json.Marshal(node.Title)
            
            // Prepare prerequisites array
            prereqs := node.Prerequisites
            if prereqs == nil {
                prereqs = []string{}
            }
            
            // We use node.ID as slug
            
            // Upsert Node
            // Note: journey_node_definitions uses unique constraint on (journey_map_id, slug)
            _, err := tx.ExecContext(ctx, `
                INSERT INTO journey_node_definitions 
                (id, journey_map_id, slug, type, title, module_key, config, prerequisites, coordinates, created_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '{"x":0,"y":0}'::jsonb, NOW())
                ON CONFLICT (journey_map_id, slug) DO UPDATE SET
                    type = EXCLUDED.type,
                    title = EXCLUDED.title,
                    module_key = EXCLUDED.module_key,
                    config = EXCLUDED.config,
                    prerequisites = EXCLUDED.prerequisites
            `, 
                uuid.New().String(), 
                journeyMapID, 
                node.ID, 
                node.Type, 
                titleJSON, 
                world.ID, 
                configJSON, 
                pq.Array(prereqs), 
            )
            
            if err != nil {
                 log.Fatalf("Failed to upsert node %s: %v", node.ID, err)
            }
            count++
        }
    }

    if err := tx.Commit(); err != nil {
        log.Fatalf("Failed to commit transaction: %v", err)
    }

    log.Printf("Successfully migrated %d nodes to database (Journey Map ID: %s).", count, journeyMapID)
}
