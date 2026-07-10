package pdf

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var pathSeq atomic.Int64

/*
this function takes a path to a pdf file create the folder structure and returns where the output image should be saved.
since we render one image by page.
the output path will follow this convention:

	{TEMP_DIR}/{hash_of_pdf_file}/{when}
	${when} = {now.UnixTime}-{CurrentPID}-{seq}

where {seq} is a process-wide atomic counter that guarantees uniqueness across concurrent
goroutines calling this function at the same unix second.

the program calling is expected to extend the path with the page number like so:

	{TEMP_DIR}/{hash_of_pdf_file}/{when}/page_{page_number}.png (,jpeg or any other format supported by the image package)


For example, assuming we're in a linux environment

#IN : "./test_samples/idiot.pdf"

#OUT: "/tmp/096c8484d1d42f3ca41c4952b5eabd69bda0fd04508c03b87d58608d84089911/1783520643-24423-1"
*/
func GetPDFOutputPath(pdfPath string) (string, error) {
	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", err
	}

	when := fmt.Sprintf("%d-%d-%d", time.Now().Unix(), os.Getpid(), pathSeq.Add(1))

	hash := sha256.Sum256(data)
	outputPath := filepath.Join(
		os.TempDir(), 
		hex.EncodeToString(hash[:]),
		when,
	)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return "", err
	}
	return outputPath, nil
}