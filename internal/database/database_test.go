package database

import (
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

func TestDatabaseConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		expected string
	}{
		{
			name:     "Check DB host",
			envVar:   "BLUEPRINT_DB_HOST",
			expected: "psql_bp", // Update with expected value
		},
		{
			name:     "Check DB port",
			envVar:   "BLUEPRINT_DB_PORT",
			expected: "5432", // Update with expected value
		},
		{
			name:     "Check DB name",
			envVar:   "BLUEPRINT_DB_DATABASE",
			expected: "blueprint", // Update with expected value
		},
		{
			name:     "Check DB user",
			envVar:   "BLUEPRINT_DB_USERNAME",
			expected: "melkey", // Update with expected value
		},
		{
			name:     "Check DB schema",
			envVar:   "BLUEPRINT_DB_SCHEMA",
			expected: "web", // Update with expected value
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := os.Getenv(tt.envVar)
			if got != tt.expected {
				t.Errorf("Expected %s=%s, got %s", tt.envVar, tt.expected, got)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successfully create new database connection",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if got == nil && !tt.wantErr {
				t.Error("New() returned nil, expected valid connection")
			}
			defer got.Close()
		})
	}
}

func TestClose(t *testing.T) {
	srv := New()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}

func TestSingletonPattern(t *testing.T) {
	// Test that New() returns the same instance
	db1 := New()
	db2 := New()

	if db1 != db2 {
		t.Error("New() should return the same instance (singleton pattern)")
	}
}
