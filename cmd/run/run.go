package run

import (
	"context"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/internal/engine"
	"github.com/nazarkhatsko/mink/internal/envfile"
	"github.com/nazarkhatsko/mink/internal/report"
	"github.com/spf13/cobra"
)

func newReporter(name string) report.Reporter {
	switch name {
	case "compact":
		return report.NewCompact()
	case "json":
		return report.NewJSON()
	case "silent":
		return report.NewSilent()
	default:
		return report.NewPretty()
	}
}

func NewCmd() *cobra.Command {
	var flowID, reporterName, envFile string

	cmd := &cobra.Command{
		Use:   "run <file>",
		Short: "Run flows from a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if envFile != "" {
				if err := envfile.Load(envFile); err != nil {
					return err
				}
			}
			cfg, err := config.Load(args[0])
			if err != nil {
				return err
			}
			eng := engine.New(cfg, newReporter(reporterName))
			return eng.Run(context.Background(), flowID)
		},
	}

	cmd.Flags().StringVar(&flowID, "flow-id", "", "Run only a specific flow by id")
	cmd.Flags().StringVar(&reporterName, "reporter", "pretty", "Output format: pretty|compact|json|silent")
	cmd.Flags().StringVar(&envFile, "env", "", "Path to .env file")

	return cmd
}
