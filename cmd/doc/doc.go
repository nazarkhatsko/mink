package doc

import (
	"github.com/nazarkhatsko/mink/cmd/doc/configuration"
	"github.com/nazarkhatsko/mink/cmd/doc/drivers"
	"github.com/nazarkhatsko/mink/cmd/doc/gettingstarted"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Show documentation",
	}
	cmd.AddCommand(drivers.NewCmd())
	cmd.AddCommand(configuration.NewCmd())
	cmd.AddCommand(gettingstarted.NewCmd())
	return cmd
}
