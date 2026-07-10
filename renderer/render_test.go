package renderer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ud20-dev/antas/renderer"
)

var samplePDFs = []string{
	"../tests_samples/copia-dsl.pdf",
	"../tests_samples/example_domain_blank.pdf",
	"../tests_samples/the_long_faces:Jane!.pdf",
}

func TestMain(m *testing.M) {
	if err := renderer.Init(); err != nil {
		panic("renderer.Init: " + err.Error())
	}
	code := m.Run()
	_ = renderer.Close()
	os.Exit(code)
}

func TestGetPageCount(t *testing.T) {
	for _, pdf := range samplePDFs {
		t.Run(filepath.Base(pdf), func(t *testing.T) {
			count, err := renderer.GetPageCount(pdf)
			if err != nil {
				t.Fatalf("GetPageCount: %v", err)
			}
			if count < 1 {
				t.Errorf("expected at least 1 page, got %d", count)
			}
		})
	}
}

func TestRenderPage(t *testing.T) {
	for _, pdf := range samplePDFs {
		t.Run(filepath.Base(pdf), func(t *testing.T) {
			outFile := filepath.Join(t.TempDir(), "page_1.png")
			if err := renderer.RenderPage(pdf, 0, outFile); err != nil {
				t.Fatalf("RenderPage: %v", err)
			}
			info, err := os.Stat(outFile)
			if err != nil {
				t.Fatalf("output file not created: %v", err)
			}
			if info.Size() == 0 {
				t.Error("output file is empty")
			}
		})
	}
}
