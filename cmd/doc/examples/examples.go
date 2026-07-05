package examples

import (
	"github.com/nazarkhatsko/mink/internal/docutil"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "examples [name]",
		Short: "List examples or show example documentation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			page := "README"
			if len(args) == 1 {
				page = args[0]
			}
			return docutil.Show("examples", page)
		},
	}
}
