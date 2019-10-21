package main

import (
	msk "aws-code/msk-operator/pkg/msk"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	vpcNeeded   bool
	vpcID       string
	brokerType  string
	brokerCount string
	createCmd   = &cobra.Command{
		Use:   "create",
		Short: "create would create a msk cluster",
		Long:  "create submits a request for msk cluster creation based on parameters provided ",
		Args: func(cmd *cobra.Command, args []string) error {
			if !vpcNeeded && &vpcID == nil {
				fmt.Println("Please provide vpc id")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//cmd.Help()
			msk.ManageVPC(thisSession, &vpcID)

		},
	}
)

func init() {
	createCmd.Flags().BoolVarP(&vpcNeeded, "vpc", "v", false, "Flag that determines if a dedicated vpc is needed.")
	createCmd.Flags().StringVarP(&vpcID, "vpcID", "i", "", "Provide VPC ID if dedicated vpc is not opted for. Subnets will be choosen automatically")
	createCmd.Flags().StringVarP(&brokerType, "brokerType", "t", "m5.xlarge", "Define type of node machine needed for brokers.")
	createCmd.Flags().StringVarP(&brokerCount, "brokerCount", "c", "3", "Define number of brokers needed.")

}
