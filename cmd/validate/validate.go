package validate

import (
	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(args[0])
			if err != nil {
				return err
			}
			console.Println("✓ config is valid")
			return nil
		},
	}
}
