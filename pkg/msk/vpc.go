package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// ManageVPC creates, deletes and operates on AWS VPC
func ManageVPC(thisSession *session.Session, vpcID *string) {
	ec2Session := ec2.New(thisSession, aws.NewConfig().WithRegion("eu-west-1"))
	ctx := context.Background()

	//ec2Session.CreateTagsWithContext(ctx)
	if vpcID == nil {
		vpcID = createVpc(ctx, ec2Session)
		createSubnets(ctx, ec2Session, vpcID)
	} else {
		GetSubnets(ctx, ec2Session, vpcID)
	}

}

func createVpc(ctx context.Context, ec2Session *ec2.EC2) *string {
	//var vpcInput *ec2.CreateVpcInput
	//vpcInput = new(ec2.CreateVpcInput)
	vpcInput := &ec2.CreateVpcInput{}
	vpcInput.SetCidrBlock("10.200.0.0/16")
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
			Value: aws.String("gosdk-vpc"),
		},
	})
	tagInput.SetDryRun(false)

	_, tagerr := ec2Session.CreateTagsWithContext(ctx, tagInput)
	if tagerr != nil {
		log.Info(tagerr.Error())
	}

	return vpcOutPut.Vpc.VpcId
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

func createSubnets(ctx context.Context, ec2Client *ec2.EC2, VpcID *string) {
	subnetInput := &ec2.CreateSubnetInput{}
	subnetInput.SetAvailabilityZone("eu-west-1a")
	subnetInput.SetCidrBlock("10.200.10.0/24")
	subnetInput.SetVpcId(*VpcID)
	subnetInput.SetDryRun(false)
	subnetInput.Validate()
	err := subnetInput.Validate()
	if err != nil {
		log.Error(err.Error())
	}

	subnetOutput, err := ec2Client.CreateSubnetWithContext(ctx, subnetInput)

	tagInput := &ec2.CreateTagsInput{}
	tagInput.SetResources([]*string{subnetOutput.Subnet.SubnetId})
	tagInput.SetTags([]*ec2.Tag{
		{
			Key:   aws.String("Name"),
			Value: aws.String("gosdk-vpc"),
		},
	})
	tagInput.SetDryRun(false)

	_, tagerr := ec2Client.CreateTagsWithContext(ctx, tagInput)
	if tagerr != nil {
		log.Info(tagerr.Error())
	}

}
