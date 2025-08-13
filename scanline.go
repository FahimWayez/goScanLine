package goscanline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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


func parseInt(src string, bitSize int, out any) error{
	v, err := strconv.ParseInt(strings.TrimSpace(src), 10, bitSize)

	if err != nil{
		return fmt.Errorf("%w: int: %v", ErrParse, err)
	}
	switch p := out.(type){
	case *int64:
		*p = v

	case *int:
		*p = int(v)
			
	default:
		return fmt.Errorf("%w: internal int out type %T", ErrUnsupported, out)
	}
	
	return nil
}

func parseUint(src string, bitSize int, out any) error{
	v, err := strconv.ParseUint(strings.TrimSpace(src), 10, bitSize)

	if err != nil{
		return fmt.Errorf("%w: uint: %v", ErrParse, err)
	}

	switch p := out.(type){
	case *uint64:
		*p = v

	case *uint:
		*p = uint(v)	
	default:
		return fmt.Errorf("%w: internal uint out type %T", ErrUnsupported, out)
	}
	
	return nil
}
