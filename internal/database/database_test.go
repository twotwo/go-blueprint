package database

import (
	"testing"
)

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
