package main

import (
	"fmt"

	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/pkg/utils"
)

func main() {
	// Simulate UpdateTicketRequest
	subject := "Updated Subject"
	priority := "high"
	status := "in_progress"
	assignedTo := "agent-123"

	req := &domain.UpdateTicketRequest{
		ID:         "ticket-456",
		Subject:    &subject,
		Priority:   &priority,
		Status:     &status,
		AssignedTo: &assignedTo,
		// Description is nil - should be skipped
	}

	fmt.Println("Original struct:")
	fmt.Printf("ID: %s\n", req.ID)
	fmt.Printf("Subject: %v\n", req.Subject)
	fmt.Printf("Priority: %v\n", req.Priority)
	fmt.Printf("Status: %v\n", req.Status)
	fmt.Printf("AssignedTo: %v\n", req.AssignedTo)
	fmt.Printf("Description: %v\n", req.Description)

	fmt.Println("\nConverted to update map:")
	updates := utils.ToUpdateMap(req)
	for key, value := range updates {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Test with StructToMap directly
	fmt.Println("\nDirect StructToMap result:")
	directMap := utils.StructToMap(req)
	for key, value := range directMap {
		fmt.Printf("%s: %v\n", key, value)
	}
}
