package terraform

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/TechnicallyJoe/sturdy-parakeet/internal/config"
)

// Runner executes terraform/tofu commands using configuration
type Runner struct {
	config *config.Config
}

// NewRunner creates a new Runner with the given configuration
func NewRunner(cfg *config.Config) *Runner {
	return &Runner{config: cfg}
}

// Binary returns the configured binary name
func (r *Runner) Binary() string {
	return r.config.Binary
}

// RunInit executes terraform/tofu init in the specified directory
func (r *Runner) RunInit(dir string) error {
	cmd := exec.Command(r.config.Binary, "init")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s init in %s\n", r.config.Binary, dir)
	return cmd.Run()
}

// RunFmt executes terraform/tofu fmt in the specified directory
func (r *Runner) RunFmt(dir string) error {
	cmd := exec.Command(r.config.Binary, "fmt")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s fmt in %s\n", r.config.Binary, dir)
	return cmd.Run()
}

// RunValidate executes terraform/tofu validate in the specified directory
func (r *Runner) RunValidate(dir string) error {
	cmd := exec.Command(r.config.Binary, "validate")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s validate in %s\n", r.config.Binary, dir)
	return cmd.Run()
}
