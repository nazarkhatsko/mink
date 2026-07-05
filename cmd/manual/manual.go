package manual

import (
	"github.com/nazarkhatsko/mink/cmd/manual/configuration"
	"github.com/nazarkhatsko/mink/cmd/manual/drivers"
	"github.com/nazarkhatsko/mink/cmd/manual/gettingstarted"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manual",
		Short: "Show documentation",
	}
	cmd.AddCommand(drivers.NewCmd())
	cmd.AddCommand(configuration.NewCmd())
	cmd.AddCommand(gettingstarted.NewCmd())
	return cmd
}
