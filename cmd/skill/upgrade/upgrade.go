package upgrade

import (
	"errors"

	"github.com/nazarkhatsko/mink/internal/console"
	skill "github.com/nazarkhatsko/mink/internal/skill"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var forClaude, forCodex, global bool
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade mink skills to the latest embedded version",
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
				upgraded, err := svc.Upgrade(name)
				if err != nil {
					if errors.Is(err, skill.ErrNotInstalled) {
						console.Printlnf("  /%s not installed", name)
						continue
					}
					return err
				}
				if upgraded {
					console.Printlnf("✓ upgraded /%s", name)
				} else {
					console.Printlnf("  /%s already up to date", name)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&forClaude, "claude", false, "target Claude Code")
	cmd.Flags().BoolVar(&forCodex, "codex", false, "target Codex")
	cmd.Flags().BoolVar(&global, "global", false, "upgrade in ~/.<platform>/skills/")
	return cmd
}
