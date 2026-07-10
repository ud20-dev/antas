package console

import (
	"encoding/json"
	"os"
)

func (JsonReporter) PageRendered(pageNum int, path string) {
}

func (JsonReporter) Done(outDir string, pageCount int) {
	result := Result{
		OK:        true,
		OutDir:    outDir,
		PageCount: pageCount,
	}
	json.NewEncoder(os.Stdout).Encode(result)
}

func (JsonReporter) Error(err error) {
	result := Result{
		OK:    false,
		Error: err.Error(),
	}
	json.NewEncoder(os.Stdout).Encode(result)
}