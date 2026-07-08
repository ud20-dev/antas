package pdf

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)


/* 
this function takes a path to a pdf file create the folder structure and returns where the output image should be saved.
since we render one image by page.
the output path will follow this convention:

{TEMP_DIR}/{UUID}/{hash_of_pdf_file}/

the program calling is expected to extend the path with the page number like so:

{TEMP_DIR}/{UUID}/{hash_of_pdf_file}/page_{page_number}.png (,jpeg or any other format supported by the image package)
	

For example, assuming we're in a linux environment

#IN : "./test_samples/idiot.pdf"

#OUT: "/tmp/3c82a360-6678-48d7-923a-8cb25a895742/hash_of_pdf_file/"
*/ 
func GetPDFOutputPath(pdfPath string) (string, error) {
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	outputPath := filepath.Join(os.TempDir(), uuid.NewString(), hex.EncodeToString(hash[:]))
	os.MkdirAll(outputPath, os.ModePerm)
	return outputPath, nil
}