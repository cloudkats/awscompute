package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func rdsAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	fmt.Println("RDS Analyzer")

	svc := rds.NewFromConfig(cfg)
	p := &rds.DescribeDBInstancesInput{MaxRecords: aws.Int32(100)}
	l := rds.NewDescribeDBInstancesPaginator(svc, p)

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
			info := rdsOfferings(iType)
			iCPU += info.CPU
			iMemory += info.Memory
		}
	}

	fmt.Println(iMap)

	return &ComputeOutput{
		CPU:    iCPU,
		Memory: iMemory,
		Type:   "rds",
		Count:  count,
	}, nil
}

func rdsOfferings(key string) ComputeResources {
	result := map[string]ComputeResources{
		"db.t3.micro":   ComputeResources{CPU: 2, Memory: 1},
		"db.t3.small":   ComputeResources{CPU: 2, Memory: 2},
		"db.t3.medium":  ComputeResources{CPU: 2, Memory: 4},
		"db.t3.large":   ComputeResources{CPU: 2, Memory: 8},
		"db.t3.xlarge":  ComputeResources{CPU: 4, Memory: 16},
		"db.t3.2xlarge": ComputeResources{CPU: 8, Memory: 32},
		"db.t4g.small":  ComputeResources{CPU: 2, Memory: 2},
		"db.t4g.medium": ComputeResources{CPU: 2, Memory: 4},
	}
	return result[key]

}
