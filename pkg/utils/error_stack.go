package utils

import (
	"bytes"
	"strconv"
)

func NewErrorStack() *ErrorStack {
	return &ErrorStack{
		errors: make([]error, 0),
	}
}

type ErrorStack struct {
	errors []error
}

func (e *ErrorStack) Error() string {
	buffer := new(bytes.Buffer)

	for i, err := range e.errors {
		buffer.WriteString(strconv.Itoa(i) + ") " + err.Error() + "\n")
	}

	return buffer.String()
}

func (e *ErrorStack) Add(err error) {
	if err == nil {
		return
	}

	if es, ok := err.(*ErrorStack); ok {
		e.errors = append(e.errors, es.errors...)
		return
	}

	e.errors = append(e.errors, err)
}

func (e *ErrorStack) Get() error {
	if len(e.errors) == 0 {
		return nil
	}

	return e
}
