package wp

import (
	"fmt"
	"strings"
)

type MultiError struct {
	errorString string
}

func (e *MultiError) Error() string {
	return strings.TrimSpace(e.errorString)
}

func (e *MultiError) Exists() bool {
	return e.errorString != ""
}

func (e *MultiError) Append(err error) {
	e.errorString = fmt.Sprintf("%s\n%s", e.errorString, err.Error())
}

func MultiErrorFromErrors(errors []error) *MultiError {
	me := &MultiError{""}
	for _, err := range errors {
		if err != nil {
			me.Append(err)
		}
	}
	return me
}
