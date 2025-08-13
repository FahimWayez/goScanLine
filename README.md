## `README.md`

# goscanline

Minimal, dependency-free helpers for **line-based input** in Go CLIs.

- Reads a **full line** from stdin (keeps spaces), trims the newline.
- Assigns into typed destinations: `*string`, `*bool`, `*int*`, `*uint*`, `*float32/64`, or any `encoding.TextUnmarshaler`.
- Extras for real CLI UX: **cancellable reads** (`ScanCtx`) and **no-echo secrets** (`ScanSecret`, via `x/term`).
- Prompts print to **stderr** by default, so **stdout** stays clean for piping.

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/FahimWayez/goScanLine.svg)](https://pkg.go.dev/github.com/FahimWayez/goScanLine)
[![Go Report Card](https://goreportcard.com/badge/github.com/FahimWayez/goScanLine)](https://goreportcard.com/report/github.com/FahimWayez/goScanLine)

---

## Install

```bash
go get github.com/FahimWayez/goScanLine
```

Go 1.21+ recommended.

---

## Quick start

```go
package main

import (
	"fmt"

	"github.com/FahimWayez/goScanLine"
)

func main() {
	var name string
	if err := goscanline.ScanPrompt("Enter your full name: ", &name); err != nil {
		panic(err)
	}
	fmt.Println("Hello,", name)
}
```

---

## API overview

### Package-level helpers

```go
func ReadLine() (string, error)
func Scan(dest any) error
func ScanPrompt(prompt string, dest any) error
func ScanCtx(ctx context.Context, dest any) error
func ScanSecret(prompt string, dest *string) error
func ScanT[T any]() (T, error)
```

* `dest` must be a **pointer** to a supported type.
* `ScanPrompt` writes to **stderr** (keeps stdout for pipelines).
* `ScanCtx` cancels a blocking read (e.g., on timeout).
* `ScanSecret` reads without echo when stdin is a terminal; otherwise falls back to normal `Scan`.

### Scanner (custom streams)

```go
type Scanner struct{ /* ... */ }

func New(r io.Reader, w io.Writer) *Scanner
func (s *Scanner) ReadLine() (string, error)
func (s *Scanner) Scan(dest any) error
func (s *Scanner) ScanPrompt(prompt string, dest any) error
func (s *Scanner) ScanCtx(ctx context.Context, dest any) error
func (s *Scanner) ScanSecret(prompt string, dest *string) error
```

Use `Scanner` when you want to inject your own reader/writer (tests, pipes, in-memory buffers).

---

## Examples

**Parse a full-line string**

```go
var name string
_ = goscanline.ScanPrompt("Name: ", &name) // "Ada Lovelace"
```

**Parse a number from the line**

```go
var age int
if err := goscanline.ScanPrompt("Age: ", &age); err != nil {
    // handle invalid integer
}
```

**Cancellable input**

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
var v string
if err := goscanline.ScanCtx(ctx, &v); err != nil {
    // context deadline exceeded (e.g., user never hit Enter)
}
```

**Secret input (no echo)**

```go
var pwd string
if err := goscanline.ScanSecret("Password: ", &pwd); err != nil { /* ... */ }
```

**Custom reader/writer (tests)**

```go
in := strings.NewReader("hello\n")
var out bytes.Buffer
s := goscanline.New(in, &out)
var v string
_ = s.ScanPrompt("Enter: ", &v) // prompt written to out
```

---

## Notes

* The global `Default` scanner serializes calls. For concurrent interactive prompts, create separate `Scanner` instances.
* Errors: parse failures wrap `ErrParse`; unsupported destination types return `ErrUnsupported`.
* Prompts are sent to **stderr**; print results to **stdout** to keep pipelines predictable.

---

## Versioning

Semantic versioning. Tag releases so users can pin versions.

```bash
git tag v0.2.0
git push --tags
```

---

## Contributing

Issues and PRs welcome. Please include tests, keep dependencies minimal, and maintain backward compatibility.

Run checks locally:

```bash
gofmt -s -l .
go vet ./...
go test ./...    # add -race for data races
```

---

## License

MIT © Fahim Wayez

````
