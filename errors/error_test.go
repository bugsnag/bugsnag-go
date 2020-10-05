package errors

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// fixture functions
func a() error {
	b(5)
	return nil
}

func b(i int) {
	c()
}

func c() {
	panic('a')
}

func TestParseStack(t *testing.T) {
	defer func() {
		err := New(recover(), 0)
		if err.Err.Error() != "97" {
			t.Errorf("Received incorrect error, expected 'a' got '%s'", err.Err.Error())
		}
		if err.TypeName() != "*errors.errorString" {
			t.Errorf("Error type was '%s'", err.TypeName())
		}
		expected := []StackFrame{
			StackFrame{Name: "TestParseStack.func1", File: "errors/error_test.go"},
			StackFrame{Name: "gopanic"},
			StackFrame{Name: "c", File: "errors/error_test.go", LineNumber: 22},
			StackFrame{Name: "c", File: "errors/error_test.go", LineNumber: 22},
			StackFrame{Name: "b", File: "errors/error_test.go", LineNumber: 18},
			StackFrame{Name: "a", File: "errors/error_test.go", LineNumber: 13},
		}
		assertStacksMatch(t, expected, err.StackFrames())
	}()

	a()
}

func TestSkipWorks(t *testing.T) {
	defer func() {
		err := New(recover(), 2)
		if err.Err.Error() != "97" {
			t.Errorf("Received incorrect error, expected 'a' got '%s'", err.Err.Error())
		}

		expected := []StackFrame{
			StackFrame{Name: "c", File: "errors/error_test.go", LineNumber: 22},
			StackFrame{Name: "c", File: "errors/error_test.go", LineNumber: 22},
			StackFrame{Name: "b", File: "errors/error_test.go", LineNumber: 18},
			StackFrame{Name: "a", File: "errors/error_test.go", LineNumber: 13},
		}

		assertStacksMatch(t, expected, err.StackFrames())
	}()

	a()
}

func assertStacksMatch(t *testing.T, expected []StackFrame, actual []StackFrame) {
	for index, frame := range expected {
		actualFrame := actual[index]
		if actualFrame.Name != frame.Name {
			t.Errorf("Frame %d method - Expected '%s' got '%s'", index, frame.Name, actualFrame.Name)
		}
		// Not using exact match as it would change depending on whether
		// the package is being tested within or outside of the $GOPATH
		if frame.File != "" && !strings.HasSuffix(actualFrame.File, frame.File) {
			t.Errorf("Frame %d file - Expected '%s' to end with '%s'", index, actualFrame.File, frame.File)
		}
		if frame.Package != "" && actualFrame.Package != frame.Package {
			t.Errorf("Frame %d package - Expected '%s' to end with '%s'", index, actualFrame.Package, frame.Package)
		}
		if frame.LineNumber != 0 && actualFrame.LineNumber != frame.LineNumber {
			t.Errorf("Frame %d line - Expected '%d' got '%d'", index, frame.LineNumber, actualFrame.LineNumber)
		}
	}
}

type testErrorWithStackFrames struct {
	Err *Error
}

func (tews *testErrorWithStackFrames) StackFrames() []StackFrame {
	return tews.Err.StackFrames()
}

func (tews *testErrorWithStackFrames) Error() string {
	return tews.Err.Error()
}

func TestNewError(t *testing.T) {

	e := func() error {
		return New("hi", 1)
	}()

	if e.Error() != "hi" {
		t.Errorf("Constructor with a string failed")
	}

	if New(fmt.Errorf("yo"), 0).Error() != "yo" {
		t.Errorf("Constructor with an error failed")
	}

	if New(e, 0) != e {
		t.Errorf("Constructor with an Error failed")
	}

	if New(nil, 0).Error() != "<nil>" {
		t.Errorf("Constructor with nil failed")
	}

	err := New("foo", 0)
	tews := &testErrorWithStackFrames{
		Err: err,
	}

	if bytes.Compare(New(tews, 0).Stack(), err.Stack()) != 0 {
		t.Errorf("Constructor with ErrorWithStackFrames failed")
	}
}

func ExampleErrorf() {
	for i := 1; i <= 2; i++ {
		if i%2 == 1 {
			e := Errorf("can only halve even numbers, got %d", i)
			fmt.Printf("Error: %+v", e)
		}
	}
	// Output:
	// Error: can only halve even numbers, got 1
}

func ExampleNew() {
	// Wrap io.EOF with the current stack-trace and return it
	e := New(io.EOF, 0)
	fmt.Printf("%+v", e)
	// Output:
	// EOF
}

func ExampleNew_skip() {
	defer func() {
		if err := recover(); err != nil {
			// skip 1 frame (the deferred function) and then return the wrapped err
			err = New(err, 1)
		}
	}()
}

func ExampleError_Stack() {
	e := New("Oh noes!", 1)
	fmt.Printf("Error: %s\n", e.Error())
	fmt.Printf("Stack is %d bytes", len(e.Stack()))
	// Output:
	// Error: Oh noes!
	// Stack is 589 bytes
}
