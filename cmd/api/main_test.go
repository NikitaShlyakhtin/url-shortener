package main

import "testing"

func TestValidateConfig(t *testing.T) {
	cfg := config{
		baseUrl: "",
		storage: storageConfig{
			storage_type: "in-memory",
		},
	}
	err := validateConfig(cfg)
	if err == nil || err.Error() != "base URL is required" {
		t.Errorf("Expected error: base URL is required, got: %v", err)
	}

	cfg = config{
		baseUrl: "https://example.com",
		storage: storageConfig{
			storage_type: "invalid",
		},
	}
	err = validateConfig(cfg)
	if err == nil || err.Error() != "invalid storage type" {
		t.Errorf("Expected error: invalid storage type, got: %v", err)
	}

	cfg = config{
		baseUrl: "https://example.com",
		storage: storageConfig{
			storage_type: "postgres",
			db: databaseConfig{
				dsn: "",
			},
		},
	}
	err = validateConfig(cfg)
	if err == nil || err.Error() != "DSN is required for PostgreSQL storage" {
		t.Errorf("Expected error: DSN is required for PostgreSQL storage, got: %v", err)
	}

	cfg = config{
		baseUrl: "https://example.com",
		storage: storageConfig{
			storage_type: "in-memory",
		},
	}
	err = validateConfig(cfg)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
