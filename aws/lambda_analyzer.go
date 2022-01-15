package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

func lambdaAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	fmt.Println("Lambda Analyzer")

	svc := lambda.NewFromConfig(cfg)
	p := &lambda.ListFunctionsInput{MaxItems: aws.Int32(MaxItems)}
	l := lambda.NewListFunctionsPaginator(svc, p)

	memCombined := 0
	count := 0

	for l.HasMorePages() {
		resp, err := l.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.Functions {
			memCombined += int(*i.MemorySize)
			count++
		}
	}

	return &ComputeOutput{
		CPU:    0,
		Memory: memCombined / 1024,
		Type:   "lambda",
		Count:  count,
	}, nil
}