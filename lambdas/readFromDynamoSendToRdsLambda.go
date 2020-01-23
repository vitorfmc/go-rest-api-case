package main

/*
    Developed by "https://github.com/vitorfmc"
    
    =======================================================
    Overview:
    =======================================================

    This Lambda Function is example of integration with RDS and DynamoDB.
    The idea is: Everytime a table in DynamoDB receive a data, it will send the event 
    information to this lambda function, which will persist in a Aurora RDS.

    DynamoDB Stream ==> This Lambda Function ==> AuroraDB

	Obs.: Remember to give SecretsManagerReadWrite, AWSLambdaDynamoDBExecutionRole and
	AmazonRDSDataFullAccess policies to your Lambda Function
*/

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rdsdataservice"
)

type Client struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(ctx context.Context, e events.DynamoDBEvent) {

	for _, record := range e.Records {
		if record.EventName == "INSERT" {
			
			var client *Client
			client = new(Client)
			client.Name = record.Change.NewImage["name"].String()
			client.Email = record.Change.NewImage["email"].String()

			sess, _ := session.NewSession(&aws.Config{
				Region: aws.String("us-east-1")},
			)
		
			rdsdataservice_client := rdsdataservice.New(sess)

			SQLStatement := "insert into client(name,email) values ('" + client.Name + "','" + client.Email + "');"
			fmt.Println("statement:", SQLStatement)
		
			req, resp := rdsdataservice_client.ExecuteStatementRequest(&rdsdataservice.ExecuteStatementInput{
				Database:    aws.String("client"),
				ResourceArn: aws.String("arn:aws:rds:us-east-1:595223775479:cluster:testdb"),
				SecretArn:   aws.String("arn:aws:secretsmanager:us-east-1:595223775479:secret:rds-db-credentials/cluster-DKYGVZQ7LDABC5HUTQERBFTWO4/superadmin-lALlHn"),
				Sql:         aws.String(SQLStatement),
			})
		
			err1 := req.Send()
			if err1 == nil { // resp is now filled
				fmt.Println("Response:", resp)
			} else {
				fmt.Println("error:", err1)
			}

		}else{
			fmt.Printf("Event: %s\n", record.EventName)
		}
	}

}