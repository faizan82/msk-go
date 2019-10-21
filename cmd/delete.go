package main

import "github.com/spf13/cobra"

var (
	clusterName string
	deleteCmd   = &cobra.Command{
		Use:   "delete",
		Short: "delete removes an msk cluster and associated resources",
		Long:  "delete submits a request for msk cluster deletion and dependent resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	deleteCmd.PersistentFlags().StringVar(&clusterName, "clusterName", "", "Name of the cluster to be deleted")

}
