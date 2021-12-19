package errorwrap

import (
	"errors"
	"fmt"
	"io"
)

// ErrorWrapper is a multi level error. It is built like a stack (from bottom to top). If there is a 3 level, then each
// level will have ErrorWrapper instance with an error (CurrentError).
//
//      ErrorWrapper Z (CurrentError contains any object that implements error, RootCause contains ErrorWrapper X, ParentError contains ErrorWrapper Y)
//           V
//      ErrorWrapper Y (CurrentError contains any object that implements error, RootCause contains ErrorWrapper X, ParentError contains ErrorWrapper X)
//           V
//      ErrorWrapper X (CurrentError contains any object that implements error, RootCause contains nil, ParentError contains nil)
//
type ErrorWrapper interface {
	// CurrentError represents as the current level error (the top level of the stack).
	CurrentError() error
	// ContextMessage is used as a modifier to the CurrentError. It is used for elaborating the CurrentError.
	ContextMessage() string
	// RootCause will return the bottom level of the stack.
	RootCause() ErrorWrapper
	// ParentError will return an ErrorWrapper below the current level.
	ParentError() ErrorWrapper
	// StackTrace return a StackTrace of the current ErrorWrapper level only.
	StackTrace() StackTrace

	error
	Unwrap() error
	Is(target error) bool
	As(target interface{}) bool
	fmt.Formatter
}

type errorWrapper struct {
	error
	contextMsg  string
	rootCause   ErrorWrapper
	parentError ErrorWrapper
	*stack
}

func (e *errorWrapper) CurrentError() error {
	return e.error
}

func (e *errorWrapper) ContextMessage() string {
	return e.contextMsg
}

func (e *errorWrapper) RootCause() ErrorWrapper {
	return e.rootCause
}

func (e *errorWrapper) ParentError() ErrorWrapper {
	return e.parentError
}

func (e *errorWrapper) StackTrace() StackTrace {
	return e.stack.StackTrace()
}

func (e *errorWrapper) Error() string {
	str := e.error.Error()
	if e.contextMsg != "" {
		str += ": " + e.contextMsg
	}
	return str
}

func (e *errorWrapper) fullError() string {
	str := e.Error()
	if e.parentError != nil {
		str += "\n"
		if ec, ok := e.parentError.(*errorWrapper); ok && ec != nil {
			str += ec.fullError()
		} else {
			str += e.parentError.Error()
		}
	}
	return str
}

func (e *errorWrapper) Unwrap() error {
	return e.parentError
}

func (e *errorWrapper) Is(target error) bool {
	return e.error == target
}

func (e *errorWrapper) As(target interface{}) bool {
	if target == nil {
		return e.error == target
	}
	if tErr, ok := target.(error); ok && e.error == tErr {
		return true
	}

	return false
}

func (e *errorWrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.fullError())
			if rc, ok := e.rootCause.(*errorWrapper); ok && rc != nil {
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

// New creates a base or root ErrorWrapper.
//
// ErrorWrapper.CurrentError is populated with errors.New(message).
func New(message string) error {
	return &errorWrapper{
		error: errors.New(message),
		stack: callStack(),
	}
}

// NewError creates a base or root ErrorWrapper.
//
// ErrorWrapper.CurrentError is populated with err. If err is nil then NewError returns nil.
func NewError(err error) error {
	if err == nil {
		return nil
	}

	if ew, ok := err.(*errorWrapper); ok {
		if ew == nil {
			return nil
		}
		return ew
	}

	return &errorWrapper{
		error: err,
		stack: callStack(),
	}
}

// NewErrorWithMessage creates a base or root ErrorWrapper.
//
// ErrorWrapper.CurrentError is populated with err.
// ErrorWrapper.ContextMessage is populated with contextMessage.
func NewErrorWithMessage(err error, contextMessage string) error {
	if err == nil {
		return nil
	}

	if ew, ok := err.(*errorWrapper); ok {
		if ew == nil {
			return nil
		}
		ew.contextMsg = contextMessage
		return ew
	}

	return &errorWrapper{
		error:      err,
		contextMsg: contextMessage,
		stack:      callStack(),
	}
}

// Wrap returns an ErrorWrapper{CurrentError: wrapper, ParentError: parent, RootCause: parent.RootCause}
func Wrap(parent error, wrapper error) error {
	ec := WrapWithMessage(parent, wrapper, "")

	if ec == nil {
		return nil
	}

	if _, ok := parent.(*errorWrapper); !ok {
		ec.(*errorWrapper).parentError.(*errorWrapper).stack = callStack()
	}
	ec.(*errorWrapper).stack = callStack()
	return ec
}

// WrapWithMessage returns an ErrorWrapper{CurrentError: wrapper, ParentError: parent, RootCause: parent.RootCause, ContextMessage: contextMessage}
func WrapWithMessage(parent error, wrapper error, contextMessage string) error {
	if parent == nil || wrapper == nil {
		return nil
	}

	parentContainer, ok := parent.(*errorWrapper)
	if !ok {
		parentContainer = &errorWrapper{
			error: parent,
			stack: callStack(),
		}
	}

	var rootCause ErrorWrapper

	if parentContainer.rootCause == nil {
		rootCause = parentContainer
	} else {
		rootCause = parentContainer.rootCause
	}

	currError := &errorWrapper{
		error:       wrapper,
		contextMsg:  contextMessage,
		rootCause:   rootCause,
		parentError: parentContainer,
		stack:       callStack(),
	}
	return currError
}

// Wrapper returns ErrorWrapper instance of err in target level
func Wrapper(err error, target error) ErrorWrapper {
	if target == nil || err == nil {
		return nil
	}

	for {
		if curContainer, ok := err.(ErrorWrapper); ok && (curContainer == target || curContainer.Is(target)) {
			return curContainer
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
