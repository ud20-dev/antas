package pdf_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ud20-dev/antas/pdf"
)

const samplePDF = "../tests_samples/example_domain_blank.pdf"

func TestGetPDFOutputPath_GoroutineUniqueness(t *testing.T) {
	const n = 50
	paths := make([]string, n)
	errs := make([]error, n)

	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			path, err := pdf.GetPDFOutputPath(samplePDF)
			paths[idx] = path
			errs[idx] = err
		}(i)
	}
	wg.Wait()

	seen := make(map[string]bool, n)
	for i, path := range paths {
		if errs[i] != nil {
			t.Fatalf("goroutine %d: %v", i, errs[i])
		}
		if seen[path] {
			t.Errorf("goroutine %d produced duplicate path: %s", i, path)
		}
		seen[path] = true
	}
}

// TestGetPDFOutputPath_Grouping verifies that all paths for the same file share
// the same parent directory (the content-hash directory under TempDir).
func TestGetPDFOutputPath_Grouping(t *testing.T) {
	const n = 20
	paths := make([]string, n)
	errs := make([]error, n)

	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			path, err := pdf.GetPDFOutputPath(samplePDF)
			paths[idx] = path
			errs[idx] = err
		}(i)
	}
	wg.Wait()

	var hashDir string
	for i, path := range paths {
		if errs[i] != nil {
			t.Fatalf("goroutine %d: %v", i, errs[i])
		}
		parent := filepath.Dir(path)
		if hashDir == "" {
			hashDir = parent
		} else if parent != hashDir {
			t.Errorf("goroutine %d: expected parent %s, got %s", i, hashDir, parent)
		}
	}

	// The hash directory must live directly under the system temp dir.
	if filepath.Dir(hashDir) != filepath.Clean(os.TempDir()) {
		t.Errorf("hash dir %s is not directly under temp dir %s", hashDir, os.TempDir())
	}
}
