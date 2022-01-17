package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Compute struct {
	context.Context
	aws.Config
}

func New(ctx context.Context, cfg aws.Config) *Compute {
	return &Compute{ctx, cfg}
}

type ComputeResources struct {
	CPU    int
	Memory int
}

type ComputeOutput struct {
	CPU       int
	Memory    int
	Type      string
	Count     int
	Resources map[string]int
}

const (
	MaxItems = 50
)

func (cmp Compute) ComputeResourcesByType(resource string) (*ComputeOutput, error) {
	switch resource {
	case "ec2":
		return ec2Analyzer(cmp.Context, cmp.Config)
	case "lambda":
		return lambdaAnalyzer(cmp.Context, cmp.Config)
	case "rds":
		return rdsAnalyzer(cmp.Context, cmp.Config)
	case "redshift":
		return redshiftAnalyzer(cmp.Context, cmp.Config)
	case "kafka":
		return kafkaAnalyzer(cmp.Context, cmp.Config)
	case "opensearch":
		return opensearchAnalyzer(cmp.Context, cmp.Config)
	default:
		return nil, fmt.Errorf("resource type is not (yet) supported: %s", resource)
	}
}
