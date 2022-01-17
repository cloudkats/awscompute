package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
)

func elasticacheAnalyzer(ctx context.Context, cfg aws.Config) (*ComputeOutput, error) {
	svc := elasticache.NewFromConfig(cfg)
	p := elasticache.NewDescribeCacheClustersPaginator(svc, &elasticache.DescribeCacheClustersInput{})
	iCPU, iMemory, instances := 0, 0, 0
	iMap := map[string]int{}

	for p.HasMorePages() {
		resp, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, i := range resp.CacheClusters {
			iType := fmt.Sprintf("%v", *i.CacheNodeType)
			var numNodes int
			numNodes = int(*i.NumCacheNodes)
			instances += numNodes
			iMap[iType]++
			info, err := elasticacheOfferings(iType)
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
		Type:      "redis",
		Count:     instances,
		Resources: iMap,
	}, nil
}

func elasticacheOfferings(key string) (*ComputeResources, error) {
	// TODO: do use more real number, as cache memory is float e.g. 1.34 and etc
	result := map[string]ComputeResources{
		"cache.t3.small":  {CPU: 2, Memory: 2},
		"cache.t3.medium": {CPU: 2, Memory: 3},
	}
	value, isPresent := result[key]
	if isPresent {
		return &value, nil
	}
	return nil, fmt.Errorf("type %v not supported. please add from https://aws.amazon.com/elasticache/pricing/", key)
}
