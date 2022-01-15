package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func ec2Analyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	fmt.Println("EC2 Analyzer")

	svc := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}

	p := ec2.NewDescribeInstancesPaginator(svc, input)
	iMap := map[string]int{}
	iTypes, err := instances(ctx, p, iMap)
	if err != nil {
		return nil, err
	}
	typesInput := &ec2.DescribeInstanceTypesInput{
		InstanceTypes: iTypes,
	}

	o := ec2.NewDescribeInstanceTypesPaginator(svc, typesInput)
	iTypesMap := map[string]ComputeResources{}
	for o.HasMorePages() {
		resp, err := o.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.InstanceTypes {
			iType := fmt.Sprintf("%v", i.InstanceType)
			iTypesMap[iType] = ComputeResources{
				CPU:    int(*i.VCpuInfo.DefaultVCpus),
				Memory: int(*i.MemoryInfo.SizeInMiB),
			}
		}
	}

	// fmt.Println(iMap)
	// fmt.Println(iTypes)

	iCPU := 0
	iMemory := 0
	instances := 0

	for iType, count := range iMap {
		cpu := iTypesMap[iType].CPU
		memory := iTypesMap[iType].Memory
		iCPU += cpu * count
		iMemory += memory * count
		instances += count
	}
	return &ComputeOutput{
		CPU:    iCPU,
		Memory: iMemory / 1024,
		Type:   "ec2",
		Count:  instances,
	}, nil
}

func instances(ctx context.Context, p *ec2.DescribeInstancesPaginator, iMap map[string]int) ([]types.InstanceType, error) {
	result := make([]types.InstanceType, 0)
	for p.HasMorePages() {
		resp, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, reservations := range resp.Reservations {
			for _, i := range reservations.Instances {
				iType := fmt.Sprintf("%v", i.InstanceType)
				iMap[iType]++
				result = append(result, types.InstanceType(iType))
			}
		}
	}
	return result, nil
}
