package main

import (
	"fmt"
	"net"

	msk "aws-code/msk-operator/pkg/msk"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	vpcNeeded   bool
	cidr        net.IPNet
	region      string
	vpcID       string
	brokerType  string
	brokerCount string
	name        string
	createCmd   = &cobra.Command{
		Use:   "create",
		Short: "create would create a msk cluster",
		Long:  "create submits a request for msk cluster creation based on parameters provided ",
		Args: func(cmd *cobra.Command, args []string) error {
			log.Infof("vpcneeded: %v, region: %v", vpcNeeded, region)
			if !vpcNeeded && &vpcID == nil {
				fmt.Println("Please provide vpc id or select vpc option to create a dedicated vpc")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			//region comes from rootCmd
			log.Info("Calling VPC manager")
			log.Infof("cidr: %v, network: %v", cidr, cidr.IP.DefaultMask())
			msk.ManageVPC(thisSession, &vpcID, region, cidr, name)

		},
	}
)

func init() {
	createCmd.Flags().BoolVarP(&vpcNeeded, "vpc", "v", false, "Flag that determines if a dedicated vpc is needed. If set a new VPC with three subnets and a Nat gateway would be provisioned")
	createCmd.Flags().StringVarP(&vpcID, "vpcID", "i", "", "Provide VPC ID if dedicated vpc is not opted for. Subnets will be created for msk cluster across available AZs")
	createCmd.Flags().StringVarP(&brokerType, "brokerType", "t", "m5.xlarge", "Define type of node machine needed for brokers.")
	createCmd.Flags().StringVarP(&brokerCount, "brokerCount", "c", "3", "Define number of brokers needed.")
	// @// TODO:  change default region to us-east-1 as usually that is the cheapest
	createCmd.Flags().StringVarP(&region, "region", "r", "eu-west-1", "Region where MSK cluster needs to be operated on.")
	createCmd.Flags().StringVarP(&name, "name", "n", "msk-cluster", "Name of msk cluster and subcomponents.")
	createCmd.Flags().IPNetVar(&cidr, "vpc-cidr", cidr, "Vpc cidr range")

}
