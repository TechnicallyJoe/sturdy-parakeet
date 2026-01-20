package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/config"
)

func TestResolveExplicitPath_AbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test directory
	testPath := filepath.Join(tmpDir, "test-module")
	if err := os.MkdirAll(testPath, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	result, err := resolveExplicitPath(testPath)
	if err != nil {
		t.Fatalf("resolveExplicitPath returned error: %v", err)
	}

	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}
}

func TestResolveExplicitPath_NonExistent(t *testing.T) {
	_, err := resolveExplicitPath("/non/existent/path")
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
}

func TestResolveTargetPath_NoFlags(t *testing.T) {
	// Reset all flags
	componentFlag = ""
	baseFlag = ""
	projectFlag = ""
	pathFlag = ""

	_, err := resolveTargetPath()
	if err == nil {
		t.Error("expected error when no flags are set")
	}
}

func TestResolveTargetPath_MultipleFlags(t *testing.T) {
	componentFlag = "storage"
	baseFlag = "k8s"
	projectFlag = ""
	pathFlag = ""

	_, err := resolveTargetPath()
	if err == nil {
		t.Error("expected error when multiple flags are set")
	}

	// Reset
	componentFlag = ""
	baseFlag = ""
}

func TestResolveTargetPath_PathMutuallyExclusive(t *testing.T) {
	componentFlag = "storage"
	baseFlag = ""
	projectFlag = ""
	pathFlag = "/some/path"

	_, err := resolveTargetPath()
	if err == nil {
		t.Error("expected error when path is combined with other flags")
	}

	// Reset
	componentFlag = ""
	pathFlag = ""
}

func TestResolveTargetPath_WithExplicitPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test directory
	testPath := filepath.Join(tmpDir, "my-module")
	if err := os.MkdirAll(testPath, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	componentFlag = ""
	baseFlag = ""
	projectFlag = ""
	pathFlag = testPath

	result, err := resolveTargetPath()
	if err != nil {
		t.Fatalf("resolveTargetPath returned error: %v", err)
	}

	if result != testPath {
		t.Errorf("expected '%s', got '%s'", testPath, result)
	}

	// Reset
	pathFlag = ""
}

func TestFindModule_ComponentFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Set up cfg to point to tmpDir
	cfg = &config.Config{
		Root:   "",
		Binary: "terraform",
	}

	// Create components directory with a module
	componentsDir := filepath.Join(tmpDir, "components")
	modulePath := filepath.Join(componentsDir, "azurerm", "storage-account")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	// Create .tf file
	tfFile := filepath.Join(modulePath, "main.tf")
	if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
		t.Fatalf("failed to create .tf file: %v", err)
	}

	// Change to tmpDir
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	result, err := findModule("components", "storage-account")
	if err != nil {
		t.Fatalf("findModule returned error: %v", err)
	}

	if result != modulePath {
		t.Errorf("expected '%s', got '%s'", modulePath, result)
	}
}

func TestFindModule_ModuleNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	cfg = &config.Config{
		Root:   "",
		Binary: "terraform",
	}

	// Create components directory without the module we're looking for
	componentsDir := filepath.Join(tmpDir, "components")
	if err := os.MkdirAll(componentsDir, 0755); err != nil {
		t.Fatalf("failed to create components directory: %v", err)
	}

	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err := findModule("components", "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent module")
	}
}

func TestFindModule_DirectoryNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	cfg = &config.Config{
		Root:   "",
		Binary: "terraform",
	}

	// Don't create the components directory
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err := findModule("components", "any-module")
	if err == nil {
		t.Error("expected error when directory does not exist")
	}
}

func TestFindModule_WithConfigRoot(t *testing.T) {
	tmpDir := t.TempDir()

	// Set cfg.Root to a subdirectory
	cfg = &config.Config{
		Root:   "iac",
		Binary: "terraform",
	}

	// Create iac/components with a module
	modulePath := filepath.Join(tmpDir, "iac", "components", "storage-account")
	if err := os.MkdirAll(modulePath, 0755); err != nil {
		t.Fatalf("failed to create module directory: %v", err)
	}

	tfFile := filepath.Join(modulePath, "main.tf")
	if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
		t.Fatalf("failed to create .tf file: %v", err)
	}

	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	result, err := findModule("components", "storage-account")
	if err != nil {
		t.Fatalf("findModule returned error: %v", err)
	}

	if result != modulePath {
		t.Errorf("expected '%s', got '%s'", modulePath, result)
	}
}

func TestFindModule_NameClash(t *testing.T) {
	tmpDir := t.TempDir()

	cfg = &config.Config{
		Root:   "",
		Binary: "terraform",
	}

	// Create two modules with the same name
	module1 := filepath.Join(tmpDir, "components", "azurerm", "storage-account")
	module2 := filepath.Join(tmpDir, "components", "aws", "storage-account")

	for _, path := range []string{module1, module2} {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create module directory: %v", err)
		}
		tfFile := filepath.Join(path, "main.tf")
		if err := os.WriteFile(tfFile, []byte("# terraform"), 0644); err != nil {
			t.Fatalf("failed to create .tf file: %v", err)
		}
	}

	originalWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	_, err := findModule("components", "storage-account")
	if err == nil {
		t.Error("expected error for name clash")
	}
}
