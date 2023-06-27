

package aws_test

import (
	"context"
	"log"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	awstrace "git.proto.group/protoobp/pobp-trace-go/contrib/aws/aws-sdk-go-v2/aws"
)

func Example() {
	awsCfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}

	awstrace.AppendMiddleware(&awsCfg)

	sqsClient := sqs.NewFromConfig(awsCfg)
	sqsClient.ListQueues(context.Background(), &sqs.ListQueuesInput{})
}
