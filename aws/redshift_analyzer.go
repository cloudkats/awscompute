package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

func redshiftAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	// fmt.Println("RDS Analyzer")

	svc := redshift.NewFromConfig(cfg)
	p := &redshift.DescribeClustersInput{MaxRecords: aws.Int32(100)}
	l := redshift.NewDescribeClustersPaginator(svc, p)

	count := 0

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
			count++
			info := redshiftOfferings(iType)
			iCPU += info.CPU
			iMemory += info.Memory
		}
	}

	// fmt.Println(iMap)

	return &ComputeOutput{
		CPU:    iCPU,
		Memory: iMemory,
		Type:   "redshift",
		Count:  count,
	}, nil
}

func redshiftOfferings(key string) ComputeResources {
	result := map[string]ComputeResources{
		"ds2.xlarge":  ComputeResources{CPU: 4, Memory: 31},
		"ds2.8xlarge": ComputeResources{CPU: 36, Memory: 244},
		"dc2.large":   ComputeResources{CPU: 2, Memory: 15},
		"dc2.8xlarge": ComputeResources{CPU: 32, Memory: 244},
		"dc1.large":   ComputeResources{CPU: 2, Memory: 15},
		"dc1.8xlarge": ComputeResources{CPU: 32, Memory: 244},
	}
	return result[key]

}
