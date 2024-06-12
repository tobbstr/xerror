/*
Package xerror provides a way to wrap errors with additional context and to add variables to the error that can be
logged at a later time. It also provides a way to categorize errors into different kinds.
*/
package xerror

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

/*
TODO:
	1. Make it easy to consume gRPC errors
	2. Make it easy to consume HTTP errors by generating typescript models from status.Status
	3. Make it easy to produce HTTP errors
	4. Make it easy to respond with HTTP errors
	5. Make it easy to produce gRPC errors
	6. Make it easy to respond with gRPC errors (unaryinterceptor)
*/

// LogLevel is used to control the way the error is logged. For example as an error, warning, notice etc.
type LogLevel uint8

const (
	LogLevelUnspecified LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// const (
// 	KindUnsupported Kind = iota
// 	KindInvalidArgs
// 	KindPreconditionViolation
// 	KindOutOfRange
// 	KindUnauthenticated
// 	KindPermissionDenied
// 	KindNotFound
// 	KindAborted
// 	KindAlreadyExists
// 	KindResourceExhausted
// 	KindCanceled
// 	KindDataLoss
// 	KindUnknown
// 	KindInternal
// 	KindUnimplemented
// 	KindUnavailable
// 	KindDeadlineExceeded
// )

// // Kind is used to categorize errors into different kinds. It's used to provide more context to the error.
// type Kind uint8

// func (k Kind) String() string {
// 	switch k {
// 	case KindUnsupported:
// 		return "unsupported"
// 	case KindInvalidArgs:
// 		return "invalid_arguments"
// 	case KindPreconditionViolation:
// 		return "precondition_violation"
// 	case KindOutOfRange:
// 		return "out of range"
// 	case KindUnauthenticated:
// 		return "unauthenticated"
// 	case KindPermissionDenied:
// 		return "permission_denied"
// 	case KindNotFound:
// 		return "not_found"
// 	case KindAborted:
// 		return "aborted"
// 	case KindAlreadyExists:
// 		return "already_exists"
// 	case KindResourceExhausted:
// 		return "resource_exhausted"
// 	case KindCanceled:
// 		return "cancelled"
// 	case KindDataLoss:
// 		return "data_loss"
// 	case KindUnknown:
// 		return "unknown"
// 	case KindInternal:
// 		return "internal"
// 	case KindUnimplemented:
// 		return "not_implemented"
// 	case KindUnavailable:
// 		return "unavailable"
// 	case KindDeadlineExceeded:
// 		return "deadline_exceeded"
// 	default:
// 		return "unsupported"
// 	}
// }

type Error struct {
	LogLevel      LogLevel
	status        status.Status
	detailsHidden bool
}

var errNotFound = errors.New("something was not found")

func (e *Error) findBadRequest() (*errdetails.BadRequest, error) {
	for _, detail := range e.status.Details() {
		switch v := detail.(type) {
		case *errdetails.BadRequest:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (e *Error) findDebugInfo() (*errdetails.DebugInfo, error) {
	for _, detail := range e.status.Details() {
		switch v := detail.(type) {
		case *errdetails.DebugInfo:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (e *Error) findPreconditionFailure() (*errdetails.PreconditionFailure, error) {
	for _, detail := range e.status.Details() {
		switch v := detail.(type) {
		case *errdetails.PreconditionFailure:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (e *Error) findErrorInfo() (*errdetails.ErrorInfo, error) {
	for _, detail := range e.status.Details() {
		switch v := detail.(type) {
		case *errdetails.ErrorInfo:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (e *Error) findQuotaFailure() (*errdetails.QuotaFailure, error) {
	for _, detail := range e.status.Details() {
		switch v := detail.(type) {
		case *errdetails.QuotaFailure:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

// AddBadRequestViolations adds a list of bad request violations to the error details. If the error details already
// contain bad request violations, the new ones are appended to the existing ones.
//
// It is recommended to include a bad request violation for the following error types:
//   - INVALID_ARGUMENT
//   - OUT_OF_RANGE
//
// See: https://cloud.google.com/apis/design/errors#error_payloads
func (e *Error) AddBadRequestViolations(violations []BadRequestViolation) {
	violationspb := make([]*errdetails.BadRequest_FieldViolation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.BadRequest_FieldViolation{Field: v.Field, Description: v.Description}
	}
	existing, err := e.findBadRequest()
	if errors.Is(err, errNotFound) {
		detail := errdetails.BadRequest{FieldViolations: violationspb}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
		return
	}
	existing.FieldViolations = append(existing.FieldViolations, violationspb...)
}

// AddPreconditionViolations adds a list of precondition violations to the error details. If the error details already
// contain precondition violations, the new ones are appended to the existing ones.
func (e *Error) AddPreconditionViolations(violations []PreconditionViolation) {
	violationspb := make([]*errdetails.PreconditionFailure_Violation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.PreconditionFailure_Violation{Description: v.Description, Subject: v.Subject, Type: v.Typ}
	}
	existing, err := e.findPreconditionFailure()
	if errors.Is(err, errNotFound) {
		detail := errdetails.PreconditionFailure{Violations: violationspb}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
		return
	}
	existing.Violations = append(existing.Violations, violationspb...)
}

// SetErrorInfo sets error info details to the error details. If the error details already contain error info
// details, they are overwritten. If the domain is empty, the domain is set to the package's default domain which
// should be set at startup-time by calling the Init() function. A non-empty domain is required.
//
// It is recommended to include an error info detail for the following error types:
//   - UNAUTHENTICATED
//   - PERMISSION_DENIED
//   - ABORTED
//
// See: https://cloud.google.com/apis/design/errors#error_payloads
func (e *Error) SetErrorInfo(domain, reason string, metadata map[string]any) {
	if reason == "" {
		return
	}
	if domain == "" {
		domain = maker.domain
	}
	metadatapb := make(map[string]string, len(metadata))
	for k, v := range metadata {
		metadatapb[k] = fmt.Sprintf("%v", v)
	}
	existing, err := e.findErrorInfo()
	if errors.Is(err, errNotFound) {
		detail := errdetails.ErrorInfo{Domain: domain, Reason: reason, Metadata: metadatapb}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
		return
	}
	existing.Domain = domain
	existing.Reason = reason
	existing.Metadata = metadatapb
}

// AddResourceInfos adds resource info details to the error details. If the error details already contain a
// resource info detail, it is overwritten.
//
// It is recommended to include a resource info detail for the following error types:
//   - NOT_FOUND
//   - ALREADY_EXISTS
//
// See: https://cloud.google.com/apis/design/errors#error_payloads
func (e *Error) AddResourceInfos(infos []ResourceInfo) {
	for _, info := range infos {
		detail := errdetails.ResourceInfo{
			Description:  info.Description,
			ResourceName: info.ResourceName,
			ResourceType: info.ResourceType,
			Owner:        info.Owner,
		}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
	}
}

// SetDebugInfoDetail sets debug info detail to the error details. If the error details already contain a debug info
// detail, it is overwritten. If the detail is empty, the operation is a no-op.
//
// It is recommended to include a debug info detail for the following error types:
//   - DATA_LOSS
//   - UNKNOWN
//   - INTERNAL
//   - UNAVAILABLE
//   - DEADLINE_EXCEEDED
//
// See: https://cloud.google.com/apis/design/errors#error_payloads
func (e *Error) SetDebugInfo(detail string, stackEntries []string) {
	if detail == "" {
		return
	}

	existing, err := e.findDebugInfo()
	if errors.Is(err, errNotFound) {
		detail := errdetails.DebugInfo{Detail: detail, StackEntries: stackEntries}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
		return
	}

	existing.Detail = detail
	existing.StackEntries = stackEntries
}

type QuotaViolation struct {
	Subject     string
	Description string
}

// AddQuotaViolations adds a list of quota violations to the error details. If the error details already contain quota
// violations, the new ones are appended to the existing ones.
func (e *Error) AddQuotaViolations(violations []QuotaViolation) {
	violationspb := make([]*errdetails.QuotaFailure_Violation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.QuotaFailure_Violation{Subject: v.Subject, Description: v.Description}
	}
	existing, err := e.findQuotaFailure()
	if errors.Is(err, errNotFound) {
		detail := errdetails.QuotaFailure{Violations: violationspb}
		status, err := e.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		e.status = *status
		return
	}
	existing.Violations = append(existing.Violations, violationspb...)
}

// // Var models what the circumstances were when the error was encountered and is used to provide additional context
// // to the error. Its purpose is to be logged and thereby give context to the error in the logs.
// type Var struct {
// 	Name  string
// 	Value any
// }

// // Root models a root error and is used when you first encounter an error in your code.
// // The idiomatic way to use this type is to create a new instance by using the "new" built-in function.
// //
// //		Ex.
// //	 	err := errors.NewRoot(err, "more context to err").AddVar("id", id).SetKind(errors.NotFound)
// //
// //nolint:errname
// type Root struct {
// 	// Err is the wrapped error.
// 	Err error
// 	// Kind is the kind of error that was encountered.
// 	Kind Kind
// 	// RuntimeState is a snapshot of the state of the application when the error was encountered.
// 	RuntimeState []Var
// 	// DirectlyRetryable is a flag that indicates whether the operation is directly retryable by the client.
// 	// Directly means that the exact same operation should be retried. For example, a failed gRPC or HTTP request
// 	// whose response indicates that the operation may be retried.
// 	DirectlyRetryable bool
// 	// IndirectlyRetryable is a flag that indicates whether the operation is indirectly retryable by the client.
// 	// Indirectly means that a retry should be attempted at a higher level in the system. For example, a database
// 	// transaction that failed because of a deadlock, or a version mismatch in a system that uses optimistic
// 	// concurrency. Such scenarios are often resolved by retrying the whole transaction-level operation, not just the
// 	// failing part. In contrast, directly retryable errors are often resolved by retrying the exact same operation.
// 	// Such retries are often only useful for long running background jobs.
// 	IndirectlyRetryable bool
// 	// Severity is the severity of the error used to control the way the error is logged. For example as an error,
// 	// warning, notice etc.
// 	Severity LogLevel
// }

// // NewInternal creates a new Root error. It's a convenience function that allows you to create a Root error with an original
// // error and a message (or message chain) that provides more context to the error.
// // If you don't want to provide a message, you can omit it. If you provide multiple messages,
// // they will be wrapped in a chain of WrappedError instances where the first message is the innermost message.
// //
// // Ex. NewInternal(err, "message1", "message2", "message3") will create a chain of WrappedError instances where "message1"
// // is the innermost message and "message3" is the outermost message. When you call Error() on the Root instance, you'll
// // get "message3: message2: message1: original error message".
// func NewInternal(original error, messages ...string) *Root {
// 	if len(messages) == 0 {
// 		return &Root{Err: original}
// 	}
// 	var root Root
// 	for i, msg := range messages {
// 		if i == 0 {
// 			root.Err = &WrappedError{Msg: msg, Err: original}
// 			continue
// 		}
// 		root.Err = &WrappedError{Msg: msg, Err: root.Err}
// 	}
// 	return &root
// }

// // Error returns the error message of the wrapped error that was encountered.
// //
// // If you have a Root instance whose 'Err' field value is a WrappedError, calling Error() on the Root instance will
// // return the message of the WrappedError instance.
// // Ex. Calling root.Error() will return "message1: original error message", given that the WrappedError instance
// // has the message "message1".
// func (r *Root) Error() string {
// 	if r.Err == nil {
// 		return ""
// 	}
// 	return r.Err.Error()
// }

// // Original returns the original error that was encountered. It unwraps the error chain and returns the original.
// func (r *Root) Original() error {
// 	err := r.Err
// 	for {
// 		unwrapped := errors.Unwrap(err)
// 		if unwrapped != nil {
// 			err = unwrapped
// 			continue
// 		}
// 		return err
// 	}
// }

// // SetError sets the Root's original error. It must only be called once.
// func (r *Root) SetError(err error) *Root {
// 	r.Err = err
// 	return r
// }

// // AddVar adds a variable to the RuntimeState slice. It should only be used on variables that will still be in scope
// // at log-time. If the variable won't be in scope at log-time, it must not be added to the RuntimeState slice.
// //
// // For out-of-scope variables, there are other methods that can be used to add them to the RuntimeState slice such as
// // 'WithIOReaderVar', which is provided for io.Reader variables.
// //
// // This method may be called multiple times.
// func (r *Root) AddVar(name string, value any) *Root {
// 	if name == "" || value == nil {
// 		return r
// 	}
// 	r.RuntimeState = append(r.RuntimeState, Var{Name: name, Value: value})
// 	return r
// }

// // AddVars allows you to add multiple variables to the RuntimeState slice. It's the bulk-variant of AddVar.
// func (r *Root) AddVars(vars ...Var) *Root {
// 	for _, v := range vars {
// 		_ = r.AddVar(v.Name, v.Value)
// 	}
// 	return r
// }

// // WithIOReaderVar reads the io.Reader and adds the contents to the RuntimeState slice.
// func (r *Root) WithIOReaderVar(name string, value io.Reader) *Root {
// 	if name == "" || value == nil {
// 		return r
// 	}
// 	b, err := io.ReadAll(value)
// 	if err != nil {
// 		return r
// 	}
// 	return r.AddVar(name, string(b))
// }

// // SetKind sets the Root's kind. It must only be called once.
// func (r *Root) SetKind(kind Kind) *Root {
// 	r.Kind = kind
// 	return r
// }

// // WithDirectRetry sets the Root's DirectlyRetryable flag to true.
// func (r *Root) WithDirectRetry() *Root {
// 	r.DirectlyRetryable = true
// 	return r
// }

// // Unwrap returns the original error that was encountered.
// func (r *Root) Unwrap() error {
// 	return r.Err
// }

// // WithSeverity sets the Root's severity.
// func (r *Root) WithSeverity(s LogLevel) *Root {
// 	r.Severity = s
// 	return r
// }

// // Is traverses err's error chain and compares the first encountered Root error for equality. If equal, it returns true.
// // If no such error is found or if they're not equal, it returns false.
// func (r *Root) Is(err error) bool {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		return root.Kind == r.Kind && root.Err.Error() == r.Err.Error()
// 	}
// 	return false
// }

// type WrappedError struct {
// 	Msg string
// 	Err error
// }

// func (e *WrappedError) Error() string {
// 	if e.Err == nil {
// 		return e.Msg
// 	}
// 	return e.Msg + ": " + e.Err.Error()
// }

// func (e *WrappedError) Unwrap() error {
// 	return e.Err
// }

// // Wrap is a utility function that makes it easier to wrap errors with a message to add more context.
// func Wrap(err error, msg string) error {
// 	if err == nil {
// 		return nil
// 	}
// 	if msg == "" {
// 		return err
// 	}
// 	return &WrappedError{Msg: msg, Err: err}
// }

// // Is traverses err's error chain and compares the first encountered Root error's Kind. If equal, it returns true.
// // If no such error is found or if their kinds differ, it returns false.
// func Is(err error, kind Kind) bool {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		return root.Kind == kind
// 	}
// 	return false
// }

// // AddVars adds variables to the RuntimeState slice of a Root error. It's a utility function that makes it easier to
// // add variables to the RuntimeState slice of a Root error.
// func AddVars(err error, vars ...Var) {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		root.RuntimeState = append(root.RuntimeState, vars...)
// 	}
// }

// // IsDirectlyRetryable returns true if the error is directly retryable, otherwise it returns false.
// func IsDirectlyRetryable(err error) bool {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		return root.DirectlyRetryable
// 	}
// 	return false
// }

// // IsIndirectlyRetryable returns true if the error is indirectly retryable, otherwise it returns false.
// func IsIndirectlyRetryable(err error) bool {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		return root.DirectlyRetryable
// 	}
// 	return false
// }

// // ClearDirectlyRetryable sets the error's DirectlyRetryable flag to false.
// func ClearDirectlyRetryable(err error) error {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		root.DirectlyRetryable = false
// 	}
// 	return err
// }

// // ClearIndirectlyRetryable sets the error's IndirectlyRetryable flag to false.
// func ClearIndirectlyRetryable(err error) error {
// 	var root *Root
// 	if errors.As(err, &root) {
// 		root.IndirectlyRetryable = false
// 	}
// 	return err
// }

// var Factory = factory{}

// type factory struct {
// 	domain          string
// 	requestIDCtxKey string
// }

// // SetFactory initializes the global factory. It must only be called at application startup-time. It is NOT thread-safe.
// func SetFactory(domain, requestIDCtxKey string) {
// 	Factory = factory{domain: domain, requestIDCtxKey: requestIDCtxKey}
// }

// func (_ factory) NewUnknownError(err error) *Error {
// 	const msg = "something unknown happened"
// 	root := *NewInternal(err, msg).SetKind(KindUnknown).WithSeverity(LogLevelError)
// 	e := &Error{
// 		Internal: root,
// 		API:      status.New(codes.Unknown, "something unknown happened"),
// 	}
// 	return e
// }

// func (_ factory) NewDeadlineExceeded(err error) *Error {
// 	const msg = "context deadline exceeded"
// 	root := *NewInternal(err, msg).SetKind(KindDeadlineExceeded).WithSeverity(LogLevelWarn)
// 	return &Error{
// 		Internal: root,
// 		API:      status.New(codes.DeadlineExceeded, msg),
// 	}
// }

// func (f factory) NewInternalError(err error) *Error {
// 	const msg = "something bad happened"
// 	if errors.Is(err, context.DeadlineExceeded) {
// 		return f.NewDeadlineExceeded(err)
// 	}
// 	return &Error{
// 		Internal: *NewInternal(err, msg).SetKind(KindInternal).WithSeverity(LogLevelError),
// 		API:      status.New(codes.Internal, msg),
// 	}
// }

// type BadRequestViolation struct {
// 	Field, Description string
// }

// func (f factory) NewInvalidArgumentError(violation BadRequestViolation) *Error {
// 	const msg = "one or more request arguments were invalid"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindInvalidArgs).
// 			AddVar("arg_field", violation.Field).AddVar("arg_desc", violation.Description),
// 		API: status.New(codes.InvalidArgument, msg),
// 	}
// 	if err := e.AddBadRequestViolation(violation.Field, violation.Description); err != nil {
// 		e := f.NewInternalError(fmt.Errorf("failed to add bad request violation: %w", err))
// 		e.Internal.AddVar("arg_field", violation.Field).AddVar("arg_desc", violation.Description)
// 		return e
// 	}
// 	return e
// }

// func (_ factory) NewInvalidArgumentErrors(violations []BadRequestViolation) *Error {
// 	const msg = "one or more request arguments were invalid"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindInvalidArgs).AddVar("arg_violations", violations),
// 		API:      status.New(codes.InvalidArgument, msg),
// 	}

// 	violationspb := make([]*errdetails.BadRequest_FieldViolation, 0, len(violations))
// 	for i, v := range violations {
// 		violationspb[i] = &errdetails.BadRequest_FieldViolation{Field: v.Field, Description: v.Description}
// 	}
// 	_ = e.AddBadRequestViolations(violationspb)

// 	return e
// }

// func (_ factory) NewPreconditionFailure(description, subject, typ string) *Error {
// 	const msg = "one or more request preconditions failed"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindPreconditionViolation).WithSeverity(LogLevelWarn),
// 		API:      status.New(codes.FailedPrecondition, msg),
// 	}
// 	_ = e.AddPreconditionViolations([]*errdetails.PreconditionFailure_Violation{{
// 		Description: description, Subject: subject, Type: typ,
// 	}})
// 	return e
// }

// func (_ factory) NewModelBindingError(err error) *Error {
// 	const msg = "failed to bind model"
// 	e := &Error{
// 		Internal: *NewInternal(err, msg).SetKind(KindInvalidArgs).WithSeverity(LogLevelInfo),
// 		API:      status.New(codes.InvalidArgument, msg),
// 	}
// 	return e
// }

// func (f factory) NewUnauthenticatedError(err error) *Error {
// 	const msg = "request could not be authenticated"
// 	e := &Error{
// 		Internal: *NewInternal(err, msg).SetKind(KindUnauthenticated).WithSeverity(LogLevelWarn),
// 		API:      status.New(codes.Unauthenticated, msg),
// 	}
// 	// TODO: Make reason and metadata dynamic
// 	_ = e.SetErrorInfo(f.domain, "UNAUTHENTICATED", map[string]string{
// 		"err": err.Error(),
// 	})
// 	return e
// }

// func (_ factory) NewNotFoundError(desc, rscName, rscType string) *Error {
// 	const msg = "requested resource not found"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindNotFound).WithSeverity(LogLevelInfo),
// 		API:      status.New(codes.NotFound, msg),
// 	}
// 	_ = e.SetResourceInfo(desc, rscName, rscType)
// 	return e
// }

// func (_ factory) NewNotFoundErrors(infos []*errdetails.ResourceInfo) *Error {
// 	const msg = "requested resources not found"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindNotFound).WithSeverity(LogLevelInfo),
// 		API:      status.New(codes.NotFound, msg),
// 	}
// 	_ = e.SetResourceInfos(infos)
// 	return e
// }

// func (f factory) NewPermissionDeniedError(reason string, metadata map[string]string) *Error { // nolint:unparam
// 	const msg = "permission denied"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindPermissionDenied).WithSeverity(LogLevelWarn),
// 		API:      status.New(codes.PermissionDenied, msg),
// 	}
// 	_ = e.SetErrorInfo(f.domain, "PERMISSION_DENIED", metadata)
// 	return e
// }

// func (_ factory) NewCanceledError() *Error {
// 	const msg = "request canceled by the client"
// 	e := &Error{
// 		Internal: *NewInternal(nil, msg).SetKind(KindCanceled),
// 		API:      status.New(codes.Canceled, msg),
// 	}
// 	return e
// }

// var errNotFound = errors.New("something was not found")

// type Error struct {
// 	Internal        Root
// 	API             *status.Status
// 	requestIDCtxKey string
// }

// func (e *Error) findBadRequest() (*errdetails.BadRequest, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.BadRequest:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// func (e *Error) findResourceInfo() (*errdetails.ResourceInfo, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.ResourceInfo:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// func (e *Error) findDebugInfo() (*errdetails.DebugInfo, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.DebugInfo:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// func (e *Error) findRequestInfo() (*errdetails.RequestInfo, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.RequestInfo:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// func (e *Error) findPreconditionFailure() (*errdetails.PreconditionFailure, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.PreconditionFailure:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// func (e *Error) findErrorInfo() (*errdetails.ErrorInfo, error) {
// 	for _, detail := range e.API.Details() {
// 		switch v := detail.(type) {
// 		case *errdetails.ErrorInfo:
// 			return v, nil
// 		default:
// 			continue
// 		}
// 	}
// 	return nil, errNotFound
// }

// // SetRequestInfoDetails sets the request ID in the error details. If the error details already contain a request ID,
// // it is overwritten.
// func (e *Error) SetRequestInfo(ctx context.Context) error {
// 	reqID, exists := ctx.Value(e.requestIDCtxKey).(string)
// 	if !exists {
// 		return nil
// 	}

// 	existing, err := e.findRequestInfo()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.RequestInfo{RequestId: reqID}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	existing.RequestId = reqID
// 	return nil
// }

// // AddBadRequestViolation adds a bad request violation to the error details. If the error details already contain a bad
// // request detail, the new field violation is appended to the existing ones.
// func (e *Error) AddBadRequestViolation(field, desc string) error {
// 	existing, err := e.findBadRequest()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.BadRequest{
// 			FieldViolations: []*errdetails.BadRequest_FieldViolation{{Description: desc, Field: field}},
// 		}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	existing.FieldViolations = append(
// 		existing.FieldViolations, &errdetails.BadRequest_FieldViolation{Description: desc, Field: field},
// 	)
// 	return nil
// }

// // AddBadRequestViolations adds a list of bad request violations to the error details. If the error details already
// // contain bad request violations, the new ones are appended to the existing ones.
// func (e *Error) AddBadRequestViolations(violations []*errdetails.BadRequest_FieldViolation) []error {
// 	var errs []error
// 	for _, v := range violations {
// 		if err := e.AddBadRequestViolation(v.Field, v.Description); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
// 	return errs
// }

// // AddPreconditionViolations adds a list of precondition violations to the error details. If the error details already
// // contain precondition violations, the new ones are appended to the existing ones.
// func (e *Error) AddPreconditionViolations(violations []*errdetails.PreconditionFailure_Violation) error {
// 	existing, err := e.findPreconditionFailure()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.PreconditionFailure{Violations: violations}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	existing.Violations = append(existing.Violations, violations...)
// 	return nil
// }

// // SetErrorInfo sets error info details to the error details. If the error details already contain error info
// // details, they are overwritten.
// func (e *Error) SetErrorInfo(domain, reason string, metadata map[string]string) error {
// 	existing, err := e.findErrorInfo()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.ErrorInfo{Domain: domain, Reason: reason, Metadata: metadata}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// 	existing.Domain = domain
// 	existing.Reason = reason
// 	existing.Metadata = metadata
// 	return nil
// }

// // SetResourceInfoDetail sets a resource info detail to the error details. If the error details already contain a
// // resource info detail, it is overwritten.
// func (e *Error) SetResourceInfo(desc, resourceName, resourceType string) error {
// 	existing, err := e.findResourceInfo()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.ResourceInfo{
// 			Description:  desc,
// 			ResourceName: resourceName,
// 			ResourceType: resourceType,
// 		}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	// Overwrite the existing resource info detail
// 	existing.ResourceName = resourceName
// 	existing.ResourceType = resourceType
// 	existing.Description = desc
// 	return nil
// }

// // SetResourceInfoDetail sets a resource info detail to the error details. If the error details already contain a
// // resource info detail, it is overwritten.
// func (e *Error) SetResourceInfos(infos []*errdetails.ResourceInfo) error {
// 	for _, info := range infos {
// 		if err := e.SetResourceInfo(info.Description, info.ResourceName, info.ResourceType); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // SetDebugInfoDetail sets debug info detail to the error details. If the error details already contain a debug info
// // detail, it is overwritten.
// func (e *Error) SetDebugInfo(detail string, stackEntries []string) error {
// 	if detail == "" {
// 		return nil
// 	}

// 	existing, err := e.findDebugInfo()
// 	if errors.Is(err, errNotFound) {
// 		detail := errdetails.DebugInfo{Detail: detail, StackEntries: stackEntries}
// 		e.API, err = e.API.WithDetails(&detail)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}

// 	existing.Detail = detail
// 	existing.StackEntries = stackEntries
// 	return nil
// }
