package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
	log "github.com/sirupsen/logrus"
)

func manageMsk(thisSession *session.Session) {

	mskSvc := kafka.New(thisSession, aws.NewConfig().WithRegion("eu-west-1"))
	input := &kafka.ListClustersInput{}

	clusters, err := mskSvc.ListClusters(input)
	if err == nil {
		log.Info(clusters)
	}

	//clusterInput := &kafka.CreateClusterInput{}
	//var clusterInput *kafka.CreateClusterInput
	//clusterInput = new(kafka.CreateClusterInput)
	clusterInput := &kafka.CreateClusterInput{}

	// clusterInput.SetBrokerNodeGroupInfo(&kafka.BrokerNodeGroupInfo{
	//   ClientSubnets: ,
	//   InstanceType: ,
	// })
	clusterInput.SetClusterName("glp1.0")
	clusterInput.SetKafkaVersion("2.2.1")
	//clusterInput.NumberOfBrokerNodes(3)

	//createInput := &kafka.CreateClusterInput{}
	op, error := mskSvc.CreateCluster(clusterInput)
	if error != nil {
		log.Fatal(error.Error())
	}

	fmt.Println(op)

}
