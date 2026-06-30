package skill

import (
	"github.com/nazarkhatsko/mink/cmd/skill/install"
	"github.com/nazarkhatsko/mink/cmd/skill/uninstall"
	"github.com/nazarkhatsko/mink/cmd/skill/upgrade"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Manage mink skills",
	}
	cmd.AddCommand(install.NewCmd())
	cmd.AddCommand(uninstall.NewCmd())
	cmd.AddCommand(upgrade.NewCmd())
	return cmd
}
