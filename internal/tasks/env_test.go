package tasks

import (
	"os"
	"strings"
	"testing"
)

func TestNewEnvBuilder(t *testing.T) {
	b := NewEnvBuilder()
	if b == nil {
		t.Fatal("NewEnvBuilder should not return nil")
	}
	if b.vars == nil {
		t.Error("vars should be initialized")
	}
}

func TestEnvBuilder_WithGitRoot(t *testing.T) {
	b := NewEnvBuilder().WithGitRoot("/path/to/repo")
	if b.vars[EnvGitRoot] != "/path/to/repo" {
		t.Errorf("expected %q, got %q", "/path/to/repo", b.vars[EnvGitRoot])
	}
}

func TestEnvBuilder_WithModulePath(t *testing.T) {
	b := NewEnvBuilder().WithModulePath("/path/to/module")
	if b.vars[EnvModulePath] != "/path/to/module" {
		t.Errorf("expected %q, got %q", "/path/to/module", b.vars[EnvModulePath])
	}
}

func TestEnvBuilder_WithModuleName(t *testing.T) {
	b := NewEnvBuilder().WithModuleName("my-module")
	if b.vars[EnvModuleName] != "my-module" {
		t.Errorf("expected %q, got %q", "my-module", b.vars[EnvModuleName])
	}
}

func TestEnvBuilder_WithConfigPath(t *testing.T) {
	b := NewEnvBuilder().WithConfigPath("/path/to/.motf.yml")
	if b.vars[EnvConfigPath] != "/path/to/.motf.yml" {
		t.Errorf("expected %q, got %q", "/path/to/.motf.yml", b.vars[EnvConfigPath])
	}
}

func TestEnvBuilder_WithBinary(t *testing.T) {
	b := NewEnvBuilder().WithBinary("tofu")
	if b.vars[EnvBinary] != "tofu" {
		t.Errorf("expected %q, got %q", "tofu", b.vars[EnvBinary])
	}
}

func TestEnvBuilder_Chaining(t *testing.T) {
	b := NewEnvBuilder().
		WithGitRoot("/repo").
		WithModulePath("/repo/components/my-module").
		WithModuleName("my-module").
		WithConfigPath("/repo/.motf.yml").
		WithBinary("terraform")

	if b.vars[EnvGitRoot] != "/repo" {
		t.Errorf("EnvGitRoot = %q, want %q", b.vars[EnvGitRoot], "/repo")
	}
	if b.vars[EnvModulePath] != "/repo/components/my-module" {
		t.Errorf("EnvModulePath = %q, want %q", b.vars[EnvModulePath], "/repo/components/my-module")
	}
	if b.vars[EnvModuleName] != "my-module" {
		t.Errorf("EnvModuleName = %q, want %q", b.vars[EnvModuleName], "my-module")
	}
	if b.vars[EnvConfigPath] != "/repo/.motf.yml" {
		t.Errorf("EnvConfigPath = %q, want %q", b.vars[EnvConfigPath], "/repo/.motf.yml")
	}
	if b.vars[EnvBinary] != "terraform" {
		t.Errorf("EnvBinary = %q, want %q", b.vars[EnvBinary], "terraform")
	}
}

func TestEnvBuilder_Build_IncludesParentEnv(t *testing.T) {
	// Set a test env var
	os.Setenv("MOTF_TEST_VAR", "test_value")
	defer os.Unsetenv("MOTF_TEST_VAR")

	env := NewEnvBuilder().Build()

	found := false
	for _, e := range env {
		if e == "MOTF_TEST_VAR=test_value" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Build() should include parent process environment")
	}
}

func TestEnvBuilder_Build_IncludesBuiltInVars(t *testing.T) {
	env := NewEnvBuilder().
		WithGitRoot("/repo").
		WithModuleName("test-module").
		Build()

	foundGitRoot := false
	foundModuleName := false
	for _, e := range env {
		if e == "MOTF_GIT_ROOT=/repo" {
			foundGitRoot = true
		}
		if e == "MOTF_MODULE_NAME=test-module" {
			foundModuleName = true
		}
	}

	if !foundGitRoot {
		t.Error("Build() should include MOTF_GIT_ROOT")
	}
	if !foundModuleName {
		t.Error("Build() should include MOTF_MODULE_NAME")
	}
}

func TestEnvBuilder_Build_EmptyValues(t *testing.T) {
	// Empty values should still be included (for soft-fail scenarios)
	env := NewEnvBuilder().WithGitRoot("").Build()

	found := false
	for _, e := range env {
		if strings.HasPrefix(e, "MOTF_GIT_ROOT=") {
			found = true
			if e != "MOTF_GIT_ROOT=" {
				t.Errorf("expected empty value, got %q", e)
			}
			break
		}
	}
	if !found {
		t.Error("Build() should include empty MOTF_GIT_ROOT")
	}
}

func TestModuleNameFromPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/path/to/storage-account", "storage-account"},
		{"/repo/components/azurerm/key-vault", "key-vault"},
		{"my-module", "my-module"},
		{"", ""},
		{"/", "/"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := ModuleNameFromPath(tt.path)
			if got != tt.want {
				t.Errorf("ModuleNameFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
