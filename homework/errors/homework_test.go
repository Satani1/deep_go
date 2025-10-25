package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	flatError := fmt.Sprintf("%d errors occured:\n", len(e.Errors))

	for _, err := range e.Errors {
		flatError += fmt.Sprintf("\t* %s", err.Error())
	}
	flatError += "\n"

	return flatError
}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	var myErr *MultiError
	if errors.As(err, &myErr) {
		myErr.Errors = append(myErr.Errors, errs...)
		return myErr
	}

	errors := make([]error, 0, len(errs)+1)
	if err != nil {
		errors = append(errors, err)
	}

	errors = append(errors, errs...)

	return &MultiError{
		Errors: errors,
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
