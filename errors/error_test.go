package errors

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// fixture functions doing work to avoid inlining
func a(i int) error {
	if b(i + 5) && b(i + 6) {
		return nil
	}
	return fmt.Errorf("not gonna happen")
}

func b(i int) bool {
	return c(i+2) > 12
}

// panicking function!
func c(i int) int {
	if i > 3 {
		panic('a')
	}
	return i * i
}

func TestParsePanicStack(t *testing.T) {
	defer func() {
		err := New(recover(), 0)
		if err.Error() != "97" {
			t.Errorf("Received incorrect error, expected 'a' got '%s'", err.Error())
		}
		if err.TypeName() != "*errors.errorString" {
			t.Errorf("Error type was '%s'", err.TypeName())
		}
		for index, frame := range err.StackFrames() {
			if frame.Func() == nil {
				t.Errorf("Failed to remove nil frame %d", index)
			}
		}
		expected := []StackFrame{
			StackFrame{Name: "TestParsePanicStack.func1", File: "errors/error_test.go"},
			StackFrame{Name: "a", File: "errors/error_test.go", LineNumber: 13},
		}
		assertStacksMatch(t, expected, err.StackFrames())
	}()

	a(1)
}

func TestParseGeneratedStack(t *testing.T) {
	err := New(fmt.Errorf("e_too_many_colander"), 0)
	expected := []StackFrame{
		StackFrame{Name: "TestParseGeneratedStack", File: "errors/error_test.go"},
	}
	if err.Error() != "e_too_many_colander" {
		t.Errorf("Error name was '%s'", err.Error())
	}
	if err.TypeName() != "*errors.errorString" {
		t.Errorf("Error type was '%s'", err.TypeName())
	}
	for index, frame := range err.StackFrames() {
		if frame.Func() == nil {
			t.Errorf("Failed to remove nil frame %d", index)
		}
	}
	assertStacksMatch(t, expected, err.StackFrames())
}

func TestSkipWorks(t *testing.T) {
	defer func() {
		err := New(recover(), 1)
		if err.Error() != "97" {
			t.Errorf("Received incorrect error, expected 'a' got '%s'", err.Error())
		}

		for index, frame := range err.StackFrames() {
			if frame.Name == "TestSkipWorks.func1" {
				t.Errorf("Failed to skip frame")
			}
			if frame.Func() == nil {
				t.Errorf("Failed to remove inlined frame %d", index)
			}
		}

		expected := []StackFrame{
			StackFrame{Name: "a", File: "errors/error_test.go", LineNumber: 13},
		}

		assertStacksMatch(t, expected, err.StackFrames())
	}()

	a(4)
}

func checkFramesMatch(expected StackFrame, actual StackFrame) bool {
	if actual.Name != expected.Name {
		return false
	}
	// Not using exact match as it would change depending on whether
	// the package is being tested within or outside of the $GOPATH
	if expected.File != "" && !strings.HasSuffix(actual.File, expected.File) {
		return false
	}
	if expected.Package != "" && actual.Package != expected.Package {
		return false
	}
	if expected.LineNumber != 0 && actual.LineNumber != expected.LineNumber {
		return false
	}
	return true
}

func assertStacksMatch(t *testing.T, expected []StackFrame, actual []StackFrame) {
	var lastmatch int = 0
	var matched int = 0
	// loop over the actual stacktrace, checking off expected frames as they
	// are found. Each one might be in the middle of the stack, but the order
	// should remain the same.
	for _, actualFrame := range actual {
		for index, expectedFrame := range expected {
			if index < lastmatch {
				continue
			}
			if checkFramesMatch(expectedFrame, actualFrame) {
				lastmatch = index
				matched += 1
				break
			}
		}
	}
	if matched != len(expected) {
		t.Fatalf("failed to find matches for %d frames: '%v'\ngot: '%v'", len(expected)-matched, expected[matched:], actual)
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
