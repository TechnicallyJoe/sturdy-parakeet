package terraform

import (
	"testing"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/config"
)

func TestNewRunner(t *testing.T) {
	cfg := &config.Config{
		Root:   "/some/path",
		Binary: "tofu",
	}

	runner := NewRunner(cfg)

	if runner == nil {
		t.Fatal("NewRunner returned nil")
	}

	if runner.config != cfg {
		t.Error("NewRunner did not store config correctly")
	}
}

func TestRunner_Binary(t *testing.T) {
	tests := []struct {
		name     string
		binary   string
		expected string
	}{
		{
			name:     "terraform binary",
			binary:   "terraform",
			expected: "terraform",
		},
		{
			name:     "tofu binary",
			binary:   "tofu",
			expected: "tofu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Binary: tt.binary,
			}

			runner := NewRunner(cfg)
			if runner.Binary() != tt.expected {
				t.Errorf("Binary() = %s, expected %s", runner.Binary(), tt.expected)
			}
		})
	}
}

func TestRunner_InheritsConfigBinary(t *testing.T) {
	// Test that Runner properly inherits the binary from config
	cfg := &config.Config{
		Root:   "/test/root",
		Binary: "tofu",
	}

	runner := NewRunner(cfg)

	// Binary should come from config
	if runner.Binary() != "tofu" {
		t.Errorf("expected Binary to be 'tofu', got '%s'", runner.Binary())
	}

	// Changing config should reflect in runner
	cfg.Binary = "terraform"
	if runner.Binary() != "terraform" {
		t.Errorf("expected Binary to be 'terraform' after config change, got '%s'", runner.Binary())
	}
}

func TestRunner_WithDefaultConfig(t *testing.T) {
	// Test that Runner works with default config values
	cfg := config.DefaultConfig()
	runner := NewRunner(cfg)

	if runner.Binary() != "terraform" {
		t.Errorf("expected default Binary to be 'terraform', got '%s'", runner.Binary())
	}
}
