package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func rdsAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := rds.NewFromConfig(cfg)
	l := rds.NewDescribeDBInstancesPaginator(svc, &rds.DescribeDBInstancesInput{MaxRecords: aws.Int32(100)})

	count := 0

	iMap := map[string]int{}
	iCPU := 0
	iMemory := 0

	for l.HasMorePages() {
		resp, err := l.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.DBInstances {
			iType := fmt.Sprintf("%v", *i.DBInstanceClass)
			iMap[iType]++
			count++
			info, err := rdsOfferings(iType)
			if err != nil {
				return nil, err
			}
			iCPU += info.CPU
			iMemory += info.Memory
		}
	}
	return &ComputeOutput{
		CPU:    iCPU,
		Memory: iMemory,
		Type:   "rds",
		Count:  count,
	}, nil
}

func rdsOfferings(key string) (*ComputeResources, error) {
	result := map[string]ComputeResources{
		"db.t3.micro":   {CPU: 2, Memory: 1},
		"db.t3.small":   {CPU: 2, Memory: 2},
		"db.t3.medium":  {CPU: 2, Memory: 4},
		"db.t3.large":   {CPU: 2, Memory: 8},
		"db.t3.xlarge":  {CPU: 4, Memory: 16},
		"db.t3.2xlarge": {CPU: 8, Memory: 32},
		"db.t4g.small":  {CPU: 2, Memory: 2},
		"db.t4g.medium": {CPU: 2, Memory: 4},
	}
	value, isPresent := result[key]
	if isPresent {
		return &value, nil
	}
	return nil, fmt.Errorf("type %v not supported", key)
}
