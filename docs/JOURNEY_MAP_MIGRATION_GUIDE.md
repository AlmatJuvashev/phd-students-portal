# Journey Map Migration Guide: JSON → Program-Based Backend

> **Document Version:** 1.0  
> **Created:** January 3, 2026  
> **Purpose:** Migrate from static `playbook.json` to backend Program-based Journey Maps while keeping UI intact and using static templates

---

## Executive Summary

### Current Architecture
```
frontend/src/playbooks/playbook.json  →  JourneyMap component  →  NodeDetails
         ↓                                      ↓
  Static 3363-line JSON             Renders nodes by "world"
         ↓                                      ↓
  assets_list.json               Node state from /api/journey/state
```

### Target Architecture
```
Backend /api/curriculum/programs/:id/journey-map  →  JourneyMap component  →  NodeDetails
              ↓                                              ↓
    Program + JourneyMap + JourneyNodeDefinition    Same UI, API-driven
              ↓                                              ↓
    Static templates (no dynamic population)    Node state from /api/journey/state
```

### Key Constraints
- ✅ Keep existing JourneyMap UI **intact**
- ✅ Use **static templates** (no dynamic field population for now)
- ✅ Migrate node definitions to backend database
- ✅ Keep assets as static files (templates)
- ✅ Maintain backward compatibility during transition

---

## 📊 Data Structure Mapping

### JSON Playbook → Backend Models

| JSON Field | Backend Model | Database Table |
|------------|---------------|----------------|
| `playbook_id` | `Program.code` | `programs` |
| `version` | `JourneyMap.version` | `journey_maps` |
| `worlds[]` | `JourneyNodeDefinition.module_key` | `journey_node_definitions` |
| `worlds[].nodes[]` | `JourneyNodeDefinition` | `journey_node_definitions` |
| `roles[]` | Static (keep in frontend) | - |
| `conditions[]` | `JourneyNodeDefinition.config` | `journey_node_definitions` |
| `assets[]` | Static file (assets_list.json) | - |

### Node Definition Mapping

**JSON Node (`playbook.json`):**
```json
{
  "id": "S1_profile",
  "title": { "ru": "Профиль докторанта", "kz": "...", "en": "..." },
  "type": "form",
  "who_can_complete": ["student"],
  "prerequisites": [],
  "next": ["S1_text_ready"],
  "outcomes": [{ "value": "done", "next": ["S1_text_ready"] }],
  "requirements": {
    "fields": [...],
    "uploads": [...]
  }
}
```

**Backend Model (`JourneyNodeDefinition`):**
```go
type JourneyNodeDefinition struct {
  ID            string         // UUID
  JourneyMapID  string         // FK to journey_maps
  Slug          string         // "S1_profile" (stable key)
  Type          string         // "form"
  Title         string         // JSONB: {"ru": "...", "kz": "...", "en": "..."}
  Description   string         // JSONB
  ModuleKey     string         // "I" (world/module identifier)
  Config        string         // JSONB: entire requirements + outcomes
  Prerequisites []string       // [""] or ["S1_text_ready"]
  Coordinates   string         // JSONB: {"x": 0, "y": 0} for visual position
}
```

### Config JSONB Structure

The `config` field stores the rich node configuration:

```json
{
  "who_can_complete": ["student"],
  "next": ["S1_text_ready"],
  "outcomes": [
    { "value": "done", "label": {...}, "next": ["S1_text_ready"] }
  ],
  "requirements": {
    "fields": [...],
    "uploads": [...]
  },
  "screen": {...},
  "actionHints": ["form"],
  "timer": null,
  "condition": null
}
```

---

## 🔄 Migration Strategy

### Phase 1: Backend API for Journey Map (Week 1)

#### 1.1 New API Endpoint

Add to `backend/internal/handlers/api.go`:

```go
// Journey Map for Student (read their program's map)
protected.GET("/journey/map", journeyHandler.GetMyJourneyMap)

// Admin: Get any program's journey map
curr.GET("/programs/:id/journey-map", curriculumHandler.GetJourneyMap)
curr.GET("/programs/:id/journey-map/nodes", curriculumHandler.ListJourneyNodes)
```

#### 1.2 Journey Map Response Structure

Create response that matches current frontend `Playbook` type:

```go
// In handlers/curriculum.go or journey.go
type JourneyMapResponse struct {
    PlaybookID string                 `json:"playbook_id"`
    Version    string                 `json:"version"`
    Worlds     []WorldResponse        `json:"worlds"`
    Roles      []RoleResponse         `json:"roles"`
    Conditions []ConditionResponse    `json:"conditions"`
}

type WorldResponse struct {
    ID    string         `json:"id"`
    Title map[string]string `json:"title"`
    Order int            `json:"order"`
    Nodes []NodeResponse `json:"nodes"`
}

type NodeResponse struct {
    ID              string                 `json:"id"`
    Title           map[string]string      `json:"title"`
    Type            string                 `json:"type"`
    WhoCanComplete  []string               `json:"who_can_complete"`
    Prerequisites   []string               `json:"prerequisites"`
    Next            []string               `json:"next,omitempty"`
    Outcomes        []OutcomeResponse      `json:"outcomes,omitempty"`
    Condition       string                 `json:"condition,omitempty"`
    Timer           *TimerConfig           `json:"timer,omitempty"`
    Requirements    *RequirementsConfig    `json:"requirements,omitempty"`
    Outputs         []OutputConfig         `json:"outputs,omitempty"`
    ActionHints     []string               `json:"actionHints,omitempty"`
    Screen          json.RawMessage        `json:"screen,omitempty"`
}
```

#### 1.3 Handler Implementation

```go
func (h *JourneyHandler) GetMyJourneyMap(c *gin.Context) {
    userID := c.GetString("userID")
    tenantID := c.GetString("tenantID")
    
    // 1. Get user's program enrollment
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(500, gin.H{"error": "failed to get user"})
        return
    }
    
    // 2. Get program's journey map
    programID := user.ProgramID // Assuming user has program_id
    if programID == "" {
        // Fallback: return static playbook or default
        c.JSON(200, h.getDefaultPlaybook())
        return
    }
    
    journeyMap, err := h.curriculumService.GetJourneyMap(c.Request.Context(), programID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 3. Get all node definitions
    nodes, err := h.curriculumService.GetNodeDefinitions(c.Request.Context(), journeyMap.ID)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 4. Transform to frontend format
    response := h.transformToPlaybookFormat(journeyMap, nodes)
    c.JSON(200, response)
}
```

---

### Phase 2: Database Seed Script (Week 1)

#### 2.1 Migration Script

Create `backend/cmd/migrate_playbook/main.go`:

```go
package main

import (
    "encoding/json"
    "io/ioutil"
    "log"
    // ...imports
)

func main() {
    // 1. Read playbook.json
    data, err := ioutil.ReadFile("../frontend/src/playbooks/playbook.json")
    if err != nil {
        log.Fatal(err)
    }
    
    var playbook Playbook
    json.Unmarshal(data, &playbook)
    
    // 2. Connect to DB
    db := connectDB()
    
    // 3. Create/Get Program
    programID := ensureProgram(db, playbook.PlaybookID, "PhD Program")
    
    // 4. Create JourneyMap
    journeyMapID := ensureJourneyMap(db, programID, playbook.Version)
    
    // 5. Migrate each world and node
    for _, world := range playbook.Worlds {
        for _, node := range world.Nodes {
            migrateNode(db, journeyMapID, world.ID, world.Order, node)
        }
    }
    
    log.Println("Migration complete!")
}

func migrateNode(db *sqlx.DB, journeyMapID, worldID string, worldOrder int, node NodeDef) {
    // Marshal config (requirements, outcomes, etc.)
    config := NodeConfig{
        WhoCanComplete: node.WhoCanComplete,
        Next:           node.Next,
        Outcomes:       node.Outcomes,
        Requirements:   node.Requirements,
        Timer:          node.Timer,
        Condition:      node.Condition,
        ActionHints:    node.ActionHints,
        Screen:         node.Screen,
    }
    configJSON, _ := json.Marshal(config)
    titleJSON, _ := json.Marshal(node.Title)
    
    _, err := db.Exec(`
        INSERT INTO journey_node_definitions 
        (id, journey_map_id, slug, type, title, module_key, config, prerequisites, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
        ON CONFLICT (journey_map_id, slug) DO UPDATE SET
            type = EXCLUDED.type,
            title = EXCLUDED.title,
            config = EXCLUDED.config,
            prerequisites = EXCLUDED.prerequisites
    `, 
        uuid.New().String(),
        journeyMapID,
        node.ID,           // Use original ID as slug
        node.Type,
        string(titleJSON),
        worldID,           // Store world ID as module_key
        string(configJSON),
        pq.Array(node.Prerequisites),
    )
    if err != nil {
        log.Printf("Failed to migrate node %s: %v", node.ID, err)
    }
}
```

#### 2.2 Database Schema (if not exists)

```sql
-- Migration: 20260103_journey_map_tables.sql

CREATE TABLE IF NOT EXISTS journey_maps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    program_id UUID NOT NULL REFERENCES programs(id),
    title JSONB DEFAULT '{}',
    version VARCHAR(20) NOT NULL DEFAULT '1.0.0',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(program_id, version)
);

CREATE TABLE IF NOT EXISTS journey_node_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journey_map_id UUID NOT NULL REFERENCES journey_maps(id) ON DELETE CASCADE,
    parent_node_id UUID REFERENCES journey_node_definitions(id),
    slug VARCHAR(100) NOT NULL,  -- Original node ID like "S1_profile"
    type VARCHAR(50) NOT NULL,   -- form, upload, decision, etc.
    title JSONB DEFAULT '{}',
    description JSONB DEFAULT '{}',
    module_key VARCHAR(20),      -- World ID: W1, W2, etc.
    coordinates JSONB DEFAULT '{"x": 0, "y": 0}',
    config JSONB DEFAULT '{}',   -- Full node config
    prerequisites TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(journey_map_id, slug)
);

CREATE INDEX idx_journey_nodes_map ON journey_node_definitions(journey_map_id);
CREATE INDEX idx_journey_nodes_module ON journey_node_definitions(module_key);
```

---

### Phase 3: Frontend Integration (Week 2)

#### 3.1 New API Client

Create `frontend/src/api/program.ts`:

```typescript
import { api } from './client';
import type { Playbook } from '@/lib/playbook';

/**
 * Fetch the current user's program journey map from backend.
 * Returns the same Playbook structure as the static JSON.
 */
export const getMyJourneyMap = async (): Promise<Playbook> => {
  return api.get<Playbook>('/journey/map');
};

/**
 * Fetch a specific program's journey map (admin).
 */
export const getProgramJourneyMap = async (programId: string): Promise<Playbook> => {
  return api.get<Playbook>(`/curriculum/programs/${programId}/journey-map`);
};
```

#### 3.2 Update DoctoralJourney Page

Replace static JSON import with API call:

```typescript
// pages/doctoral.journey.tsx
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";

import { JourneyMap } from "@/components/map/JourneyMap";
import { ResetBar } from "@/features/journey/components/ResetBar";
import { useJourneyState } from "@/features/journey/hooks";
import type { Playbook } from "@/lib/playbook";
import { useRequireAuth } from '@/hooks/useRequireAuth';
import { getMyJourneyMap } from "@/api/program";

// Static fallback for development/offline
import staticPlaybook from "@/playbooks/playbook.json";

export function DoctoralJourney() {
  const { t: T, i18n } = useTranslation("common");
  const { isLoading: authLoading } = useRequireAuth();
  const { state = {}, refetch } = useJourneyState();

  // Fetch journey map from API, fallback to static
  const { data: playbook, isLoading: mapLoading, error } = useQuery({
    queryKey: ['journey', 'map'],
    queryFn: getMyJourneyMap,
    staleTime: 5 * 60 * 1000, // 5 minutes cache
    retry: 1,
    // Fallback to static playbook on error
    placeholderData: staticPlaybook as Playbook,
  });

  // Use static playbook if API fails
  const effectivePlaybook = error ? (staticPlaybook as Playbook) : playbook;

  if (authLoading || mapLoading || !effectivePlaybook) {
    return (
      <div className="flex items-center justify-center py-16">
        <p className="text-sm text-muted-foreground animate-pulse">
          {T("map.loading", { defaultValue: "Loading dissertation map…" })}
        </p>
      </div>
    );
  }

  return (
    <div>
      <JourneyMap
        playbook={effectivePlaybook}
        locale={i18n.language}
        stateByNodeId={state as any}
        onStateChanged={refetch}
      />
      <ResetBar />
    </div>
  );
}
```

#### 3.3 No Changes to JourneyMap Component

The `JourneyMap` component remains **unchanged** because:
- It receives `playbook: Playbook` as a prop
- The API returns the exact same structure
- All rendering logic stays the same

---

### Phase 4: Assets & Templates (Keep Static)

#### 4.1 Keep assets_list.json

Templates remain as static files. The `assets_list.json` continues to be loaded client-side:

```typescript
// lib/assets.ts - NO CHANGES NEEDED
import assetsList from "@/playbooks/assets_list.json";

export function getAssetUrl(assetId: string): string {
  const asset = assetsList.assets.find(a => a.id === assetId);
  if (!asset) return '';
  return `/templates/${asset.storage.key}`;
}
```

#### 4.2 Template Files Location

Keep templates in `frontend/public/templates/` or serve from S3:
```
frontend/public/templates/
├── Zayavlenie_v_OMiD_ru.docx
├── Zayavlenie_v_OMiD_kz.docx
├── Letter_to_LCB_ru.docx
└── ...
```

No dynamic population - users download and fill manually.

---

## 📋 Implementation Checklist

### Backend Tasks

- [ ] Add database tables (`journey_maps`, `journey_node_definitions`)
- [ ] Create migration script to import `playbook.json`
- [ ] Add `GET /api/journey/map` endpoint
- [ ] Add `GET /api/curriculum/programs/:id/journey-map` endpoint
- [ ] Transform backend models to frontend `Playbook` format
- [ ] Add user → program association (if not exists)
- [ ] Write tests for new endpoints

### Frontend Tasks

- [ ] Create `api/program.ts` with `getMyJourneyMap()`
- [ ] Update `DoctoralJourney` to use API with fallback
- [ ] Add React Query caching for journey map
- [ ] Keep static `playbook.json` as fallback
- [ ] Test offline/error scenarios
- [ ] No changes to `JourneyMap.tsx` or `lib/playbook.ts`

### Migration Tasks

- [ ] Run migration script on staging
- [ ] Verify all nodes imported correctly
- [ ] Compare API response with static JSON
- [ ] Test student journey progression
- [ ] Test admin journey map viewing
- [ ] Deploy to production

---

## 🗺️ World/Module Structure Preservation

The current playbook has these "worlds":

| World ID | Title (RU) | Order | Node Count |
|----------|------------|-------|------------|
| W1 | I — Подготовка | 1 | ~15 nodes |
| W2 | II — Экспертиза | 2 | ~20 nodes |
| W3 | III — Условный | 3 | ~5 nodes (conditional) |
| W4 | IV — Защита | 4 | ~15 nodes |
| W5 | V — Пост-защита | 5 | ~10 nodes |
| W6 | VI — Аттестация | 6 | ~8 nodes |

These map to `journey_node_definitions.module_key`:
- `W1` → `module_key = "W1"`
- Frontend groups by `module_key` to reconstruct worlds

---

## 🔀 Transition Strategy

### Option A: Big Bang (Simple)
1. Migrate all data to backend
2. Switch frontend to API
3. Remove static JSON dependency

### Option B: Gradual (Recommended)
1. **Week 1**: Backend API + migration script
2. **Week 2**: Frontend uses API with static fallback
3. **Week 3**: Monitor, fix issues
4. **Week 4**: Remove static fallback, delete `playbook.json`

### Feature Flag Approach

```typescript
// config/features.ts
export const FEATURES = {
  USE_BACKEND_JOURNEY_MAP: import.meta.env.VITE_USE_BACKEND_MAP === 'true',
};

// In DoctoralJourney
const playbook = FEATURES.USE_BACKEND_JOURNEY_MAP 
  ? await getMyJourneyMap()
  : staticPlaybook;
```

---

## 📈 Benefits After Migration

1. **Multi-Program Support**: Different programs have different journey maps
2. **Version Control**: Track changes to journey maps over time
3. **Admin Editing**: Future ability to edit nodes via admin UI
4. **Per-Tenant Customization**: Each tenant can have custom journeys
5. **Analytics**: Better tracking of which nodes are used
6. **Reduced Bundle Size**: Remove 3363-line JSON from frontend bundle

---

## ⚠️ Risk Mitigation

| Risk | Mitigation |
|------|------------|
| API failure | Static JSON fallback |
| Data mismatch | Migration validation script |
| Performance | Cache API response (5 min) |
| Breaking changes | Feature flag for gradual rollout |
| Missing nodes | Compare node counts pre/post migration |

---

## 🧪 Testing Plan

### Unit Tests
- Backend: Test `GetJourneyMap` handler
- Backend: Test node definition transformation
- Frontend: Test API client error handling

### Integration Tests
- API returns valid Playbook structure
- All node types render correctly
- Prerequisites work as expected
- State transitions function properly

### Manual Tests
- Complete full journey flow
- Test each node type (form, upload, decision, etc.)
- Test conditional nodes (W3)
- Test as different roles (student, advisor)

---

## 📅 Timeline

| Week | Tasks |
|------|-------|
| 1 | Backend API, database schema, migration script |
| 2 | Frontend integration, testing |
| 3 | Staging deployment, bug fixes |
| 4 | Production deployment, remove static fallback |

## 🏁 Next Steps: Roadmap

1.  **Backend Setup**: Apply the `20260103_journey_map_tables.sql` migration.
2.  **Data Ingestion**: Run the `migrate_playbook` script to populate the database from the static `playbook.json`.
3.  **API Verification**: Test the `GET /api/journey/map` endpoint to ensure it returns the valid Playbook structure.
4.  **Frontend Switch**: Update `DoctoralJourney` to fetch data from the API while maintaining the static fallback.
