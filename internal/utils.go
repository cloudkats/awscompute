package internal

import (
	"errors"
	"fmt"
)

var supported = map[string]bool{
	"ec2":      true,
	"lambda":   true,
	"rds":      true,
	"redshift": true,
}

func MatchSupportedTypes(resources []string) error {
	for _, resource := range resources {
		if _, found := supported[resource]; !found {
			msg := fmt.Sprintf("Resource %v not supported", resource)
			return errors.New(msg)
		}
	}
	return nil
}
