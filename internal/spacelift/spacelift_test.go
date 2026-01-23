package spacelift

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadModuleVersion_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()

	version := ReadModuleVersion(tmpDir)
	if version != "" {
		t.Errorf("expected empty version, got '%s'", version)
	}
}

func TestReadModuleVersion_WithConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .spacelift directory and config.yml
	spaceliftDir := filepath.Join(tmpDir, DirSpacelift)
	if err := os.MkdirAll(spaceliftDir, 0755); err != nil {
		t.Fatalf("failed to create .spacelift dir: %v", err)
	}

	configContent := `module_version: "1.2.3"`
	configPath := filepath.Join(spaceliftDir, FileConfig)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	version := ReadModuleVersion(tmpDir)
	if version != "1.2.3" {
		t.Errorf("expected '1.2.3', got '%s'", version)
	}
}

func TestReadModuleVersion_InvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()

	spaceliftDir := filepath.Join(tmpDir, DirSpacelift)
	if err := os.MkdirAll(spaceliftDir, 0755); err != nil {
		t.Fatalf("failed to create .spacelift dir: %v", err)
	}

	// Write invalid YAML
	configPath := filepath.Join(spaceliftDir, FileConfig)
	if err := os.WriteFile(configPath, []byte("not: valid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	version := ReadModuleVersion(tmpDir)
	if version != "" {
		t.Errorf("expected empty version for invalid yaml, got '%s'", version)
	}
}

func TestReadModuleVersion_EmptyVersion(t *testing.T) {
	tmpDir := t.TempDir()

	spaceliftDir := filepath.Join(tmpDir, DirSpacelift)
	if err := os.MkdirAll(spaceliftDir, 0755); err != nil {
		t.Fatalf("failed to create .spacelift dir: %v", err)
	}

	// Config without module_version
	configContent := `other_field: "value"`
	configPath := filepath.Join(spaceliftDir, FileConfig)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	version := ReadModuleVersion(tmpDir)
	if version != "" {
		t.Errorf("expected empty version, got '%s'", version)
	}
}
