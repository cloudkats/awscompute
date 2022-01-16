package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	. "awscompute/internal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fatih/color"
	flag "github.com/spf13/pflag"

	compute "awscompute/aws"
)

func main() {
	var saveToFile bool
	var resources CommaSeparatedListFlag

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.Usage = func() {
		printHelp(flags)
	}

	flags.BoolVar(&saveToFile, "save", false, "Enable persistence, save to file")
	flags.VarP(&resources, "resources", "r", "Comma-separated list of resources to compute CPU/Memory")

	_ = flags.Parse(os.Args[0:])
	args := flags.Args()

	if len(args) == 0 {
		printHelp(flags)
		return
	}

	err := MatchSupportedTypes(resources)
	if err != nil {
		CheckError("resources", err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	CheckError("Error: failed to call", err)

	ctx := context.TODO()
	cmp := compute.New(ctx, cfg)
	acc, err := cmp.AccountAnalyzer()
	CheckError("iam", err)

	// done: ec2, lambda, rds, redshift
	data := make([][]string, 0)
	totalCPU := 0
	totalMemory := 0

	for _, res := range resources {
		result, err := cmp.ComputeResourcesByType(res)
		CheckError(res, err)
		_, _ = fmt.Fprint(os.Stdout, color.YellowString("\tType: %s ::  ", strings.ToUpper(result.Type)),
			color.BlueString("Count: %d. ", result.Count),
			color.GreenString("CPU: %d. Memory: %d GiB\n", result.CPU, result.Memory))
		data = append(data, []string{acc.Account, acc.Aliases, res, strconv.Itoa(result.Count), strconv.Itoa(result.CPU), strconv.Itoa(result.Memory)})
		totalMemory += result.Memory
		totalCPU += result.CPU
	}

	if saveToFile {
		fw := FileWriter{}
		file := fmt.Sprintf("data/%s.csv", acc.Account)
		err = fw.WriteToFile(file, data)
		if err != nil {
			ExitErrorf("Write to file", err)
		}
	}
	_, _ = fmt.Fprint(os.Stdout,
		color.MagentaString("\t\tTOTAL::CPU: %d. Memory: %d GiB\n", totalCPU, totalMemory))
}

func printHelp(fs *flag.FlagSet) {
	_, _ = fmt.Fprintf(os.Stderr, "\n"+strings.TrimSpace(help)+"\n")
	fs.PrintDefaults()
}

const help = `
awsls - list AWS resources.

USAGE:
    $ compute [flags]

FLAGS:
`
