package errorwrap

import (
	"errors"
	"fmt"
	"io"
)

var (
	multilineSeparator = " -  "
	multilineIndent    = "    "
)

// ErrorDefinition is a definition of an error. It is used like an immutable variable.
type ErrorDefinition interface {
	error
	fmt.Formatter
}

// ErrorWrapper is a multi level error. It is built like a stack (vertical from bottom to top) that contains multiple
// ErrorDefinition (horizontal) in each level. If there is a 3 level, then each level will have an ErrorWrapper
// instance with one ErrorDefinition (CurrentError) minimum.
//
//      ErrorWrapper Z (CurrentError contains any objects that implement error, RootCause contains ErrorWrapper X, ParentError contains ErrorWrapper Y)
//           V
//      ErrorWrapper Y (CurrentError contains any objects that implement error, RootCause contains ErrorWrapper X, ParentError contains ErrorWrapper X)
//           V
//      ErrorWrapper X (CurrentError contains any objects that implement error, RootCause contains nil, ParentError contains nil)
//
type ErrorWrapper interface {
	// CurrentError represents as the current level error (the top level of the stack).
	CurrentError() []error
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

type errorDefinition struct {
	msg string
}

func (e *errorDefinition) Error() string {
	return e.msg
}

func (e *errorDefinition) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

type errorWrapper struct {
	errors      []error
	contextMsg  string
	rootCause   ErrorWrapper
	parentError ErrorWrapper
	*stack
}

func (e *errorWrapper) CurrentError() []error {
	return e.errors
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
	str := ""
	if len(e.errors) <= 0 {
		return ""
	} else if len(e.errors) == 1 {
		str += multilineSeparator + e.errors[0].Error()
	} else {
		for i, err := range e.errors {
			switch i {
			case 0:
				str += multilineSeparator + err.Error()
			default:
				str += "\n" + multilineIndent + err.Error()
			}
		}
	}

	if e.contextMsg != "" {
		str += "\n" + multilineIndent + "context: " + e.contextMsg
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
	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *errorWrapper) As(target interface{}) bool {
	if target == nil {
		if e.errors == nil {
			return true
		}
		return false
	}
	for _, err := range e.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

func (e *errorWrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n\n", e.fullError())
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

// New creates an ErrorDefinition
func New(message string) error {
	return &errorDefinition{
		msg: message,
	}
}

func newErrorWrapper(err ...error) *errorWrapper {
	var errs []error
	for _, e := range err {
		if e == nil {
			continue
		}
		errs = append(errs, e)
	}

	if len(errs) == 0 {
		return nil
	}

	return &errorWrapper{
		errors: errs,
	}
}

// NewError creates a base or root ErrorWrapper.
//
// ErrorWrapper.CurrentError is populated with err. If err is nil then NewError returns nil.
// It is recommended to pass ErrorDefinition as err arguments.
func NewError(err ...error) error {
	errWrap := newErrorWrapper(err...)
	if errWrap == nil {
		return nil
	}
	errWrap.stack = callStack()
	return errWrap
}

// NewErrorWithMessage creates a base or root ErrorWrapper.
//
// ErrorWrapper.CurrentError is populated with err.
// ErrorWrapper.ContextMessage is populated with contextMessage.
// It is recommended to pass ErrorDefinition as err arguments.
func NewErrorWithMessage(contextMessage string, err ...error) error {
	errWrap := newErrorWrapper(err...)
	if errWrap == nil {
		return nil
	}
	errWrap.stack = callStack()
	errWrap.contextMsg = contextMessage
	return errWrap
}

// AppendInto appends err into errWrapper. It will return the same errWrapper instance or new instance if errWrapper is nil.
// It is recommended to pass ErrorDefinition as err arguments.
func AppendInto(errWrapper error, err ...error) error {
	ew, ok := errWrapper.(*errorWrapper)
	if ok && ew != nil {
		ew.errors = append(ew.errors, err...)
	} else {
		ew = newErrorWrapper(err...)
		if ew == nil {
			return nil
		}
		ew.stack = callStack()
	}
	return ew
}

// Wrap returns an ErrorWrapper{CurrentError: wrapper, ParentError: parent, RootCause: parent.RootCause}.
// It is recommended to pass ErrorDefinition as err arguments.
func Wrap(parent error, err ...error) error {
	ec := WrapWithMessage(parent, "", err...)
	if ec == nil {
		return nil
	}

	if p, ok := parent.(*errorWrapper); !ok || p == nil {
		if parentErr, ok2 := ec.(*errorWrapper).parentError.(*errorWrapper); ok2 && parentErr != nil {
			parentErr.stack = callStack()
		}
	}
	ec.(*errorWrapper).stack = callStack()
	return ec
}

// WrapWithMessage returns an ErrorWrapper{CurrentError: wrapper, ParentError: parent, RootCause: parent.RootCause, ContextMessage: contextMessage}.
// It is recommended to pass ErrorDefinition as err arguments.
func WrapWithMessage(parent error, contextMessage string, err ...error) error {
	if parent == nil && len(err) == 0 {
		return nil
	}

	parentWrapper, ok := parent.(*errorWrapper)
	if !ok {
		parentWrapper = newErrorWrapper(parent)
		if parentWrapper != nil {
			parentWrapper.stack = callStack()
		}
	}

	currError := newErrorWrapper(err...)
	if currError == nil {
		if parentWrapper != nil {
			return parentWrapper
		}
		return nil
	}

	if parentWrapper != nil {
		if parentWrapper.rootCause == nil {
			currError.rootCause = parentWrapper
		} else {
			currError.rootCause = parentWrapper.rootCause
		}
		currError.parentError = parentWrapper
	}

	currError.contextMsg = contextMessage
	currError.stack = callStack()

	return currError
}

// Wrapper returns ErrorWrapper instance of err in target level
func Wrapper(err error, target error) ErrorWrapper {
	if target == nil || err == nil {
		return nil
	}

	for {
		if curWrapper, ok := err.(ErrorWrapper); ok && (curWrapper == target || curWrapper.Is(target)) {
			return curWrapper
		}

		if err = Unwrap(err); err == nil {
			return nil
		}
	}
}

// Is checks whether target is placed in err (ErrorWrapper) or not. It will find recursively from current level
// to the root of ErrorWrapper stack.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// IsExact checks whether target is placed in current level error of err (ErrorWrapper) or not.
func IsExact(err error, target error) bool {
	if curWrapper, ok := err.(ErrorWrapper); ok && (curWrapper == target || curWrapper.Is(target)) {
		return true
	}
	return false
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap will return the parent error (ErrorWrapper) of err (ErrorWrapper)
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
