package errorwrap

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
	"strings"
)

type stack []uintptr

func (st *stack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, pc := range *st {
				f := frame(pc)
				fmt.Fprintf(s, "\n%+v", f)
			}
		}
	}
}

func callStack() *stack {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	var st stack = pcs[0:n]
	return &st
}

type frame uintptr

func (f frame) pc() uintptr {
	return uintptr(f) - 1
}

func (f frame) fileName() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

func (f frame) lineNumber() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

func (f frame) functionName() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.functionName())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.fileName())
		default:
			io.WriteString(s, path.Base(f.fileName()))
		}
	case 'd':
		io.WriteString(s, strconv.Itoa(f.lineNumber()))
	case 'n':
		io.WriteString(s, funcname(f.functionName()))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
