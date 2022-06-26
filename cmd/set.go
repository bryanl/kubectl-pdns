package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/kubectl-pdns/pkg/pdns"
)

func createSetCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "set <zone> <name> <rrType> <contents>",
		Short: "Set a DNS value",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 4 {
				return cmd.Help()
			}

			config := pdns.SetConfig{
				Zone:        args[0],
				Name:        args[1],
				Type:        args[2],
				RawContents: args[3],
				Namespace:   "",
			}

			return pdns.Set(config)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace")

	return cmd
}
