package goscanline

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"time"
)

func TestReadLine_TrimLF(t *testing.T) {
	s := New(strings.NewReader("hello world\n"), nil)
	got, err := s.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello world" {
		t.Fatalf("got %q", got)
	}
}

func TestReadLine_TrimCRLF(t *testing.T) {
	s := New(strings.NewReader("hello\r\n"), nil)
	got, err := s.ReadLine()
	if err != nil {
		t.Fatal(err)
	}
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestScan_String(t *testing.T) {
	s := New(strings.NewReader("Fahim Wayez\n"), nil)
	var v string
	if err := s.Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != "Fahim Wayez" {
		t.Fatalf("got %q", v)
	}
}

func TestScan_Int(t *testing.T) {
	s := New(strings.NewReader("  42 \n"), nil)
	var v int
	if err := s.Scan(&v); err != nil {
		t.Fatal(err)
	}
	if v != 42 {
		t.Fatalf("got %d", v)
	}
}

func TestScanPrompt_WritesToPromptWriter(t *testing.T) {
	in := strings.NewReader("hello\n")
	var out bytes.Buffer
	s := New(in, &out)
	var v string
	if err := s.ScanPrompt("Enter: ", &v); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out.String(), "Enter: ") {
		t.Fatalf("prompt missing, got %q", out.String())
	}
	if v != "hello" {
		t.Fatalf("got %q", v)
	}
}

// when no newline or no EOF, ScanCtx should return ctx error
func TestScanCtx_Timeout(t *testing.T) {
	pr, _ := io.Pipe()
	s := New(pr, nil)

	var v string
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := s.ScanCtx(ctx, &v)
	if err == nil {
		t.Fatal("expected context deadline exceeded")
	}
	if !errorsIsDeadline(err) {
		t.Fatalf("expected deadline error, got: %v", err)
	}
}

func errorsIsDeadline(err error) bool {
	return err == context.DeadlineExceeded || strings.Contains(err.Error(), "deadline exceeded")
}
