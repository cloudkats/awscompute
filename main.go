package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fatih/color"

	compute "awsconfig/aws"
)

func main() {

	// nameFilter := os.Args[1]
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		exitErrorf("Error: failed to call", err)
	}

	ctx := context.TODO()

	cmp := compute.New(ctx, cfg)
	result, err := cmp.ComputeResourcesByType("ec2")
	if err != nil {
		exitErrorf("ec2", err.Error())
	}
	fmt.Fprint(os.Stdout, color.YellowString("\tType: %s ::  ", result.Type),
		color.GreenString("CPU: %d. Memory: %d GiB\n", result.CPU, result.Memory))
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprint(os.Stderr, color.RedString("%s: %s\n", msg, args))
	os.Exit(1)
}
