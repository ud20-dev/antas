package main

import (
	"fmt"
	"os"

	pflag "github.com/spf13/pflag"

	"github.com/ud20-dev/antas/console"
	"github.com/ud20-dev/antas/cmd"
)


func main() {
	var ctx RunContext
	pflag.StringVarP(&ctx.Format,"format", "f", "human", 
		fmt.Sprintf("the format to use when printing/logging one of %s", console.GetReportersUsage()),
	)
	pflag.BoolVarP(&ctx.Help,"help", "h", false, "Show help information")
	pflag.BoolVarP(&ctx.Version,"version", "v", false, "Show version information")
	pflag.Parse()
	ctx.Args = pflag.Args()

	return_code, err := Dispatch(ctx)
	if return_code == ExitBadCLIUsage {
		console.PrintWithStyle(
			console.ErrorStyle,
			"%v",
			err,
		)
	}
	os.Exit(return_code)
}

func Dispatch(ctx RunContext) (int, error) {
	if ctx.Help{
		cmd.PrintHelp()
		return ExitSuccess, nil
	} else if ctx.Version {
		cmd.PrintVersion()
		return ExitSuccess, nil
	} 

	reporter, err := console.GetReporter(ctx.Format)
	if err != nil {
		return ExitBadCLIUsage, err
	}

	err = cmd.CanonicalRun(ctx.Args, reporter)
	if err != nil {
		reporter.Error(err)
		return ExitGenericFailure, err
	}
	return 0, nil
}