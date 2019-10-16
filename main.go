package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	log "github.com/sirupsen/logrus"
)

var thisSession *session.Session

func init() {
	// Create a Session with a custom region
	thisSession = session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
}

func main() {

	// test if configuration works
	creds, err := thisSession.Config.Credentials.Get()
	if err != nil {
		log.Fatalf("Unable to setup session due to error: %v", err.Error())
	} else {
		log.Info(creds)
	}

	manageVPC(thisSession)

}
