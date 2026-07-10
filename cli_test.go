package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/ud20-dev/antas/console"
)

// captureStdout redirects os.Stdout while fn runs and returns what was written.
// Not safe for parallel tests — do not call t.Parallel() in tests that use it.
func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// --- Exit code correctness ---

func TestDispatch_HelpFlag_ExitSuccess(t *testing.T) {
	code, _ := Dispatch(RunContext{Help: true, Format: "human"})
	if code != ExitSuccess {
		t.Errorf("expected exit %d, got %d", ExitSuccess, code)
	}
}

func TestDispatch_VersionFlag_ExitSuccess(t *testing.T) {
	code, _ := Dispatch(RunContext{Version: true, Format: "human"})
	if code != ExitSuccess {
		t.Errorf("expected exit %d, got %d", ExitSuccess, code)
	}
}

func TestDispatch_UnknownFormat_ExitBadCLIUsage(t *testing.T) {
	code, _ := Dispatch(RunContext{Format: "xml"})
	if code != ExitBadCLIUsage {
		t.Errorf("expected exit %d, got %d", ExitBadCLIUsage, code)
	}
}

func TestDispatch_NoArgs_ExitGenericFailure(t *testing.T) {
	code, _ := Dispatch(RunContext{Format: "human"})
	if code != ExitGenericFailure {
		t.Errorf("expected exit %d, got %d", ExitGenericFailure, code)
	}
}

func TestDispatch_FileNotFound_ExitGenericFailure(t *testing.T) {
	code, _ := Dispatch(RunContext{Format: "human", Args: []string{"/nonexistent/file.pdf"}})
	if code != ExitGenericFailure {
		t.Errorf("expected exit %d, got %d", ExitGenericFailure, code)
	}
}

// --- JSON format: stdout must always be valid JSON for exit 0 and 1 ---

func TestDispatch_JSONFormat_NoArgs_JSONStdout(t *testing.T) {
	var code int
	out := captureStdout(func() {
		code, _ = Dispatch(RunContext{Format: "json"})
	})
	if code != ExitGenericFailure {
		t.Errorf("expected exit %d, got %d", ExitGenericFailure, code)
	}
	assertJSONError(t, out)
}

func TestDispatch_JSONFormat_FileNotFound_JSONStdout(t *testing.T) {
	var code int
	out := captureStdout(func() {
		code, _ = Dispatch(RunContext{Format: "json", Args: []string{"/nonexistent/file.pdf"}})
	})
	if code != ExitGenericFailure {
		t.Errorf("expected exit %d, got %d", ExitGenericFailure, code)
	}
	assertJSONError(t, out)
}

func TestDispatch_JSONFormat_ValidPDF_JSONStdout(t *testing.T) {
	var code int
	out := captureStdout(func() {
		code, _ = Dispatch(RunContext{Format: "json", Args: []string{"tests_samples/example_domain_blank.pdf"}})
	})
	if code != ExitSuccess {
		t.Errorf("expected exit %d, got %d", ExitSuccess, code)
	}
	var result console.Result
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\noutput: %q", err, out)
	}
	if !result.OK {
		t.Errorf("expected ok=true, got error: %s", result.Error)
	}
	if result.PageCount < 1 {
		t.Errorf("expected page_count >= 1, got %d", result.PageCount)
	}
}

// TestDispatch_Exit2_NeverJSON verifies that exit code 2 (bad CLI usage) never
// causes Dispatch to emit JSON on stdout. The human-readable error is printed
// by main() after Dispatch returns, not by Dispatch itself.
func TestDispatch_Exit2_NeverJSON(t *testing.T) {
	cases := []RunContext{
		{Format: "xml"},
		{Format: "xml", Args: []string{"any.pdf"}},
		{Format: ""},
	}
	for _, ctx := range cases {
		var code int
		out := captureStdout(func() {
			code, _ = Dispatch(ctx)
		})
		if code != ExitBadCLIUsage {
			t.Errorf("format=%q: expected exit %d, got %d", ctx.Format, ExitBadCLIUsage, code)
		}
		trimmed := strings.TrimSpace(out)
		if trimmed != "" && json.Valid([]byte(trimmed)) {
			t.Errorf("format=%q: exit 2 must not produce JSON stdout, got: %q", ctx.Format, out)
		}
	}
}

func assertJSONError(t *testing.T, out string) {
	t.Helper()
	var result console.Result
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &result); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\noutput: %q", err, out)
	}
	if result.OK {
		t.Error("expected ok=false")
	}
	if result.Error == "" {
		t.Error("expected non-empty error field")
	}
}
