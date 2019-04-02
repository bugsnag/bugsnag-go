package errors

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

// A StackFrame contains all necessary information about to generate a line
// in a callstack.
type StackFrame struct {
	File           string
	LineNumber     int
	Name           string
	Package        string
	ProgramCounter uintptr
}

// NewStackFrame popoulates a stack frame object from the program counter.
func NewStackFrame(pc uintptr) (frame StackFrame) {
	frame = StackFrame{ProgramCounter: pc}
	if frame.Func() == nil {
		return
	}
	frame.Package, frame.Name = packageAndName(frame.Func().Name())

	// pc -1 because the program counters we use are usually return addresses,
	// and we want to show the line that corresponds to the function call
	frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
	return
}

// NewStackFrameFromRuntime populates a stack frame object from a runtime.Frame object.
func NewStackFrameFromRuntime(frame runtime.Frame) StackFrame {
	var pkg, name string
	if frame.Func != nil {
		pkg, name = packageAndName(frame.Func.Name())
	} else if frame.Function != "" {
		pkg, name = packageAndName(frame.Function)
	} else {
		return StackFrame{}
	}
	return StackFrame{
		File:           frame.File,
		LineNumber:     frame.Line,
		Name:           name,
		Package:        pkg,
		ProgramCounter: frame.PC,
	}
}

// Func returns the function that this stackframe corresponds to
func (frame *StackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}
	return runtime.FuncForPC(frame.ProgramCounter)
}

// String returns the stackframe formatted in the same way as go does
// in runtime/debug.Stack()
func (frame *StackFrame) String() string {
	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.ProgramCounter)

	source, err := frame.SourceLine()
	if err != nil {
		return str
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Name, source)
}

// SourceLine gets the line of code (from File and Line) of the original source if possible
func (frame *StackFrame) SourceLine() (string, error) {
	data, err := ioutil.ReadFile(frame.File)

	if err != nil {
		return "", err
	}

	lines := bytes.Split(data, []byte{'\n'})
	if frame.LineNumber <= 0 || frame.LineNumber >= len(lines) {
		return "???", nil
	}
	// -1 because line-numbers are 1 based, but our array is 0 based
	return string(bytes.Trim(lines[frame.LineNumber-1], " \t")), nil
}

// packageAndName splits a package path-qualified function name into the package path and function name.
func packageAndName(qualifiedName string) (pkg string, name string) {
	name = qualifiedName
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return
}

func pcsToFrames(pcs []uintptr) []runtime.Frame {
	frames := runtime.CallersFrames(pcs)
	s := make([]runtime.Frame, 0, len(pcs))
	for {
		frame, more := frames.Next()
		s = append(s, frame)
		if !more {
			break
		}
	}
	return s
}

func runtimeToErrorFrames(rtFrames []runtime.Frame) []StackFrame {
	frames := make([]StackFrame, len(rtFrames))
	for i, f := range rtFrames {
		frames[i] = NewStackFrameFromRuntime(f)
	}
	return frames
}
