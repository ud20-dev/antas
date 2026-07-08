package renderer

import (
	"os"
	
	"github.com/klippa-app/go-pdfium/requests"
)

func GetPageCount(filePath string) (int, error) {
	// Load the PDF file into a byte array.
	pdfBytes, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Open the PDF using PDFium (and claim a worker)
	doc, err := instance.OpenDocument(&requests.OpenDocument{
		File: &pdfBytes,
	})
	if err != nil {
		return 0, err
	}

	// Always close the document, this will release its resources.
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: doc.Document,
	})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
		Document: doc.Document,
	})
	if err != nil {
		return 0, err
	}

	return pageCount.PageCount, nil
}