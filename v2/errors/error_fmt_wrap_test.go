// +build go1.13

package errors

import (
	"fmt"
	"runtime"
	"testing"
)

func TestUnwrapErrorsCause(t *testing.T) {
	_, _, line, ok := runtime.Caller(0) // grab line immediately before error generators
	err1 := fmt.Errorf("invalid token")
	err2 := fmt.Errorf("login failed: %w", err1)
	err3 := fmt.Errorf("terminate process: %w", err2)
	unwrapped := New(err3, 0)
	if !ok {
		t.Fatalf("Something has gone wrong with loading the current stack")
	}
	if unwrapped.Error() != "terminate process: login failed: invalid token" {
		t.Errorf("Failed to unwrap error: %s", unwrapped.Error())
	}
	assertStacksMatch(t, []StackFrame{
		StackFrame{Name: "TestUnwrapErrorsCause", File: "errors/error_fmt_wrap_test.go", LineNumber: line + 4},
	}, unwrapped.StackFrames())
	if unwrapped.Cause == nil {
		t.Fatalf("Failed to capture cause error")
	}
	if unwrapped.Cause.Error() != "login failed: invalid token" {
		t.Errorf("Failed to unwrap cause error: %s", unwrapped.Cause.Error())
	}
	if len(unwrapped.Cause.StackFrames()) > 0 {
		t.Errorf("Did not expect cause to have a stack: %v", unwrapped.Cause.StackFrames())
	}
	if unwrapped.Cause.Cause == nil {
		t.Fatalf("Failed to capture nested cause error")
	}
	if len(unwrapped.Cause.Cause.StackFrames()) > 0 {
		t.Errorf("Did not expect cause to have a stack: %v", unwrapped.Cause.Cause.StackFrames())
	}
	if unwrapped.Cause.Cause.Cause != nil {
		t.Fatalf("Extra cause detected: %v", unwrapped.Cause.Cause.Cause)
	}
}
