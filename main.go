package main

import (
	"os"

	"github.com/nazarkhatsko/mink/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
