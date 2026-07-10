package main

import (
	"fmt"
	"os"

	pflag "github.com/spf13/pflag"

	"github.com/ud20-dev/antas/console"
	"github.com/ud20-dev/antas/pdf"
	"github.com/ud20-dev/antas/renderer"
)


func main() {
	format := pflag.StringP("format", "f", "human", console.GetReportersUsage())
	pflag.Parse()
	reporter, err := console.GetReporter(*format)
	if err != nil {
		console.PrintWithStyle(
			console.ErrorStyle,
			"%v",
			err,
		)
		os.Exit(2)
	}
	args := pflag.Args()
	if err := Run(args, reporter); err != nil {
		reporter.Error(err)
		os.Exit(1)
	}
}

func Run(Args []string, reporter console.Reporter) error {
	if len(Args) != 1 {
		return fmt.Errorf("Usage: antas <path/to/file.pdf>, %v", Args)
	}
	inputFile := Args[0]
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: %s", inputFile)
	}
	
	if err := renderer.Init(); err != nil{
		return err
	}
	defer renderer.Close()
	outputDir, err := pdf.GetPDFOutputPath(Args[0])
	if err != nil {
		return fmt.Errorf("Error getting PDF output path: %v", err)
	}
	pageCount, err := renderer.GetPageCount(inputFile)
	if err != nil {
		return fmt.Errorf("Error getting page count: %v", err)
	}
	for i := range pageCount {
		outputFile := fmt.Sprintf("%s/page_%d.png", outputDir, i+1)
		err = renderer.RenderPage(inputFile, i, outputFile)
		if err != nil {
			return fmt.Errorf("Error rendering page %d: %v", i+1, err)
		}
		reporter.PageRendered(i+1, outputFile)
	}
	reporter.Done(outputDir, pageCount)
	return nil
}