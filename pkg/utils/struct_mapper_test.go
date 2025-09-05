package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test struct to simulate UpdateTicketRequest
type TestUpdateRequest struct {
	ID          string  `json:"id" gorm:"column:id"`
	Subject     *string `json:"subject" gorm:"column:subject"`
	Description *string `json:"description" gorm:"column:description"`
	Priority    *string `json:"priority" gorm:"column:priority"`
	Status      *string `json:"status" gorm:"column:status"`
	AssignedTo  *string `json:"assigned_to" gorm:"column:assigned_to_id"`
	UpdatedBy   *string `json:"updated_by" gorm:"column:updated_by"`
}

func TestStructToMap(t *testing.T) {
	subject := "Test Subject"
	priority := "high"

	req := &TestUpdateRequest{
		ID:       "ticket-123",
		Subject:  &subject,
		Priority: &priority,
		// Description is nil (should be skipped)
		// Status is nil (should be skipped)
		// AssignedTo is nil (should be skipped)
	}

	result := StructToMap(req)

	// Should include non-nil pointer fields
	assert.Equal(t, "ticket-123", result["id"])
	assert.Equal(t, "Test Subject", result["subject"])
	assert.Equal(t, "high", result["priority"])

	// Should not include nil pointer fields
	assert.NotContains(t, result, "description")
	assert.NotContains(t, result, "status")
	assert.NotContains(t, result, "assigned_to_id")
}

func TestToUpdateMap(t *testing.T) {
	subject := "Updated Subject"
	priority := "urgent"
	assignedTo := "agent-456"

	req := &TestUpdateRequest{
		ID:         "ticket-123",
		Subject:    &subject,
		Priority:   &priority,
		AssignedTo: &assignedTo,
	}

	result := ToUpdateMap(req)

	// Should include updatable fields
	assert.Equal(t, "Updated Subject", result["subject"])
	assert.Equal(t, "urgent", result["priority"])
	assert.Equal(t, "agent-456", result["assigned_to_id"])

	// Should exclude ID and other non-updatable fields
	assert.NotContains(t, result, "id")

	// Should automatically add updated_at
	assert.Contains(t, result, "updated_at")
	assert.IsType(t, time.Time{}, result["updated_at"])
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ID", "id"},
		{"Subject", "subject"},
		{"AssignedTo", "assigned_to"},
		{"UpdatedBy", "updated_by"},
		{"CustomerName", "customer_name"},
	}

	for _, test := range tests {
		result := toSnakeCase(test.input)
		assert.Equal(t, test.expected, result, "toSnakeCase(%s) should return %s", test.input, test.expected)
	}
}
