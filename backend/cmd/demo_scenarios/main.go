package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/repository"
	"github.com/AlmatJuvashev/phd-students-portal/backend/internal/services"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	db, _ := sqlx.Connect("postgres", dbURL)
	defer db.Close()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [approve|reject|notify]")
		return
	}

	command := os.Args[1]
	
	// Repos & Services
	workflowRepo := repository.NewSQLWorkflowRepository(db)
	workflowService := services.NewWorkflowService(workflowRepo)

	ctx := context.Background()

	switch command {
	case "approve":
		// Find first pending approval for Dean role
		var approvalID, instanceID, deanID uuid.UUID
		err := db.QueryRow(`
			SELECT a.id, a.instance_id, u.id 
			FROM workflow_approvals a
			JOIN users u ON u.role = 'dean'
			WHERE a.approver_role = 'dean' AND (a.decision IS NULL OR a.decision = '')
			LIMIT 1`).Scan(&approvalID, &instanceID, &deanID)
		
		if err != nil {
			fmt.Println("No pending Dean approvals found.")
			return
		}

		fmt.Printf("Approving workflow %s for Dean %s...\n", instanceID, deanID)
		err = workflowService.ApproveStep(ctx, approvalID, deanID, "Approved via Scenario Controller")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success!")

	case "notify":
		fmt.Println("Broadcasting notification to all students...")
		_, _ = db.Exec(`INSERT INTO notifications (recipient_id, title, message, type) 
			SELECT id, 'Emergency Update', 'Please check your portal for new dissertation guidelines.', 'warning' 
			FROM users WHERE role = 'student'`)
		fmt.Println("Done.")

	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
