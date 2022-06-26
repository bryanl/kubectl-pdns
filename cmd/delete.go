package cmd

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/kubectl-pdns/pkg/pdns"
)

func createDeleteCmd() *cobra.Command {
	var namespace string

	cmd := &cobra.Command{
		Use:   "delete <zone> <name> <rrType>",
		Short: "Delete a DNS value",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			config := pdns.DeleteConfig{
				Zone:      args[0],
				Name:      args[1],
				Type:      args[2],
				Namespace: "",
			}

			return pdns.Delete(config)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace")

	return cmd
}
