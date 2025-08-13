### `README.md`

# goscanline

Read a **full line** from stdin (including spaces) without ceremony.  
Typed assignment, cancellable reads, secret input (no echo). Zero deps beyond `x/term`.

[![Go Reference](https://pkg.go.dev/badge/github.com/FahimWayez/goscanline.svg)](https://pkg.go.dev/github.com/FahimWayez/goscanline)

## Install

```bash
go get github.com/FahimWayez/goscanline
```

Go 1.21+ recommended.

## Quick start

```go
package main

import (
	"fmt"

	"github.com/FahimWayez/goscanline"
)

func main() {
	var name string
	if err := goscanline.ScanPrompt("Enter your full name: ", &name); err != nil {
		panic(err)
	}
	fmt.Println("Hello,", name)
}
```

## Features

* Full-line read with precise `\n`/`\r\n` trimming
* Typed assignment into:

  * `*string`, `*bool`
  * `*int`, `*int8/16/32/64`
  * `*uint`, `*uint8/16/32/64`
  * `*float32`, `*float64`
* **Cancellable reads**: `ScanCtx(ctx, &v)`
* **Secret input** (no echo): `ScanSecret("Password: ", &pwd)` (uses `x/term`)
* **Prompts to stderr** by default (keeps stdout clean for piping)
* Supports custom types via `encoding.TextUnmarshaler`
* Package-level helpers and configurable `Scanner`

## API (high level)

```go
// Package-level
func Scan(dest any) error
func ScanPrompt(prompt string, dest any) error
func ScanCtx(ctx context.Context, dest any) error
func ScanSecret(prompt string, dest *string) error
func ScanT[T any]() (T, error)

// Scanner for custom streams (tests/pipes)
type Scanner struct { /* ... */ }
func New(r io.Reader, w io.Writer) *Scanner
func (s *Scanner) ReadLine() (string, error)
func (s *Scanner) Scan(dest any) error
func (s *Scanner) ScanPrompt(prompt string, dest any) error
func (s *Scanner) ScanCtx(ctx context.Context, dest any) error
func (s *Scanner) ScanSecret(prompt string, dest *string) error
```

## Notes

* Global `Default` scanner serializes calls; for concurrent interactive prompts, create separate `Scanner`s.
* Prompts go to **stderr**; print your program’s results to **stdout** for pipelines.
* Errors: parse failures wrap `ErrParse`; unsupported destinations return `ErrUnsupported`.

## License

MIT © Fahim Wayez

````

---

### `LICENSE`
```text
MIT License

Copyright (c) 2025 Fahim Wayez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

---
