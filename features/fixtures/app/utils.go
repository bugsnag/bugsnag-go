package main

import (
	"fmt"
	"runtime"
)

type CustomErr struct {
	msg     string
	cause   error
	callers []uintptr
}

func NewCustomErr(msg string, cause error) error {
	callers := make([]uintptr, 8)
	runtime.Callers(2, callers)
	return CustomErr{
		msg:     msg,
		cause:   cause,
		callers: callers,
	}
}

func (err CustomErr) Error() string {
	return err.msg
}

func (err CustomErr) Unwrap() error {
	return err.cause
}

func (err CustomErr) Callers() []uintptr {
	return err.callers
}

func Login(token string) error {
	val, err := CheckValue(len(token) * -1)
	if err != nil {
		return NewCustomErr("login failed", err)
	}
	fmt.Printf("val: %d\n", val)
	return nil
}

func CheckValue(i int) (int, error) {
	if i < 0 {
		return 0, NewCustomErr("invalid token", nil)
	} else if i%2 == 0 {
		return i / 2, nil
	} else if i < 9 {
		return i * 3, nil
	}

	return i * 4, nil
}