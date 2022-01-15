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
	return &Compute{
		ctx,
		cfg,
	}
}

type ComputeResources struct {
	CPU    int
	Memory int
}

type ComputeOutput struct {
	CPU    int
	Memory int
	Type   string
}

func (cmp Compute) ComputeResourcesByType(resource string) (*ComputeOutput, error) {
	switch resource {
	case "ec2":
		return ec2Analyzer(cmp.Context, cmp.Config)
	default:
		return nil, fmt.Errorf("resource type is not (yet) supported: %s", resource)
	}
}
