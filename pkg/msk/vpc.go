package aws

import (
	"context"
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// ManageVPC creates, deletes and operates on AWS VPC
func ManageVPC(thisSession *session.Session, vpcID *string, region string, cidrBlock net.IPNet, name string) {
	ec2Session := ec2.New(thisSession, aws.NewConfig().WithRegion(region))
	ctx := context.Background()

	// if no vpc-id is provided, we create a new vpc
	if *vpcID == "" {
		vpcID = createVpc(ctx, ec2Session, cidrBlock, name)
		log.Infof("Created VPC's ID: %v,", vpcID)
		createSubnets(ctx, ec2Session, vpcID, region, &cidrBlock, name)
	} else {
		GetSubnets(ctx, ec2Session, vpcID)
	}

}

func createVpc(ctx context.Context, ec2Session *ec2.EC2, cidrBlock net.IPNet, name string) *string {
	log.Info("Creating VPC")
	vpcInput := &ec2.CreateVpcInput{}
	vpcInput.SetCidrBlock(cidrBlock.String())
	vpcInput.SetDryRun(false)
	fmt.Println(vpcInput.GoString())

	vpcerr := vpcInput.Validate()
	if vpcerr != nil {
		log.Error(vpcerr)
	}

	vpcOutPut, err := ec2Session.CreateVpcWithContext(ctx, vpcInput)

	if err != nil {
		log.Fatal(err)
	}

	log.Info(vpcOutPut)

	tagInput := &ec2.CreateTagsInput{}
	tagInput.SetResources([]*string{vpcOutPut.Vpc.VpcId})
	tagInput.SetTags([]*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String(name),
		},
	})
	tagInput.SetDryRun(false)

	_, tagerr := ec2Session.CreateTagsWithContext(ctx, tagInput)
	if tagerr != nil {
		log.Info(tagerr.Error())
	}

	return vpcOutPut.Vpc.VpcId
}

func createSubnets(ctx context.Context, ec2Client *ec2.EC2, VpcID *string, region string, cidrBlock *net.IPNet, name string) {
	azs := getAZS(ctx, ec2Client)
	azCount := len(azs)
	log.Infof("Number of AZs in region: %v", azCount)
	// if we have 2 AZs in a region, we will create two subnets in one of the AZs
	// else we equally distribute subnets acrorss AZs
	switch {
	case azCount <= 2:
		log.Info("2 AZs")
	case azCount >= 3:
		log.Info("More than 2 AZs")
	}

	calculatedSubnet, err := cidr.Subnet(cidrBlock, 4, 0)
	if err == nil {
		log.Infof("calculatedSubnet: %v", calculatedSubnet)
	}

	func() {
		for i := 0; i < azCount; i++ {
			subnetInput := &ec2.CreateSubnetInput{}
			subnetInput.SetAvailabilityZoneId(*azs[i].ZoneId)
			subnetInput.SetVpcId(*VpcID)
			subnetInput.SetDryRun(false)
			calculatedSubnet, err := cidr.Subnet(cidrBlock, 4, i)
			if err != nil {
				log.Fatalf("Error calculating subnet range: %v", err)
			}
			subnetInput.SetCidrBlock(calculatedSubnet.String())
			subnetInput.Validate()
			err = subnetInput.Validate()
			if err != nil {
				log.Error(err.Error())
			}
			log.Infof("Subnet input: %v", subnetInput)
			subnetOutput, err := ec2Client.CreateSubnetWithContext(ctx, subnetInput)
			if err != nil {
				log.Fatalf("Unable to create subnets due to error: %v", err)
			}
			log.Infof("Subnet output: %v", subnetOutput)

			tagInput := &ec2.CreateTagsInput{}
			tagInput.SetResources([]*string{subnetOutput.Subnet.SubnetId})
			tagInput.SetTags([]*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(name),
				},
			})
			tagInput.SetDryRun(false)
			_, tagerr := ec2Client.CreateTagsWithContext(ctx, tagInput)
			if tagerr != nil {
				log.Info(tagerr.Error())
			}

		}
	}()

}

func getAZS(ctx context.Context, ec2Client *ec2.EC2) []*ec2.AvailabilityZone {
	getAZInput := &ec2.DescribeAvailabilityZonesInput{}

	getAZOutput, err := ec2Client.DescribeAvailabilityZonesWithContext(ctx, getAZInput)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(getAZOutput.AvailabilityZones)
	return getAZOutput.AvailabilityZones

}

// GetSubnets returns a list of subnets for provided VPC
func GetSubnets(ctx context.Context, ec2Session *ec2.EC2, VpcID *string) {
	subnetInput := &ec2.DescribeSubnetsInput{}
	subnetFilter := []*ec2.Filter{
		{
			Name: aws.String("vpc-id"),
			Values: []*string{
				aws.String(*VpcID),
			},
		},
	}
	subnetInput.SetFilters(subnetFilter)
	subnetsOp, _ := ec2Session.DescribeSubnets(subnetInput)
	fmt.Println(subnetsOp)
}
