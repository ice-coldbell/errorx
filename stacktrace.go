package errorx

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

type frame struct {
	pc       uintptr
	file     string
	function string
	line     int
}

func frameForPC(pc uintptr) frame {
	f := frame{pc: pc}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return f
	}

	f.function = fn.Name()
	f.file, f.line = fn.FileLine(pc)
	return f
}

func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.function)
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file)
		default:
			io.WriteString(s, path.Base(f.file))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.line))
	case 'n':
		io.WriteString(s, f.function)
	case 'v':
		if s.Flag('#') {
			fmt.Fprintf(s, "errorx.frame{pc: %#x, file: %s, function:%s, line:%d}", f.pc, f.file, f.function, f.line)
			return
		}
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

func (f frame) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s %s:%d", f.function, f.file, f.line)), nil
}

type stack []frame

func (s stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, f := range s {
				io.WriteString(st, "\n")
				f.Format(st, verb)
			}
		case st.Flag('#'):
			io.WriteString(st, "[]errorx.frame{\n")
			for _, f := range s {
				fmt.Fprintf(st, "\t%#v,\n", f)
			}
			io.WriteString(st, "}")
		default:
			s.formatSlice(st, verb)
		}
	case 's':
		s.formatSlice(st, verb)
	}
}

func (s stack) formatSlice(st fmt.State, verb rune) {
	io.WriteString(st, "[")
	for i, f := range s {
		if i > 0 {
			io.WriteString(st, " ")
		}
		f.Format(st, verb)
	}
	io.WriteString(st, "]")
}

func (s stack) StackTrace() []uintptr {
	pcs := make([]uintptr, len(s))
	for i := 0; i < len(s); i++ {
		pcs[i] = s[i].pc
	}
	return pcs
}

// for test injection
var runtimeCallers = runtime.Callers

func callers(skip int) stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtimeCallers(2+skip, pcs[:])
	var st stack
	for i := range pcs[0:n] {
		st = append(st, frameForPC(pcs[i]))
	}
	return st
}
