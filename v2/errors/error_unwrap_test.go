// +build go1.13

package errors

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestFindingErrorInChain(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := errors.Wrap(baseErr, "failed")
	err := New(wrappedErr, 0)

	if !errors.Is(err, baseErr) {
		t.Errorf("Failed to find base error: %s", err.Error())
	}
}

func TestErrorUnwrapping(t *testing.T) {
	baseErr := errors.New("base error")
	wrappedErr := fmt.Errorf("failed: %w", baseErr)
	err := New(wrappedErr, 0)

	unwrapped := errors.Unwrap(err)

	if unwrapped != baseErr {
		t.Errorf("Failed to find base error: %s", unwrapped.Error())
	}
}
