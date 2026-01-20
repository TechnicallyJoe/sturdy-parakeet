package terraform

import (
	"fmt"
	"os"
	"os/exec"
)

// Binary is the terraform/tofu binary to use
var Binary = "terraform"

// SetBinary sets the binary to use (terraform or tofu)
func SetBinary(binary string) {
	Binary = binary
}

// RunInit runs terraform/tofu init in the specified directory
func RunInit(dir string) error {
	cmd := exec.Command(Binary, "init")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s init in %s\n", Binary, dir)
	return cmd.Run()
}

// RunFmt runs terraform/tofu fmt in the specified directory
func RunFmt(dir string) error {
	cmd := exec.Command(Binary, "fmt")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s fmt in %s\n", Binary, dir)
	return cmd.Run()
}

// RunValidate runs terraform/tofu validate in the specified directory
func RunValidate(dir string) error {
	cmd := exec.Command(Binary, "validate")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s validate in %s\n", Binary, dir)
	return cmd.Run()
}
