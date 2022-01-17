package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
)

func kafkaAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := kafka.NewFromConfig(cfg)
	p := kafka.NewListClustersV2Paginator(svc, &kafka.ListClustersV2Input{})

	iCPU, iMemory := 0, 0
	iMap := map[string]int{}
	instances := 0

	for p.HasMorePages() {
		resp, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.ClusterInfoList {
			if i.Serverless != nil {
				return nil, fmt.Errorf("serverless kafka not yet supported")
			}
			iType := fmt.Sprintf("%v", *i.Provisioned.BrokerNodeGroupInfo.InstanceType)
			var numNodes int
			numNodes = int(i.Provisioned.NumberOfBrokerNodes)
			instances += numNodes
			iMap[iType]++
			info, err := kafkaOfferings(iType)
			if err != nil {
				return nil, err
			}
			iCPU += info.CPU * numNodes
			iMemory += info.Memory * numNodes
		}
	}

	return &ComputeOutput{
		CPU:       iCPU,
		Memory:    iMemory,
		Type:      "kafka",
		Count:     instances,
		Resources: iMap,
	}, nil
}

func kafkaOfferings(key string) (*ComputeResources, error) {
	result := map[string]ComputeResources{
		"kafka.t3.small":    {CPU: 2, Memory: 2},
		"kafka.m5.large":    {CPU: 2, Memory: 8},
		"kafka.m5.xlarge":   {CPU: 4, Memory: 16},
		"kafka.m5.2xlarge":  {CPU: 8, Memory: 32},
		"kafka.m5.4xlarge":  {CPU: 16, Memory: 64},
		"kafka.m5.8xlarge":  {CPU: 32, Memory: 128},
		"kafka.m5.12xlarge": {CPU: 48, Memory: 192},
		"kafka.m5.16xlarge": {CPU: 64, Memory: 256},
		"kafka.m5.24xlarge": {CPU: 96, Memory: 394},
	}
	value, isPresent := result[key]
	if isPresent {
		return &value, nil
	}
	return nil, fmt.Errorf("type %v not supported", key)
}
