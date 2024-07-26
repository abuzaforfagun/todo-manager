package models_test

import (
	"restful-service/models"
	"testing"
)

func TestToString(t *testing.T) {
	statusTests := []struct {
		name     string
		status   models.Status
		expected string
	}{
		{"Should return Completed", models.Completed, "Completed"},
		{"Should return In Progress", models.InProgress, "In Progress"},
		{"Should return Pending", models.Pending, "Pending"},
	}

	for _, tt := range statusTests {
		t.Run(tt.name, func(t *testing.T) {

			status := tt.status
			result := status.ToString()

			if result != tt.expected {
				t.Errorf("Got %s, expected %s", result, tt.expected)
			}
		})
	}
}
