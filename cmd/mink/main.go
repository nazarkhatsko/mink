package main

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/internal/envfile"
	"github.com/nazarkhatsko/mink/pkg/driver"
	driverhttp "github.com/nazarkhatsko/mink/internal/drivers/http"
	"github.com/nazarkhatsko/mink/internal/drivers/generate"
	"github.com/nazarkhatsko/mink/internal/drivers/sleep"
	"github.com/nazarkhatsko/mink/internal/drivers/validate"
	"github.com/nazarkhatsko/mink/internal/engine"
	"github.com/nazarkhatsko/mink/internal/report"
	"github.com/spf13/cobra"
)

func init() {
	driver.Register(driverhttp.New())
	driver.Register(generate.New())
	driver.Register(validate.New())
	driver.Register(sleep.New())
}

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

func main() {
	root := &cobra.Command{
		Use:   "mink",
		Short: "E2E backend testing tool",
	}

	var flowName string
	var reporterName string
	var envFile string

	runCmd := &cobra.Command{
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
			return eng.Run(context.Background(), flowName)
		},
	}
	runCmd.Flags().StringVar(&flowName, "flow", "", "Run only a specific flow by name")
	runCmd.Flags().StringVar(&reporterName, "reporter", "pretty", "Output format: pretty|compact|json|silent")
	runCmd.Flags().StringVar(&envFile, "env", "", "Path to .env file")

	validateCmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate a config file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config.Load(args[0])
			if err != nil {
				return err
			}
			fmt.Println("✓ config is valid")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list-drivers",
		Short: "List registered drivers",
		Run: func(cmd *cobra.Command, args []string) {
			names := driver.List()
			sort.Strings(names)
			for _, name := range names {
				fmt.Println(" -", name)
			}
		},
	}

	root.AddCommand(runCmd, validateCmd, listCmd)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
