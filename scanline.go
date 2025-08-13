package goscanline

import (
	"bufio"
	"encoding"
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


func assign(src string, dest any) error{
	switch d := dest.(type){
	case *string:
		*d = src
		return nil

	case *bool:
		v, err := strconv.ParseBool(strings.TrimSpace(src))
		if err != nil{
			return fmt.Errorf("%w: bool: %v", ErrParse, err)
		}
		*d = v
	
		return nil

	case *int:
		return parseInt(src, 0, d)

	case *int8:
		var v int64
		if err := parseInt(src, 8, &v); err != nil{return err}
		*d = int8(v); return nil

	case *int16:
		var v  int64
		if err := parseInt(src, 16, &v); err != nil{return err}
		*d = int16(v); return nil

	case *int32:
		var v int64
		if err := parseInt(src, 32, &v); err != nil{return err}
		*d = int32(v); return nil

	case *int64:
		return parseInt(src, 64, d)

	case *uint:
		return parseUint(src, 0, d)
	
	case *uint8:
		var v uint64
		if err := parseUint(src, 8, &v); err != nil { return err }
		*d = uint8(v); return nil

	case *uint16:
		var v uint64
		if err := parseUint(src, 16, &v); err != nil { return err }
		*d = uint16(v); return nil

	case *uint32:
		var v uint64
		if err := parseUint(src, 32, &v); err != nil { return err }
		*d = uint32(v); return nil

	case *uint64:
		return parseUint(src, 64, d)

	case *float32:
		f, err := strconv.ParseFloat(strings.TrimSpace(src), 32)
		if err != nil{
			return fmt.Errorf("%w: float32: %v", ErrParse, err)
		}
		*d = float32(f); return nil
		
	case *float64:
		f, err := strconv.ParseFloat(strings.TrimSpace(src), 64)
		if err != nil {
			return fmt.Errorf("%w: float64: %v", ErrParse, err)
		}
		*d = f; return nil
	}

	if u, ok := dest.(encoding.TextUnmarshaler); ok{
		if err := u.UnmarshalText([]byte(src)); err != nil{
			return fmt.Errorf("%w: text: %v", ErrParse, err)
		}

		return nil
	}

	return fmt.Errorf("%w: %T", ErrUnsupported, dest)
}

//reads a line and assigns it into dest (basically pointer to string/bool/int*/uint*/float* or TextUnmarshaler)
func (s *Scanner) Scan(dest any) error{
	line, err := s.ReadLine()

	if err != nil{
		return err
	}

	return assign(line, dest)
}

//writes prompt without newline to the prompt writer, then behaves similar to scan
func (s *Scanner) ScanPrompt(prompt string, dest any) error{
	fmt.Fprint(s.w, prompt)
	return s.Scan(dest)
}