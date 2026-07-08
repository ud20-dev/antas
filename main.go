package main

import (
	"fmt"
	"os"

	"antas.com/antas/console"
	"antas.com/antas/pdf"
	"antas.com/antas/renderer"
)

func main() {	
	args := os.Args[1:] // Skip the program name
	if err := Run(args); err != nil {
		console.PrintWithStyle(
			console.ErrorStyle,
			"%s",
			err,
		)
		os.Exit(1)
	}
}

func Run(Args []string) error {
	if len(Args) != 1 {
		return fmt.Errorf("Usage: antas <path/to/file.pdf>")
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
		console.PrintWithStyle(console.InfoStyle, "Rendered page %d to %s", i+1, outputFile)
	}
	console.PrintWithStyle(console.SuccessStyle, "All pages rendered to %s", outputDir)
	return nil
}