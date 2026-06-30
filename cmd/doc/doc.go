package doc

import (
	"github.com/nazarkhatsko/mink/cmd/doc/drivers"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Show documentation",
	}
	cmd.AddCommand(drivers.NewCmd())
	return cmd
}
