package gettingstarted

import (
	"github.com/nazarkhatsko/mink/internal/manual"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "getting-started",
		Short: "Show getting started documentation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return manual.Show("getting-started")
		},
	}
}
