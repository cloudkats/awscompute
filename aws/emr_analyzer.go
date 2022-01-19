package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
)

func emrAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := emr.NewFromConfig(cfg)
	p, err := svc.ListClusters(ctx,
		&emr.ListClustersInput{ClusterStates: []types.ClusterState{types.ClusterStateRunning}})
	if err != nil {
		return nil, err
	}
	iCPU, iMemory, instances := 0, 0, 0
	iMap := map[string]int{}

	for _, cl := range p.Clusters {
		c, err := svc.DescribeCluster(ctx, &emr.DescribeClusterInput{ClusterId: cl.Id})
		if err != nil {
			return nil, err
		}
		instanceGroups := c.Cluster.Ec2InstanceAttributes
		fmt.Println(instanceGroups)

	}

	return &ComputeOutput{
		CPU:       iCPU,
		Memory:    iMemory,
		Type:      "emr",
		Count:     instances,
		Resources: iMap,
	}, nil
}
