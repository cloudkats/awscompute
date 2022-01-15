package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
	// done: ec2, lambda, rds
	result, err := cmp.ComputeResourcesByType("rds")
	if err != nil {
		exitErrorf("ec2", err.Error())
	}
	fmt.Fprint(os.Stdout, color.YellowString("\tType: %s ::  ", strings.ToUpper(result.Type)),
		color.BlueString("Count: %d. ", result.Count),
		color.GreenString("CPU: %d. Memory: %d GiB\n", result.CPU, result.Memory))

	// compute everythyng at the very end
	// write to a file
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprint(os.Stderr, color.RedString("%s: %s\n", msg, args))
	os.Exit(1)
}
