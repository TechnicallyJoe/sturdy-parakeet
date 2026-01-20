package finder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindModule_SingleMatch(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create components/azurerm/storage-account with .tf file
	modulePath := filepath.Join(tmpDir, "azurerm", "storage-account")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	tfFile := filepath.Join(modulePath, "main.tf")
	if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
		t.Fatalf("failed to create .tf file: %v", err)
	}

	matches, err := FindModule(tmpDir, "storage-account")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	if matches[0] != modulePath {
		t.Errorf("expected match to be '%s', got '%s'", modulePath, matches[0])
	}
}

func TestFindModule_MultipleMatches(t *testing.T) {
	// Create temp directory structure with name clash
	tmpDir := t.TempDir()

	// Create two modules with the same name under different providers
	module1 := filepath.Join(tmpDir, "azurerm", "storage-account")
	module2 := filepath.Join(tmpDir, "aws", "storage-account")

	for _, path := range []string{module1, module2} {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create module directory: %v", err)
		}
		tfFile := filepath.Join(path, "main.tf")
		if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
			t.Fatalf("failed to create .tf file: %v", err)
		}
	}

	matches, err := FindModule(tmpDir, "storage-account")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestFindModule_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a module with a different name
	modulePath := filepath.Join(tmpDir, "other-module")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	tfFile := filepath.Join(modulePath, "main.tf")
	if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
		t.Fatalf("failed to create .tf file: %v", err)
	}

	matches, err := FindModule(tmpDir, "storage-account")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("expected 0 matches, got %d", len(matches))
	}
}

func TestFindModule_IgnoresDirectoriesWithoutTfFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory with matching name but no .tf files
	modulePath := filepath.Join(tmpDir, "storage-account")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	// Create a non-.tf file
	otherFile := filepath.Join(modulePath, "README.md")
	if err := os.WriteFile(otherFile, []byte("# README"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	matches, err := FindModule(tmpDir, "storage-account")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("expected 0 matches (no .tf files), got %d", len(matches))
	}
}

func TestFindModule_MatchesTfJsonFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create module with .tf.json file
	modulePath := filepath.Join(tmpDir, "storage-account")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	tfJsonFile := filepath.Join(modulePath, "main.tf.json")
	if err := os.WriteFile(tfJsonFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to create .tf.json file: %v", err)
	}

	matches, err := FindModule(tmpDir, "storage-account")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match for .tf.json, got %d", len(matches))
	}
}

func TestFindModule_NestedDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create deeply nested module
	modulePath := filepath.Join(tmpDir, "level1", "level2", "level3", "my-module")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	tfFile := filepath.Join(modulePath, "main.tf")
	if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
		t.Fatalf("failed to create .tf file: %v", err)
	}

	matches, err := FindModule(tmpDir, "my-module")
	if err != nil {
		t.Fatalf("FindModule returned error: %v", err)
	}

	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}

	if matches[0] != modulePath {
		t.Errorf("expected match to be '%s', got '%s'", modulePath, matches[0])
	}
}

func TestFindModule_SearchPathDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist")

	_, err := FindModule(nonExistentPath, "any-module")
	if err == nil {
		t.Error("expected error for non-existent search path, got nil")
	}
}

func TestHasTerraformFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected bool
	}{
		{
			name:     "has .tf file",
			files:    []string{"main.tf"},
			expected: true,
		},
		{
			name:     "has .tf.json file",
			files:    []string{"main.tf.json"},
			expected: true,
		},
		{
			name:     "has multiple .tf files",
			files:    []string{"main.tf", "variables.tf", "outputs.tf"},
			expected: true,
		},
		{
			name:     "has mixed files",
			files:    []string{"main.tf", "README.md"},
			expected: true,
		},
		{
			name:     "no .tf files",
			files:    []string{"README.md", "config.yaml"},
			expected: false,
		},
		{
			name:     "empty directory",
			files:    []string{},
			expected: false,
		},
		{
			name:     "only .json file (not .tf.json)",
			files:    []string{"config.json"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			result := hasTerraformFiles(tmpDir)
			if result != tt.expected {
				t.Errorf("hasTerraformFiles() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestHasTerraformFiles_NonExistentDir(t *testing.T) {
	result := hasTerraformFiles("/non/existent/path")
	if result != false {
		t.Error("expected false for non-existent directory")
	}
}
