package renderer

import (
	"image/png"

	"os"
	
	"github.com/klippa-app/go-pdfium/requests"
)

// RenderPage renders a page of a PDF file to an image and saves it to the specified output path.
// the specified page is 0-indexed, so page 1 is index 0, page 2 is index 1, etc.
// The output image will be in PNG format and will truncate/create the output file if it already exists.
func RenderPage(filePath string, page int, output string) error {
	// Load the PDF file into a byte array.
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	// Render the page in DPI DPI.
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: DPI, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: doc.Document,
				Index:    page,
			},
		}, // The page to render, 0-indexed.
	})
	if err != nil {
		return err
	}

	// The Render* methods return a cleanup function that has to be called when
	// using webassembly to make sure resources are cleaned up. Do this after
	// you are done with the returned image object.
	defer pageRender.Cleanup()

	// Write the output to a file.
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, pageRender.Result.Image)
	if err != nil {
		return err
	}

	return nil
}