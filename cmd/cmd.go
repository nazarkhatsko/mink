package cmd

import (
	"github.com/nazarkhatsko/mink/cmd/doc"
	"github.com/nazarkhatsko/mink/cmd/run"
	"github.com/nazarkhatsko/mink/cmd/skill"
	"github.com/nazarkhatsko/mink/cmd/validate"
	"github.com/nazarkhatsko/mink/cmd/version"
	driverclaude "github.com/nazarkhatsko/mink/internal/drivers/claude"
	"github.com/nazarkhatsko/mink/internal/drivers/generate"
	driverhttp "github.com/nazarkhatsko/mink/internal/drivers/http"
	driversleep "github.com/nazarkhatsko/mink/internal/drivers/sleep"
	drivervalidate "github.com/nazarkhatsko/mink/internal/drivers/validate"
	"github.com/nazarkhatsko/mink/pkg/driver"
	"github.com/spf13/cobra"
)

func init() {
	driver.Register(driverhttp.New())
	driver.Register(generate.New())
	driver.Register(drivervalidate.New())
	driver.Register(driversleep.New())
	driver.Register(driverclaude.New())
}

func Execute() error {
	root := &cobra.Command{
		Use:   "mink",
		Short: "E2E backend testing tool",
	}
	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		run.NewCmd(),
		validate.NewCmd(),
		doc.NewCmd(),
		skill.NewCmd(),
		version.NewCmd(),
	)

	return root.Execute()
}
