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
         // It might be that tenants table doesn't exist or other error.
         // Assuming it exists because programs has tenant_id FK.
         log.Fatalf("Failed to query tenants: %v", err)
    }

    // 3. Ensure Program exists
    var programID string
    titleJson := `{"en": "PhD Program", "ru": "PhD Программа"}`
    
    // Check if program exists
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

    // 4. Ensure Program Version (aka Journey Map) exists
    var programVersionID string
    mapTitleJson := `{"en": "Doctoral Journey Map", "ru": "Карта докторского пути"}`
    
    // Prepare Config (Phases)
    // Map Worlds to Phases for Builder
    type Phase struct {
        ID       string                 `json:"id"`
        Title    string                 `json:"title"` // simplified for now, or LocalizedString
        Order    int                    `json:"order"`
        Color    string                 `json:"color,omitempty"`
        Position map[string]float64     `json:"position"`
    }
    
    phases := []Phase{}
    xPos := 0.0
    for _, w := range playbook.Worlds {
         // Best effort localization
         title := ""
         if t, ok := w.Title["en"]; ok { title = t } else if t, ok := w.Title["ru"]; ok { title = t }
         
         phases = append(phases, Phase{
             ID: w.ID,
             Title: title,
             Order: w.Order,
             Color: "#6366f1", // Default color
             Position: map[string]float64{"x": float64(xPos), "y": 50},
         })
         xPos += 400 // space them out
    }
    
    mapConfig := map[string]interface{}{
        "phases": phases,
    }
    mapConfigJson, _ := json.Marshal(mapConfig)

    // Check if map exists
    err = tx.QueryRowContext(ctx, "SELECT id FROM program_versions WHERE program_id = $1 AND version = $2", programID, playbook.Version).Scan(&programVersionID)
    if err == sql.ErrNoRows {
        programVersionID = uuid.New().String()
        _, err = tx.ExecContext(ctx, `
            INSERT INTO program_versions (id, program_id, title, version, config, is_active, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, true, NOW(), NOW())
        `, programVersionID, programID, mapTitleJson, playbook.Version, mapConfigJson)
        if err != nil {
             log.Fatalf("Failed to insert journey map: %v", err)
        }
        log.Printf("Inserted new Program Version: %s", programVersionID)
    } else if err != nil {
        log.Fatalf("Failed to query journey map: %v", err)
    } else {
        log.Printf("Found existing Program Version: %s", programVersionID)
         // Update config
        _, err = tx.ExecContext(ctx, "UPDATE program_versions SET config = $2, updated_at = NOW() WHERE id = $1", programVersionID, mapConfigJson)
        if err != nil {
             log.Printf("Warning: Failed to update map config: %v", err)
        }
    }
    log.Printf("Program Version ID: %s", programVersionID)

    // 5. Delete existing nodes for this map version to ensure clean slate (or we can upsert one by one)
    // Upsert is safer to preserve IDs if we were using consistent UUIDs, but here we generate new UUIDs unless we have a mapping.
    // However, the text says "Use original ID as slug".
    // Let's just UPSERT based on (program_version_id, slug).

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
            
            // Upsert Node
            // Note: program_version_node_definitions uses unique constraint on (program_version_id, slug)
            _, err := tx.ExecContext(ctx, `
                INSERT INTO program_version_node_definitions 
                (id, program_version_id, slug, type, title, module_key, config, prerequisites, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
                ON CONFLICT (program_version_id, slug) DO UPDATE SET
                    type = EXCLUDED.type,
                    title = EXCLUDED.title,
                    module_key = EXCLUDED.module_key,
                    config = EXCLUDED.config,
                    prerequisites = EXCLUDED.prerequisites,
                    updated_at = NOW()
            `, 
                uuid.New().String(), 
                programVersionID, 
                node.ID, 
                node.Type, 
                titleJSON, 
                world.ID, 
                configJSON, 
                pq.Array(prereqs), 
            )
            
            // Wait, standard pgx/sql driver doesn't support array types automatically without special handling or pq (which I imported via jackc adapter?)
            // Actually I'm using "github.com/jackc/pgx/v5/stdlib".
            // Arrays in standard sql can be tricky. Let's assume the column is text[] or jsonb?
            // Migration script usually uses `text[]`.
            // For pgx stdlib, passing []string usually works if driver handles it, or I might need `pq.Array`.
            // Let's try passing the slice directly first. If it fails, I'll switch to JSON string if column is JSONB, or use pq.Array if I import lib/pq.
            // I'll add a helper `toJsonArray` to be safe if it's JSONB, or just try raw.
            // Wait, looking at current `go.mod`, do we have `lib/pq`?
            // The migration guide used `pq.Array`. I should check if I can use it.
            // I'll stick to `pq` import if available, or just standard slice.
            
            if err != nil {
                 // Try again with different array handling if needed, but for now log fatal
                 log.Fatalf("Failed to upsert node %s: %v", node.ID, err)
            }
            count++
        }
    }

    if err := tx.Commit(); err != nil {
        log.Fatalf("Failed to commit transaction: %v", err)
    }

    log.Printf("Successfully migrated %d nodes to database.", count)
}

func toJsonArray(arr []string) interface{} {
    // If the driver supports []string for text[], return it.
    // pgx generally does.
    return arr
}
