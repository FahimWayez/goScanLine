package goscanline

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
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

//until '\n' or end of line it will read one line and trims a single trailing newline
func (s *Scanner) ReadLine() (string, error){
	s.mu.Lock()
	defer s.mu.Unlock()

	line, err := s.r.ReadString('\n')
	if errors.Is(err, io.EOF){
		if len(line) == 0{
			return "", io.EOF
		}
		err = nil
	}
	return trimNewLine(line), err
}

func trimNewLine(s string) string{
	if strings.HasSuffix(s, "\r\n"){
		return strings.TrimSuffix(s, "\r\n")
	}

	if strings.HasSuffix(s, "\n"){
		return strings.TrimSuffix(s, "\n")
	}

	return s
}
