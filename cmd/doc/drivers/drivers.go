package drivers

import (
	"sort"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drivers [name]",
		Short: "List drivers or show driver documentation",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				names := driver.List()
				sort.Strings(names)
				console.Printlnf("%-12s  %s", "NAME", "DESCRIPTION")
				for _, name := range names {
					d, _ := driver.Get(name)
					console.Printlnf("%-12s  %s", name, d.Describe().Description)
				}
				return nil
			}

			d, err := driver.Get(args[0])
			if err != nil {
				return err
			}
			doc := d.Describe()

			console.Printlnf("Driver: %s", d.Name())
			console.Printlnf("%s", doc.Description)

			if len(doc.Config) > 0 {
				console.Println("\nConfig:")
				for _, f := range doc.Config {
					req := "no"
					if f.Required {
						req = "yes"
					}
					console.Printlnf("  %-12s  %-8s  %-4s  %s", f.Name, f.Type, req, f.Description)
				}
			}

			if len(doc.Options) > 0 {
				console.Println("\nOptions:")
				for _, f := range doc.Options {
					req := "no"
					if f.Required {
						req = "yes"
					}
					console.Printlnf("  %-10s  %-8s  %-4s  %s", f.Name, f.Type, req, f.Description)
				}
			}

			if len(doc.Output) > 0 {
				console.Println("\nOutput:")
				for _, f := range doc.Output {
					console.Printlnf("  %-10s  %-8s  %s", f.Name, f.Type, f.Description)
				}
			}
			return nil
		},
	}
}
