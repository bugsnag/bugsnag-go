package errors

import (
	"runtime"
	"testing"

	"github.com/pkg/errors"
)

func Test_StackFrames(t *testing.T) {
	err := boom(5)
	stack := getStack(err)
	frames := getFrames(stack)
	// for _, frame := range frames {
	// 	fmt.Printf("%s:%d\n\t%s\n", frame.File, frame.Line, frame.Function)
	// }

	err2 := &Error{Err: err, stack: stack}
	// fmt.Println(string(err2.Stack()))
	frames2 := err2.StackFrames()
	if name, name2 := frames[0].Function, frames2[0].Func().Name(); name != name2 {
		t.Errorf("top frames don't match\n%s\n%s", name, name2)
	}
}

func boom(depth int) error {
	if depth > 0 {
		return boom(depth - 1)
	}
	return errors.New("boom")
}

func getStack(err error) []uintptr {
	type withStackTrace interface {
		StackTrace() errors.StackTrace
	}
	frames := err.(withStackTrace).StackTrace()
	stack := make([]uintptr, len(frames))
	for i, f := range frames {
		stack[i] = uintptr(f)
	}
	return stack
}

func getFrames(stack []uintptr) (frames []runtime.Frame) {
	callers := runtime.CallersFrames(stack)
	for {
		frame, more := callers.Next()
		frames = append(frames, frame)
		if !more {
			break
		}
	}
	return frames
}
