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
func ActivateNextNodes(db *sqlx.DB, pb *playbook.Manager, userID, completedNodeID, tenantID string) error {
	log.Printf("[ActivateNextNodes] Starting for user=%s node=%s tenant=%s version=%s", userID, completedNodeID, tenantID, pb.VersionID)
	
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
		log.Printf("[ActivateNextNodes] No next nodes found for %s", completedNodeID)
		return nil
	}
	
	log.Printf("[ActivateNextNodes] Activating next nodes %v for user %s tenant %s", nextNodes, userID, tenantID)
	
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
		
		// Update journey_states WITH tenant_id
		_, _ = db.Exec(`INSERT INTO journey_states (tenant_id, user_id, node_id, state)
			VALUES ($1, $2, $3, 'active')
			ON CONFLICT (user_id, node_id) DO UPDATE SET state='active', updated_at=now()`,
			tenantID, userID, nodeID)
	}
	
	return nil
}
