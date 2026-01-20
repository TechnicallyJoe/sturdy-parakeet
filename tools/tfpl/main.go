package main

import (
	"os"

	"github.com/TechnicallyJoe/sturdy-parakeet/tools/tfpl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
