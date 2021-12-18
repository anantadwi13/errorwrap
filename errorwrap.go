package errorwrap

import (
	"errors"
	"fmt"
	"io"
)

type errorContainer struct {
	error
	contextMsg  string
	rootCause   error
	parentError error
	*stack
}

func (e *errorContainer) Error() string {
	str := e.error.Error()
	if e.contextMsg != "" {
		str += ": " + e.contextMsg
	}
	return str
}

func (e *errorContainer) fullError() string {
	str := e.Error()
	if e.parentError != nil {
		str += "\n"
		if ec, ok := e.parentError.(*errorContainer); ok && ec != nil {
			str += ec.fullError()
		} else {
			str += e.parentError.Error()
		}
	}
	return str
}

func (e *errorContainer) Unwrap() error {
	return e.parentError
}

func (e *errorContainer) Is(target error) bool {
	return e.error == target
}

func (e *errorContainer) As(target interface{}) bool {
	if target == nil {
		return e.error == target
	}
	if tErr, ok := target.(error); ok && e.error == tErr {
		return true
	}

	return false
}

func (e *errorContainer) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.fullError())
			if rc, ok := e.rootCause.(*errorContainer); ok && rc != nil {
				rc.stack.Format(s, verb)
			} else {
				e.stack.Format(s, verb)
			}
			return
		}
		fallthrough
	case 's':
		if s.Flag('+') {
			io.WriteString(s, e.fullError())
			return
		}
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func New(message string) error {
	return &errorContainer{
		error: errors.New(message),
		stack: callStack(),
	}
}

func Wrap(parent error, wrapper error) error {
	ec := WrapWithMessage(parent, wrapper, "")
	ec.(*errorContainer).parentError.(*errorContainer).stack = callStack()
	return ec
}

func WrapWithMessage(parent error, wrapper error, contextMessage string) error {
	parentContainer, ok := parent.(*errorContainer)
	var rootCause error

	if !ok {
		parentContainer = &errorContainer{
			error: parent,
			stack: callStack(),
		}
		rootCause = parentContainer
	} else {
		rootCause = parentContainer.rootCause
	}

	currError := &errorContainer{
		error:       wrapper,
		contextMsg:  contextMessage,
		rootCause:   rootCause,
		parentError: parentContainer,
		stack:       callStack(),
	}
	return currError
}

func WrapString(parent error, wrapperMsg string) error {
	ec := Wrap(parent, errors.New(wrapperMsg))
	ec.(*errorContainer).parentError.(*errorContainer).stack = callStack()
	return ec
}

func Wrapper(err error, target error) error {
	if target == nil || err == nil {
		return nil
	}

	for {
		if curContainer, ok := err.(*errorContainer); ok && curContainer.Is(target) {
			return err
		}

		if err = Unwrap(err); err == nil {
			return nil
		}
	}
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
