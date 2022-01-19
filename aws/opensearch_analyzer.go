package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

func opensearchAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := opensearch.NewFromConfig(cfg)

	ldn, err := svc.ListDomainNames(ctx, nil)
	if err != nil {
		return nil, err
	}
	iCPU, iMemory, instances := 0, 0, 0
	iMap := map[string]int{}

	for _, domain := range ldn.DomainNames {
		dm, err := svc.DescribeDomain(ctx, &opensearch.DescribeDomainInput{DomainName: domain.DomainName})
		if err != nil {
			return nil, err
		}
		iType := fmt.Sprintf("%v", dm.DomainStatus.ClusterConfig.InstanceType)
		info, err := searchOfferings(iType)
		if err != nil {
			return nil, err
		}
		var numNodes int
		numNodes = int(*dm.DomainStatus.ClusterConfig.InstanceCount)
		instances += numNodes
		iMap[iType]++

		iCPU += info.CPU * numNodes
		iMemory += info.Memory * numNodes
	}

	return &ComputeOutput{
		CPU:       iCPU,
		Memory:    iMemory,
		Type:      "opensearch",
		Count:     instances,
		Resources: iMap,
	}, nil
}

func searchOfferings(key string) (*ComputeResources, error) {
	result := map[string]ComputeResources{
		"t3.small.search":    {CPU: 2, Memory: 2},
		"t3.medium.search":   {CPU: 2, Memory: 4},
		"m6g.large.search":   {CPU: 2, Memory: 8},
		"m6g.xlarge.search":  {CPU: 4, Memory: 16},
		"m6g.2xlarge.search": {CPU: 8, Memory: 32},
		"m6g.4xlarge.search": {CPU: 16, Memory: 64},
		"r5.large.search":    {CPU: 2, Memory: 16},
		"r5.xlarge.search":   {CPU: 4, Memory: 32},
		"m5.large.search":    {CPU: 2, Memory: 8},
	}
	value, isPresent := result[key]
	if isPresent {
		return &value, nil
	}
	return nil, fmt.Errorf("type %v not supported. please add from https://aws.amazon.com/opensearch-service/pricing/", key)
}
