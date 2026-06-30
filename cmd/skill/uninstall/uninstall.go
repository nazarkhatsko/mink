package uninstall

import (
	"errors"

	"github.com/nazarkhatsko/mink/internal/console"
	skill "github.com/nazarkhatsko/mink/internal/skill"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var forClaude, forCodex, global bool
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall mink skills",
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
				if err := svc.Uninstall(name); err != nil {
					if errors.Is(err, skill.ErrNotInstalled) {
						console.Printlnf("  /%s not installed", name)
						continue
					}
					return err
				}
				console.Printlnf("✓ uninstalled /%s", name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&forClaude, "claude", false, "target Claude Code")
	cmd.Flags().BoolVar(&forCodex, "codex", false, "target Codex")
	cmd.Flags().BoolVar(&global, "global", false, "uninstall from ~/.<platform>/skills/")
	return cmd
}
