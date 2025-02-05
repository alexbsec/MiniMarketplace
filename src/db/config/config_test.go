package config

import (
	"testing"
)

func TestInitService(t *testing.T) {
	// Call InitService to test if the service initializes correctly
	service, err := InitService()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatalf("Expected service to be initialized, got nil")
	}

	// Check if the database connection is valid
	sqlDB, err := service.db.DB()
	if err != nil {
		t.Fatalf("Failed to get SQL DB from GORM: %v", err)
	}

	// Ping the database to ensure it's connected
	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Database connection test failed: %v", err)
	}
}

