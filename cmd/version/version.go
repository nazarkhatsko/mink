package version

import (
	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/spf13/cobra"
)

var Version = "v1.1.0"

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print mink version",
		Run: func(cmd *cobra.Command, args []string) {
			console.Printlnf("mink %s", Version)
		},
	}
}
