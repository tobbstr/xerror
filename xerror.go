/*
Package xerror provides a way to wrap errors with additional context and to add variables to the error that can be
logged at a later time. It also provides a way to categorize errors into different kinds.
*/
package xerror

import (
	"errors"
	"io"
)

// Severity is used to control the way the error is logged. For example as an error, warning, notice etc.
type Severity uint8

const (
	SeverityUnspecified Severity = iota
	SeverityDebug
	SeverityInfo
	SeverityWarn
	SeverityError
)

const (
	KindUnsupported Kind = iota
	KindInvalidArgs
	KindPreconditionViolation
	KindOutOfRange
	KindUnauthenticated
	KindPermissionDenied
	KindNotFound
	KindAborted
	KindAlreadyExists
	KindResourceExhausted
	KindCancelled
	KindDataLoss
	KindUnknown
	KindInternal
	KindNotImplemented
	KindUnavailable
	KindDeadlineExceeded
)

// Kind is used to categorize errors into different kinds. It's used to provide more context to the error.
type Kind uint8

func (k Kind) String() string {
	switch k {
	case KindUnsupported:
		return "unsupported"
	case KindInvalidArgs:
		return "invalid_arguments"
	case KindPreconditionViolation:
		return "precondition_violation"
	case KindOutOfRange:
		return "out of range"
	case KindUnauthenticated:
		return "unauthenticated"
	case KindPermissionDenied:
		return "permission_denied"
	case KindNotFound:
		return "not_found"
	case KindAborted:
		return "aborted"
	case KindAlreadyExists:
		return "already_exists"
	case KindResourceExhausted:
		return "resource_exhausted"
	case KindCancelled:
		return "cancelled"
	case KindDataLoss:
		return "data_loss"
	case KindUnknown:
		return "unknown"
	case KindInternal:
		return "internal"
	case KindNotImplemented:
		return "not_implemented"
	case KindUnavailable:
		return "unavailable"
	case KindDeadlineExceeded:
		return "deadline_exceeded"
	default:
		return "unsupported"
	}
}

// Var models what the circumstances were when the error was encountered and is used to provide additional context
// to the error. Its purpose is to be logged and thereby give context to the error in the logs.
type Var struct {
	Name  string
	Value any
}

// Root models a root error and is used when you first encounter an error in your code.
// The idiomatic way to use this type is to create a new instance by using the "new" built-in function.
//
//		Ex.
//	 	err := errors.NewRoot(err, "more context to err").WithVar("id", id).WithKind(errors.NotFound)
//
//nolint:errname
type Root struct {
	// Err is the wrapped error.
	Err error
	// Kind is the kind of error that was encountered.
	Kind Kind
	// RuntimeState is a snapshot of the state of the application when the error was encountered.
	RuntimeState []Var
	// Retryable is a flag that indicates whether the error is retryable.
	Retryable bool
	// Severity is the severity of the error used to control the way the error is logged. For example as an error,
	// warning, notice etc.
	Severity Severity
}

// NewRoot creates a new Root error. It's a convenience function that allows you to create a Root error with an original
// error and a message (or message chain) that provides more context to the error.
// If you don't want to provide a message, you can omit it. If you provide multiple messages,
// they will be wrapped in a chain of WrappedError instances where the first message is the innermost message.
//
// Ex. NewRoot(err, "message1", "message2", "message3") will create a chain of WrappedError instances where "message1"
// is the innermost message and "message3" is the outermost message. When you call Error() on the Root instance, you'll
// get "message3: message2: message1: original error message".
func NewRoot(original error, messages ...string) *Root {
	if len(messages) == 0 {
		return &Root{Err: original}
	}
	var root Root
	for i, msg := range messages {
		if i == 0 {
			root.Err = &WrappedError{Msg: msg, Err: original}
			continue
		}
		root.Err = &WrappedError{Msg: msg, Err: root.Err}
	}
	return &root
}

// Error returns the error message of the wrapped error that was encountered.
//
// If you have a Root instance whose 'Err' field value is a WrappedError, calling Error() on the Root instance will
// return the message of the WrappedError instance.
// Ex. Calling root.Error() will return "message1: original error message", given that the WrappedError instance
// has the message "message1".
func (r *Root) Error() string {
	if r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

// Original returns the original error that was encountered. It unwraps the error chain and returns the original.
func (r *Root) Original() error {
	err := r.Err
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped != nil {
			err = unwrapped
			continue
		}
		return err
	}
}

// WithError sets the Root's original error. It must only be called once.
func (r *Root) WithError(err error) *Root {
	r.Err = err
	return r
}

// WithVar adds a variable to the RuntimeState slice. It should only be used on variables that will still be in scope
// at log-time. If the variable won't be in scope at log-time, it must not be added to the RuntimeState slice.
//
// For out-of-scope variables, there are other methods that can be used to add them to the RuntimeState slice such as
// 'WithIOReaderVar', which is provided for io.Reader variables.
//
// This method may be called multiple times.
func (r *Root) WithVar(name string, value any) *Root {
	if name == "" || value == nil {
		return r
	}
	r.RuntimeState = append(r.RuntimeState, Var{Name: name, Value: value})
	return r
}

// WithVars allows you to add multiple variables to the RuntimeState slice. It's the bulk-variant of WithVar.
func (r *Root) WithVars(vars ...Var) *Root {
	for _, v := range vars {
		_ = r.WithVar(v.Name, v.Value)
	}
	return r
}

// WithIOReaderVar reads the io.Reader and adds the contents to the RuntimeState slice.
func (r *Root) WithIOReaderVar(name string, value io.Reader) *Root {
	if name == "" || value == nil {
		return r
	}
	b, err := io.ReadAll(value)
	if err != nil {
		return r
	}
	return r.WithVar(name, string(b))
}

// WithKind sets the Root's kind. It must only be called once.
func (r *Root) WithKind(kind Kind) *Root {
	r.Kind = kind
	return r
}

// WithRetry sets the Root's Retryable flag to true.
func (r *Root) WithRetry() *Root {
	r.Retryable = true
	return r
}

// Unwrap returns the original error that was encountered.
func (r *Root) Unwrap() error {
	return r.Err
}

// WithSeverity sets the Root's severity.
func (r *Root) WithSeverity(s Severity) *Root {
	r.Severity = s
	return r
}

// Is traverses err's error chain and compares the first encountered Root error for equality. If equal, it returns true.
// If no such error is found or if they're not equal, it returns false.
func (r *Root) Is(err error) bool {
	var root *Root
	if errors.As(err, &root) {
		return root.Kind == r.Kind && root.Err.Error() == r.Err.Error()
	}
	return false
}

type WrappedError struct {
	Msg string
	Err error
}

func (e *WrappedError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return e.Msg + ": " + e.Err.Error()
}

func (e *WrappedError) Unwrap() error {
	return e.Err
}

// Wrap is a utility function that makes it easier to wrap errors with a message to add more context.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if msg == "" {
		return err
	}
	return &WrappedError{Msg: msg, Err: err}
}

// Is traverses err's error chain and compares the first encountered Root error's Kind. If equal, it returns true.
// If no such error is found or if their kinds differ, it returns false.
func Is(err error, kind Kind) bool {
	var root *Root
	if errors.As(err, &root) {
		return root.Kind == kind
	}
	return false
}

// AddVars adds variables to the RuntimeState slice of a Root error. It's a utility function that makes it easier to
// add variables to the RuntimeState slice of a Root error.
func AddVars(err error, vars ...Var) {
	var root *Root
	if errors.As(err, &root) {
		root.RuntimeState = append(root.RuntimeState, vars...)
	}
}

// IsRetryable returns true if the error is retryable, otherwise it returns false.
func IsRetryable(err error) bool {
	var root *Root
	if errors.As(err, &root) {
		return root.Retryable
	}
	return false
}

// MakeRetryable sets the error's Retryable flag to false.
func ClearRetryable(err error) error {
	var root *Root
	if errors.As(err, &root) {
		root.Retryable = false
	}
	return err
}
