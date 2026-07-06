package drivers

import (
	"github.com/nazarkhatsko/mink/internal/manual"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drivers [name]",
		Short: "List drivers or show driver documentation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return manual.Show("drivers", args[0])
			}
			return manual.Show("drivers")
		},
	}
}
