/*
Package xerror provides a way to wrap errors with additional context and to add variables to the error that can be
logged at a later time. It also provides a way to categorize errors into different kinds.
*/
package xerror

import (
	"errors"
	"fmt"
	"slices"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
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

var errNotFound = errors.New("something was not found")

// Var models what the circumstances were when the error was encountered and is used to provide additional context
// to the error. Its purpose is to be logged and thereby give context to the error in the logs.
type Var struct {
	Name  string
	Value any
}

// LogLevel is used to control the way the error is logged. For example as an error, warning, notice etc.
type LogLevel uint8

const (
	LogLevelUnspecified LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type Error struct {
	logLevel      LogLevel
	status        status.Status
	detailsHidden bool
	// runtimeState is a snapshot of the state of the application when the error was encountered. It is used to provide
	// additional context to the error and is used to log the circumstances when the error was encountered.
	runtimeState []Var
}

func (xerr *Error) Error() string {
	return xerr.status.String()
}

func (xerr *Error) findBadRequest() (*errdetails.BadRequest, error) {
	for _, detail := range xerr.status.Details() {
		switch v := detail.(type) {
		case *errdetails.BadRequest:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (xerr *Error) findDebugInfo() (*errdetails.DebugInfo, error) {
	for _, detail := range xerr.status.Details() {
		switch v := detail.(type) {
		case *errdetails.DebugInfo:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (xerr *Error) findPreconditionFailure() (*errdetails.PreconditionFailure, error) {
	for _, detail := range xerr.status.Details() {
		switch v := detail.(type) {
		case *errdetails.PreconditionFailure:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (xerr *Error) findErrorInfo() (*errdetails.ErrorInfo, error) {
	for _, detail := range xerr.status.Details() {
		switch v := detail.(type) {
		case *errdetails.ErrorInfo:
			return v, nil
		default:
			continue
		}
	}
	return nil, errNotFound
}

func (xerr *Error) findQuotaFailure() (*errdetails.QuotaFailure, error) {
	for _, detail := range xerr.status.Details() {
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
func (xerr *Error) AddBadRequestViolations(violations []BadRequestViolation) *Error {
	violationspb := make([]*errdetails.BadRequest_FieldViolation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.BadRequest_FieldViolation{Field: v.Field, Description: v.Description}
	}
	existing, err := xerr.findBadRequest()
	if errors.Is(err, errNotFound) {
		detail := errdetails.BadRequest{FieldViolations: violationspb}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
		return xerr
	}
	existing.FieldViolations = append(existing.FieldViolations, violationspb...)
	return xerr
}

// AddPreconditionViolations adds a list of precondition violations to the error details. If the error details already
// contain precondition violations, the new ones are appended to the existing ones.
func (xerr *Error) AddPreconditionViolations(violations []PreconditionViolation) *Error {
	violationspb := make([]*errdetails.PreconditionFailure_Violation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.PreconditionFailure_Violation{Description: v.Description, Subject: v.Subject, Type: v.Typ}
	}
	existing, err := xerr.findPreconditionFailure()
	if errors.Is(err, errNotFound) {
		detail := errdetails.PreconditionFailure{Violations: violationspb}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
		return xerr
	}
	existing.Violations = append(existing.Violations, violationspb...)
	return xerr
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
func (xerr *Error) SetErrorInfo(domain, reason string, metadata map[string]any) *Error {
	if reason == "" {
		return xerr
	}
	if domain == "" {
		domain = maker.domain
	}
	metadatapb := make(map[string]string, len(metadata))
	for k, v := range metadata {
		metadatapb[k] = fmt.Sprintf("%v", v)
	}
	existing, err := xerr.findErrorInfo()
	if errors.Is(err, errNotFound) {
		detail := errdetails.ErrorInfo{Domain: domain, Reason: reason, Metadata: metadatapb}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
		return xerr
	}
	existing.Domain = domain
	existing.Reason = reason
	existing.Metadata = metadatapb
	return xerr
}

// AddResourceInfos adds resource info details to the error details. If the error details already contain a
// resource info detail, it is overwritten.
//
// It is recommended to include a resource info detail for the following error types:
//   - NOT_FOUND
//   - ALREADY_EXISTS
//
// See: https://cloud.google.com/apis/design/errors#error_payloads
func (xerr *Error) AddResourceInfos(infos []ResourceInfo) *Error {
	for _, info := range infos {
		detail := errdetails.ResourceInfo{
			Description:  info.Description,
			ResourceName: info.ResourceName,
			ResourceType: info.ResourceType,
			Owner:        info.Owner,
		}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
	}
	return xerr
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
func (xerr *Error) SetDebugInfo(detail string, stackEntries []string) *Error {
	if detail == "" {
		return xerr
	}

	existing, err := xerr.findDebugInfo()
	if errors.Is(err, errNotFound) {
		detail := errdetails.DebugInfo{Detail: detail, StackEntries: stackEntries}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
		return xerr
	}

	existing.Detail = detail
	existing.StackEntries = stackEntries
	return xerr
}

type QuotaViolation struct {
	Subject     string
	Description string
}

// AddQuotaViolations adds a list of quota violations to the error details. If the error details already contain quota
// violations, the new ones are appended to the existing ones.
func (xerr *Error) AddQuotaViolations(violations []QuotaViolation) *Error {
	violationspb := make([]*errdetails.QuotaFailure_Violation, len(violations))
	for i, v := range violations {
		violationspb[i] = &errdetails.QuotaFailure_Violation{Subject: v.Subject, Description: v.Description}
	}
	existing, err := xerr.findQuotaFailure()
	if errors.Is(err, errNotFound) {
		detail := errdetails.QuotaFailure{Violations: violationspb}
		status, err := xerr.status.WithDetails(&detail)
		if err != nil {
			panic(fmt.Errorf("%v: %w", err, ErrFailedToAddErrorDetails))
		}
		xerr.status = *status
		return xerr
	}
	existing.Violations = append(existing.Violations, violationspb...)
	return xerr
}

// AddVar adds a variable to the runtime state.
func (xerr *Error) AddVar(name string, value any) *Error {
	if name == "" || value == nil {
		return xerr
	}
	xerr.runtimeState = append(xerr.runtimeState, Var{Name: name, Value: value})
	return xerr
}

// AddVars adds multiple variables to the runtime state.
func (xerr *Error) AddVars(vars ...Var) *Error {
	for _, v := range vars {
		_ = xerr.AddVar(v.Name, v.Value)
	}
	return xerr
}

// RuntimeState returns the runtime state of the error. This is used when you want to log the circumstances when the
// error was encountered.
func (xerr *Error) RuntimeState() []Var {
	return xerr.runtimeState
}

// HideDetails marks the error as having hidden details. This is useful when you want to hide the details of the error
// from external callers. Effectively, this means that the "debug info" and "error info" details are removed from the
// error when returned to the caller. For this to work, the server has to use the implementation-specific functionality
// such as the unary interceptor for gRPC.
func (xerr *Error) HideDetails() *Error {
	xerr.detailsHidden = true
	return xerr
}

// LogLevel returns the log level of the error.
func (xerr *Error) LogLevel() LogLevel {
	return xerr.logLevel
}

// SetLogLevel sets the log level of the error.
func (xerr *Error) SetLogLevel(level LogLevel) *Error {
	xerr.logLevel = level
	return xerr
}

// IsDirectlyRetryable returns true if the call that caused the error is directly retryable, otherwise it returns false.
// Retries should be attempted using an exponential backoff strategy.
func (xerr *Error) IsDirectlyRetryable() bool {
	if xerr == nil {
		return false
	}
	if xerr.status.Code() == codes.Unavailable {
		return true
	}
	return false
}

// IsRetryableAtHigherLevel returns true if the call that caused the error cannot be directly retried, but instead
// should be retried at a higher level in the system. Example: an optimistic concurrency error in a database
// transaction, where the whole transaction should be retried, not just the failing part.
func (xerr *Error) IsRetryableAtHigherLevel() bool {
	if xerr == nil {
		return false
	}
	switch xerr.status.Code() {
	case codes.ResourceExhausted, codes.Aborted:
		return true
	default:
		return false
	}
}

// IsDetailsHidden returns true if the error details are hidden, otherwise it returns false.
func (xerr *Error) IsDetailsHidden() bool {
	return xerr.detailsHidden
}

// RemoveSensitiveDetails removes sensitive details from the error. This is useful when you want to return the error
// to the client, but you don't want to expose sensitive details such as debug info or error info.
func (xerr *Error) RemoveSensitiveDetails() *Error {
	// Find the indexes of the details that should be deleted
	var deletingDetails []int
	for i, detail := range xerr.status.Details() {
		switch detail.(type) {
		case *errdetails.DebugInfo, *errdetails.ErrorInfo:
			deletingDetails = append(deletingDetails, i)
		default:
			continue
		}
	}

	// Initialize a slice to hold the remaining details
	remainingDetails := make([]protoiface.MessageV1, 0, len(xerr.status.Details())-len(deletingDetails))
	for i, detail := range xerr.status.Details() {
		if slices.Contains(deletingDetails, i) {
			continue // skip the details that should be deleted
		}
		if d, ok := detail.(protoiface.MessageV1); ok {
			remainingDetails = append(remainingDetails, d)
		}
	}

	// Create a new status with the remaining detail that replaces the old status
	newStatus := status.New(xerr.status.Code(), xerr.status.Message())
	for _, detail := range remainingDetails {
		newStatus, _ = newStatus.WithDetails(detail)
	}
	xerr.status = *newStatus

	return xerr
}

// SetStatus sets the status of the error.
func (xerr *Error) SetStatus(s *status.Status) *Error {
	xerr.status = *s
	return xerr
}

// Status returns a copy of the status contained in the error.
func (xerr *Error) Status() *status.Status {
	return status.FromProto(xerr.status.Proto())
}

// EqualsDomainError compares the error with the provided domain-specific error details (the domain and reason).
// The reason is machine-readable and most importantly, it is unique within a particular domain of errors. This
// method is used to check if a returned error is a particular domain-specific error. This is useful when decisions
// need to be made based on the error type.
//
// Note! Sometimes, decisions can be made based on the status code alone, but when that is not granular enough,
// the domain and reason should be used to make decisions.
//
// Ex.
//
//	 err := othersystempb.SomeMethod(ctx, req)
//	 if err != nil {
//		 xerr := grpc.XErrorFrom(err)
//		 if xerr.EqualsDomainError(othersystemerror.Domain, othersystemerror.NO_STOCK) {
//			 requestMoreStock() // decision based on the error type
//		 }
//	 }
func (xerr *Error) EqualsDomainError(domain, reason string) bool {
	info, err := xerr.findErrorInfo()
	if errors.Is(err, errNotFound) {
		return false
	}
	return domain == info.Domain && reason == info.Reason
}

// DomainType returns a unique error type based on the domain and reason. This is used to enable switch-case statements.
//
// Ex.
//
//	 err := othersystempb.SomeMethod(ctx, req)
//	 if err != nil {
//		 xerr := grpc.XErrorFrom(err)
//		 switch xerr.DomainType() {
//		 case xerror.DomainType(othersystemerror.Domain, othersystemerror.NO_STOCK):
//			 requestMoreStock() // decision based on the error type
func (xerr *Error) DomainType() string {
	info, err := xerr.findErrorInfo()
	if errors.Is(err, errNotFound) {
		return ""
	}
	return DomainType(info.Domain, info.Reason)
}

type WrappedError struct {
	Msg string
	Err error
}

func (wr *WrappedError) Error() string {
	if wr.Err == nil {
		return wr.Msg
	}
	return wr.Msg + ": " + wr.Err.Error()
}

// Unwrap returns the directly wrapped error.
func (wr *WrappedError) Unwrap() error {
	return wr.Err
}

// AddVar adds runtime state information to the wrapped Error instance, if there is one.
func (wr *WrappedError) AddVar(name string, value any) *WrappedError {
	if name == "" || value == nil {
		return wr
	}
	var xerr *Error
	if !errors.As(wr.Err, &wr) {
		return wr
	}
	_ = xerr.AddVar(name, value)
	return wr
}

// AddVars adds multiple runtime state information to the wrapped Error instance, if there is one.
func (wr *WrappedError) AddVars(vars ...Var) *WrappedError {
	if len(vars) == 0 {
		return wr
	}
	var xerr *Error
	if !errors.As(wr.Err, &wr) {
		return wr
	}
	_ = xerr.AddVars(vars...)
	return wr
}

// XError returns the Error instance that is wrapped by the WrappedError instance, if there is one. Otherwise, it
// returns nil.
func (wr *WrappedError) XError() *Error {
	var xerr *Error
	if !errors.As(wr.Err, &xerr) {
		return nil
	}
	return xerr
}

// Wrap wrap errors with a message to add more context to the error. It is used when receiving an error from a
// call that is already an Error instance and you want to add more context to the error.
//
// Ex.
//
//	 err := pkg.Func() // returns an Error instance
//	 if err != nil {
//		  return errors.Wrap(err, "more context to err")
//	 }
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	if msg == "" {
		return err
	}
	return &WrappedError{Msg: msg, Err: err}
}

// DomainType returns a unique error type based on the domain and reason. This is used to enable switch-case statements.
func DomainType(domain, reason string) string {
	return domain + reason
}
