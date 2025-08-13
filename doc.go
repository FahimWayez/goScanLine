// Package goscanline provides minimal helpers for line-based input in Go CLIs.
//
// It reads a full line from an input stream (including spaces), trims the trailing
// newline, and assigns the result into a typed destination (e.g., *string, *bool,
// *int*, *uint*, *float32/64, or any type implementing encoding.TextUnmarshaler).
//
// The package includes helpers for interactive command-line programs:
//
//   - ScanPrompt prints a prompt (to stderr) and then reads a line.
//   - ScanCtx reads with a context; the read cancels if the context is done.
//   - ScanSecret reads a line without echo when stdin is a terminal (uses x/term).
//
// Package-level functions operate on a default Scanner wired to stdin/stderr.
// For custom I/O (tests, pipes), construct a Scanner with New(r, w).
package goscanline
