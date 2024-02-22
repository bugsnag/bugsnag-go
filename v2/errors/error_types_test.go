package errors

import (
	"fmt"
	"testing"

	pkgerror "github.com/pkg/errors"
)

const (
	bugsnagType     = "bugsnagError"
	callersType     = "callersType"
	stackFramesType = "stackFramesType"
	stackType       = "stackType"
	internalType    = "internalType"
	stringType      = "stringType"
)

// Prepared to test inlining
func AInternal() interface{} { return fmt.Errorf("pure golang error") }
func BInternal() interface{} { return AInternal() }
func CInternal() interface{} { return BInternal() }

func AString() interface{} { defer func() interface{} { return recover() }(); panic("panic") }
func BString() interface{} { return AString() }
func CString() interface{} { return BString() }

func AStack() interface{} { return pkgerror.Errorf("from package") }
func BStack() interface{} { return AStack() }
func CStack() interface{} { return BStack() }

func ACallers() interface{} { return newCustomErr("oh no an error", fmt.Errorf("parent error")) }
func BCallers() interface{} { return ACallers() }
func CCallers() interface{} { return BCallers() }

func AFrames() interface{} { return &testErrorWithStackFrames{Err: New("foo", 0)} }
func BFrames() interface{} { return AFrames() }
func CFrames() interface{} { return BFrames() }

// Golang internal errors don't have stacktrace
// StackFrames are going to report only the line where internal golang error was wrapped in Bugsnag error
func TestInternalError(t *testing.T) {
	err := CInternal()
	typeAssert(t, err, internalType)

	bgError := New(err, 0)
	actualStack := bgError.StackFrames()
	expected := []StackFrame{
		{Name: "TestInternalError", File: "errors/error_types_test.go", LineNumber: 46},
	}
	assertStacksMatch(t, expected, actualStack)
}

// Errors from panic contain only the message about panic
// Same as above - StackFrames are going to contain only line numer of wrapping
func TestStringError(t *testing.T) {
	err := CString()
	typeAssert(t, err, stringType)

	bgError := New(err, 0)
	actualStack := bgError.StackFrames()
	expected := []StackFrame{
		{Name: "TestStringError", File: "errors/error_types_test.go", LineNumber: 60},
	}
	assertStacksMatch(t, expected, actualStack)
}

// Errors from pkg/errors have their own stack
// Inlined functions should be visible in StackFrames
func TestStackError(t *testing.T) {
	err := CStack()
	typeAssert(t, err, stackType)

	bgError := New(err, 0)
	actualStack := bgError.StackFrames()
	expected := []StackFrame{
		{Name: "AStack", File: "errors/error_types_test.go", LineNumber: 28},
		{Name: "BStack", File: "errors/error_types_test.go", LineNumber: 29},
		{Name: "CStack", File: "errors/error_types_test.go", LineNumber: 30},
		{Name: "TestStackError", File: "errors/error_types_test.go", LineNumber: 71},
	}

	assertStacksMatch(t, expected, actualStack)
}

// Errors implementing Callers() interface should have their own stack
// Inlined functions should be visible in StackFrames
func TestCallersError(t *testing.T) {
	err := CCallers()
	typeAssert(t, err, callersType)

	bgError := New(err, 0)
	actualStack := bgError.StackFrames()
	expected := []StackFrame{
		{Name: "ACallers", File: "errors/error_types_test.go", LineNumber: 32},
		{Name: "BCallers", File: "errors/error_types_test.go", LineNumber: 33},
		{Name: "CCallers", File: "errors/error_types_test.go", LineNumber: 34},
		{Name: "TestCallersError", File: "errors/error_types_test.go", LineNumber: 89},
	}
	assertStacksMatch(t, expected, actualStack)
}

// Errors with StackFrames are explicilty adding stacktrace to error
// Inlined functions should be visible in StackFrames
func TestFramesError(t *testing.T) {
	err := CFrames()
	typeAssert(t, err, stackFramesType)

	bgError := New(err, 0)
	actualStack := bgError.StackFrames()
	expected := []StackFrame{
		{Name: "AFrames", File: "errors/error_types_test.go", LineNumber: 36},
		{Name: "BFrames", File: "errors/error_types_test.go", LineNumber: 37},
		{Name: "CFrames", File: "errors/error_types_test.go", LineNumber: 38},
		{Name: "TestFramesError", File: "errors/error_types_test.go", LineNumber: 106},
	}

	assertStacksMatch(t, expected, actualStack)
}

func typeAssert(t *testing.T, err interface{}, expectedType string) {
	actualType := checkType(err)
	if actualType != expectedType {
		t.Errorf("Types don't match. Actual: %+v and expected: %+v\n", actualType, expectedType)
	}
}

func checkType(err interface{}) string {
	switch err.(type) {
	case *Error:
		return bugsnagType
	case ErrorWithCallers:
		return callersType
	case errorWithStack:
		return stackType
	case ErrorWithStackFrames:
		return stackFramesType
	case error:
		return internalType
	default:
		return stringType
	}
}
