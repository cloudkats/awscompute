package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

func redshiftAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := redshift.NewFromConfig(cfg)
	p := &redshift.DescribeClustersInput{MaxRecords: aws.Int32(100)}
	l := redshift.NewDescribeClustersPaginator(svc, p)

	instances := 0

	iMap := map[string]int{}
	iCPU := 0
	iMemory := 0

	for l.HasMorePages() {
		resp, err := l.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.Clusters {
			iType := fmt.Sprintf("%v", *i.NodeType)
			iMap[iType]++
			info, err := redshiftOfferings(iType)
			if err != nil {
				return nil, err
			}
			var numNodes int
			numNodes = int(i.NumberOfNodes)
			instances += numNodes
			iCPU += info.CPU * numNodes
			iMemory += info.Memory * numNodes
		}
	}

	return &ComputeOutput{
		CPU:       iCPU,
		Memory:    iMemory,
		Type:      "redshift",
		Count:     instances,
		Resources: iMap,
	}, nil
}

func redshiftOfferings(key string) (*ComputeResources, error) {
	result := map[string]ComputeResources{
		"ds2.xlarge":   {CPU: 4, Memory: 31},
		"ds2.8xlarge":  {CPU: 36, Memory: 244},
		"dc2.large":    {CPU: 2, Memory: 15},
		"dc2.8xlarge":  {CPU: 32, Memory: 244},
		"dc1.large":    {CPU: 2, Memory: 15},
		"dc1.8xlarge":  {CPU: 32, Memory: 244},
		"ra3.xlplus":   {CPU: 4, Memory: 32},
		"ra3.4xlarge":  {CPU: 12, Memory: 96},
		"ra3.16xlarge": {CPU: 48, Memory: 384},
	}
	value, isPresent := result[key]
	if isPresent {
		return &value, nil
	}
	return nil, fmt.Errorf("type %v not supported", key)
}
