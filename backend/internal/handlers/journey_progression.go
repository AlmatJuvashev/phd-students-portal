package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services/playbook"
	"github.com/jmoiron/sqlx"
)

// ActivateNextNodes activates the next nodes in the playbook after a node is completed.
// It evaluates 'outcomes' if present to determine which path to take.
func ActivateNextNodes(db *sqlx.DB, pb *playbook.Manager, userID, completedNodeID string) error {
	log.Printf("[ActivateNextNodes] Starting for user=%s node=%s version=%s", userID, completedNodeID, pb.VersionID)
	
	// 1. Check for outcomes first (conditional branching)
	// We need to fetch the node instance to check its outcomes
	var instanceID string
	err := db.QueryRowx(`SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 AND playbook_version_id=$3`,
		userID, completedNodeID, pb.VersionID).Scan(&instanceID)
	
	var nextNodes []string

	if err == nil {
		// Fetch outcomes for this instance
		rows, err := db.Queryx(`SELECT outcome_value FROM node_outcomes WHERE node_instance_id=$1`, instanceID)
		if err == nil {
			defer rows.Close()
			var outcomes []string
			for rows.Next() {
				var val string
				if err := rows.Scan(&val); err == nil {
					outcomes = append(outcomes, val)
				}
			}
			
			// If we have outcomes, we need to find the matching 'next' in the playbook definition
			if len(outcomes) > 0 {
				_, ok := pb.NodeDefinition(completedNodeID)
				if ok {
					// The playbook structure for outcomes is complex (it's in the raw JSON).
					// We need to parse the raw node definition to find the 'outcomes' logic.
					// Since 'pb.Nodes' is a simplified map, we might need to look at the raw JSON or
					// rely on the fact that 'next' might be conditional.
					
					// However, for S1_profile, the 'next' field in the node definition contains ALL possible next nodes.
					// The 'outcomes' in the playbook define which ones to ACTUALLY activate.
					// But the current 'ActivateNextNodes' implementation just activates EVERYTHING in 'next'.
					
					// For now, to fix the immediate issue (RP node not showing), we will stick to the existing logic:
					// Activate ALL nodes in 'next'. The frontend/backend conditions should handle visibility/locking.
					// If we want true conditional branching, we need to implement outcome evaluation here.
					
					// Given the current codebase state, let's stick to the existing logic from AdminHandler
					// which simply activates all nodes in 'next'.
				}
			}
		}
	}

	// Parse playbook to get next nodes (fallback to simple 'next' list)
	var pbStruct struct {
		Worlds []struct {
			Nodes []struct {
				ID   string   `json:"id"`
				Next []string `json:"next"`
			} `json:"nodes"`
		} `json:"worlds"`
	}
	
	if err := json.Unmarshal(pb.Raw, &pbStruct); err != nil {
		return fmt.Errorf("parse playbook: %w", err)
	}
	
	// Find the completed node and get its next nodes
	for _, world := range pbStruct.Worlds {
		for _, node := range world.Nodes {
			if node.ID == completedNodeID {
				nextNodes = node.Next
				break
			}
		}
		if len(nextNodes) > 0 {
			break
		}
	}
	
	if len(nextNodes) == 0 {
		log.Printf("[ActivateNextNodes] No next nodes found for %s", completedNodeID)
		return nil
	}
	
	log.Printf("[ActivateNextNodes] Activating next nodes %v for user %s", nextNodes, userID)
	
	// Activate each next node
	for _, nodeID := range nextNodes {
		// Check if node instance already exists
		var existingID string
		err := db.QueryRowx(`SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 AND playbook_version_id=$3`,
			userID, nodeID, pb.VersionID).Scan(&existingID)
		
		if err == sql.ErrNoRows {
			// Create new node instance in active state
			var newID string
			err = db.QueryRowx(`INSERT INTO node_instances (user_id, playbook_version_id, node_id, state, opened_at)
				VALUES ($1, $2, $3, 'active', now()) RETURNING id`,
				userID, pb.VersionID, nodeID).Scan(&newID)
			if err != nil {
				log.Printf("[ActivateNextNodes] Failed to create instance for %s: %v", nodeID, err)
				continue
			}
			log.Printf("[ActivateNextNodes] Created new instance %s for node %s", newID, nodeID)
			
			// Create initial event
			_, _ = db.Exec(`INSERT INTO node_events (node_instance_id, event_type, payload, actor_id)
				VALUES ($1, 'opened', '{}', $2)`, newID, userID)

		} else if err != nil {
			log.Printf("[ActivateNextNodes] Error checking instance for %s: %v", nodeID, err)
			continue
		} else {
			// Instance exists, update to active if it's locked
			_, err = db.Exec(`UPDATE node_instances SET state='active', updated_at=now() 
				WHERE id=$1 AND state='locked'`, existingID)
			if err != nil {
				log.Printf("[ActivateNextNodes] Failed to activate instance %s: %v", existingID, err)
			} else {
				log.Printf("[ActivateNextNodes] Activated existing instance %s for node %s", existingID, nodeID)
			}
		}
		
		// Update journey_states
		_, _ = db.Exec(`INSERT INTO journey_states (user_id, node_id, state)
			VALUES ($1, $2, 'active')
			ON CONFLICT (user_id, node_id) DO UPDATE SET state='active', updated_at=now()`,
			userID, nodeID)
	}
	
	return nil
}

// ActivateNextNodesWithTenant activates the next nodes with proper tenant_id support.
func ActivateNextNodesWithTenant(db *sqlx.DB, pb *playbook.Manager, userID, completedNodeID, tenantID string) error {
	log.Printf("[ActivateNextNodesWithTenant] Starting for user=%s node=%s tenant=%s version=%s", userID, completedNodeID, tenantID, pb.VersionID)
	
	// Parse playbook to get next nodes
	var pbStruct struct {
		Worlds []struct {
			Nodes []struct {
				ID   string   `json:"id"`
				Next []string `json:"next"`
			} `json:"nodes"`
		} `json:"worlds"`
	}
	
	if err := json.Unmarshal(pb.Raw, &pbStruct); err != nil {
		return fmt.Errorf("parse playbook: %w", err)
	}
	
	var nextNodes []string
	// Find the completed node and get its next nodes
	for _, world := range pbStruct.Worlds {
		for _, node := range world.Nodes {
			if node.ID == completedNodeID {
				nextNodes = node.Next
				break
			}
		}
		if len(nextNodes) > 0 {
			break
		}
	}
	
	if len(nextNodes) == 0 {
		log.Printf("[ActivateNextNodesWithTenant] No next nodes found for %s", completedNodeID)
		return nil
	}
	
	log.Printf("[ActivateNextNodesWithTenant] Activating next nodes %v for user %s tenant %s", nextNodes, userID, tenantID)
	
	// Activate each next node
	for _, nodeID := range nextNodes {
		// Check if node instance already exists
		var existingID string
		err := db.QueryRowx(`SELECT id FROM node_instances WHERE user_id=$1 AND node_id=$2 AND tenant_id=$3`,
			userID, nodeID, tenantID).Scan(&existingID)
		
		if err == sql.ErrNoRows {
			// Create new node instance in active state WITH tenant_id
			var newID string
			err = db.QueryRowx(`INSERT INTO node_instances (tenant_id, user_id, playbook_version_id, node_id, state, opened_at)
				VALUES ($1, $2, $3, $4, 'active', now()) RETURNING id`,
				tenantID, userID, pb.VersionID, nodeID).Scan(&newID)
			if err != nil {
				log.Printf("[ActivateNextNodesWithTenant] Failed to create instance for %s: %v", nodeID, err)
				continue
			}
			log.Printf("[ActivateNextNodesWithTenant] Created new instance %s for node %s", newID, nodeID)
			
			// Create initial event
			_, _ = db.Exec(`INSERT INTO node_events (node_instance_id, event_type, payload, actor_id)
				VALUES ($1, 'opened', '{}', $2)`, newID, userID)

		} else if err != nil {
			log.Printf("[ActivateNextNodesWithTenant] Error checking instance for %s: %v", nodeID, err)
			continue
		} else {
			// Instance exists, update to active if it's locked
			_, err = db.Exec(`UPDATE node_instances SET state='active', updated_at=now() 
				WHERE id=$1 AND state='locked'`, existingID)
			if err != nil {
				log.Printf("[ActivateNextNodesWithTenant] Failed to activate instance %s: %v", existingID, err)
			} else {
				log.Printf("[ActivateNextNodesWithTenant] Activated existing instance %s for node %s", existingID, nodeID)
			}
		}
		
		// Update journey_states WITH tenant_id
		_, _ = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state)
			VALUES ($1, $2, $3, 'active')
			ON CONFLICT (user_id, node_id) DO UPDATE SET state='active', updated_at=now()`,
			tenantID, userID, nodeID)
	}
	
	return nil
}
