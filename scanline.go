package goscanline

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
)

//errors returned by parsing/usage
var (
	ErrUnsupported = errors.New("goscanline: unsupported destination type")
	ErrParse = errors.New("goscanline: parse error")
)

type Scanner struct {
	mu sync.Mutex
	r *bufio.Reader
	src io.Reader //the original reader which is used for terminal detection
	w io.Writer
}

func New(r io.Reader, w io.Writer) *Scanner{
	if w == nil{
		w = os.Stderr
	}

	return &Scanner{
		r: bufio.NewReader(r),
		src: r,
		w: w, 
	}
}

var Default = New(os.Stdin, os.Stderr)

func SetDefaultReader(r io.Reader){
	Default.mu.Lock()
	defer Default.mu.Unlock()
	Default.r = bufio.NewReader(r)
	Default.src = r
}

func SetDefaultWriter(w io.Writer){
	Default.mu.Lock()
	defer Default.mu.Unlock()
	if w == nil{
		w = os.Stderr
	}
	Default.w = w
}