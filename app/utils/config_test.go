package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// Save original environment variables
	origEnv := map[string]string{}
	for _, envVar := range [...]string{"SERVICE_PORT", "SERVICE_NAME"} {
		origEnv[envVar] = os.Getenv(envVar)
	}

	// Set test values
	os.Setenv("SERVICE_PORT", "9090")
	os.Setenv("SERVICE_NAME", "test-service")

	fmt.Fprintln(os.Stderr, "SERVICE_PORT:", os.Getenv("SERVICE_PORT"))

	cfg, err := GetConfigInstance()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ServicePort != 9090 {
		t.Errorf("Expected SERVICE_PORT=9090, got %d", cfg.ServicePort)
	}
	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected SERVICE_NAME=test-service, got %s", cfg.ServiceName)
	}

	// Restore original environment variables
	for k, v := range origEnv {
		os.Setenv(k, v)
	}
}
