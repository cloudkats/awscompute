package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"awsconfig/internal"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"

	compute "awsconfig/aws"
)

func main() {
	var resources internal.CommaSeparatedListFlag

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		printHelp(flags)
	}
	flags.VarP(&resources, "resources", "r", "Comma-separated list of resources to compute CPU/Memory")

	_ = flags.Parse(os.Args[0:])
	args := flags.Args()

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, color.RedString("Error: resource types expected\n"))
		printHelp(flags)
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	checkError("Error: failed to call", err)

	ctx := context.TODO()

	cmp := compute.New(ctx, cfg)

	acc, err := cmp.AccountAnalyzer()
	checkError("iam", err)

	// done: ec2, lambda, rds, redshift
	data := make([][]string, 0)
	totalCPU := 0
	totalMemory := 0

	for _, res := range resources {
		result, err := cmp.ComputeResourcesByType(res)
		checkError("ec2", err)

		fmt.Fprint(os.Stdout, color.YellowString("\tType: %s ::  ", strings.ToUpper(result.Type)),
			color.BlueString("Count: %d. ", result.Count),
			color.GreenString("CPU: %d. Memory: %d GiB\n", result.CPU, result.Memory))
		data = append(data, []string{acc.Account, acc.Aliases, res, strconv.Itoa(result.CPU), strconv.Itoa(result.Memory)})
		totalMemory += result.Memory
		totalCPU += result.CPU
	}

	fw := internal.FileWriter{}
	file := fmt.Sprintf("data/%s/compute.csv", acc.Account)
	err = fw.WriteToFile(file, data)
	if err != nil {
		exitErrorf("Write to file", err)
	}
	fmt.Fprint(os.Stdout,
		color.MagentaString("\tTOTAL::CPU: %d. Memory: %d GiB\n", totalCPU, totalMemory))
}

func checkError(message string, err error) {
	if err != nil {
		exitErrorf(message, err)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprint(os.Stderr, color.RedString("%s: %s\n", msg, args))
	os.Exit(1)
}

func printHelp(fs *flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "\n"+strings.TrimSpace(help)+"\n")
	fs.PrintDefaults()
}

const help = `
awsls - list AWS resources.

USAGE:
    $ compute [flags]

FLAGS:
`
