package cmd

import (
	"fmt"
	"os"

	
	"github.com/ud20-dev/antas/pdf"
	"github.com/ud20-dev/antas/renderer"
	"github.com/ud20-dev/antas/console"
)


func CanonicalRun(inputFile string, reporter console.Reporter) error {
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", inputFile)
	}

	if err := renderer.Init(); err != nil{
		return err
	}
	defer func() { _ = renderer.Close() }()
	outputDir, err := pdf.GetPDFOutputPath(inputFile)
	if err != nil {
		return fmt.Errorf("error getting PDF output path: %v", err)
	}
	pageCount, err := renderer.GetPageCount(inputFile)
	if err != nil {
		return fmt.Errorf("error getting page count: %v", err)
	}
	for i := range pageCount {
		outputFile := fmt.Sprintf("%s/page_%d.png", outputDir, i+1)
		err = renderer.RenderPage(inputFile, i, outputFile)
		if err != nil {
			return fmt.Errorf("error rendering page %d: %v", i+1, err)
		}
		reporter.PageRendered(i+1, outputFile)
	}
	reporter.Done(outputDir, pageCount)
	return nil
}