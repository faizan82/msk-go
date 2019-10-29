package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
)

var (
	thisSession *session.Session
	rootCmd     = &cobra.Command{
		Use:   "mskmanger",
		Short: "mskmanager is a controller for operating Amazon's managed kafka service",
		Long:  `Run as an operator or as using CLI to manage Amazon's managed kafka service`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
			cmd.Help()
		},
	}
)

// Execute calls method on Command data and returns error object
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	thisSession = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
}

func main() {

	// test if configuration works
	// creds, err := thisSession.Config.Credentials.Get()
	// if err != nil {
	// 	log.Fatalf("Unable to setup session due to error: %v", err.Error())
	// } else {
	// 	log.Info(creds)
	// }

	//manageVPC(thisSession)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.Execute()

}
