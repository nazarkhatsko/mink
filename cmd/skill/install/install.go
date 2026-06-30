package install

import (
	"github.com/nazarkhatsko/mink/internal/console"
	skill "github.com/nazarkhatsko/mink/internal/skill"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var forClaude, forCodex, global bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install mink skills",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			svc, err := skill.NewFromFlags(forClaude, forCodex, global)
			if err != nil {
				return err
			}
			names, err := svc.Names()
			if err != nil {
				return err
			}
			for _, name := range names {
				installed, err := svc.Install(name)
				if err != nil {
					return err
				}
				if installed {
					console.Printlnf("✓ installed /%s", name)
				} else {
					console.Printlnf("  /%s already installed", name)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&forClaude, "claude", false, "target Claude Code")
	cmd.Flags().BoolVar(&forCodex, "codex", false, "target Codex")
	cmd.Flags().BoolVar(&global, "global", false, "install to ~/.<platform>/skills/")
	return cmd
}
