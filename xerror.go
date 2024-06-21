/*
Package xerror provides a way to wrap errors with additional context and to add variables to the error that can be
logged at a later time. It also provides a way to categorize errors into different kinds.
*/
package xerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

/*
TODO:
	1. Make it easy to consume gRPC errors
	2. Make it easy to consume HTTP errors by generating typescript models from status.Status
	5. Make it easy to produce gRPC errors
	6. Make it easy to respond with gRPC errors (unaryinterceptor)
*/

var errNotFound = errors.New("something was not found")

// Optional is a type used to model a value that may or may not be present. It is used to model the return
// value of a function that may return a value or not. It is used to avoid returning nil values.
type Optional[T any] struct {
	// Value is the value that may or may not be present.
	Value T
	// Valid is true if the value is present, otherwise it is false.
	Valid bool
}

func newValidOptional[T any](value T) Optional[T] {
	return Optional[T]{Value: value, Valid: true}
}

func newInvalidOptional[T any]() Optional[T] {
	return Optional[T]{Valid: false}
}

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

func (xerr *Error) findResourceInfos() ([]*errdetails.ResourceInfo, error) {
	var infos []*errdetails.ResourceInfo
	for _, detail := range xerr.status.Details() {
		switch v := detail.(type) {
		case *errdetails.ResourceInfo:
			infos = append(infos, v)
		default:
			continue
		}
	}
	if len(infos) == 0 {
		return nil, errNotFound
	}
	return infos, nil
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
// should be set at startup-time by calling the Init() function.
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
// It is NOT recommended to include a debug info detail since it'll be returned to the caller. If you need to log
// debug info, use the runtime state instead (the AddVar() and AddVars() methods). It's is however possible to include
// a debug info detail and still not return it to the caller by calling the HideDetails() method.
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

// QuotaViolation is a message type used to describe a single quota violation.  For example, a
// daily quota or a custom quota that was exceeded.
type QuotaViolation struct {
	//Subject on which the quota check failed.
	// For example, "clientip:<ip address of client>" or "project:<Google
	// developer project id>".
	Subject string
	// Descriptions is a description of how the quota check failed. Clients can use this
	// description to find more about the quota configuration in the service's
	// public documentation, or find the relevant quota limit to adjust through
	// developer console.
	//
	// For example: "Service disabled" or "Daily Limit for read operations
	// exceeded".
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

// BadRequestViolations returns a list of bad request violations. If the error details do not contain bad request
// violations, it returns nil.
func (xerr *Error) BadRequestViolations() []BadRequestViolation {
	pb, err := xerr.findBadRequest()
	if errors.Is(err, errNotFound) {
		return nil
	}
	violations := make([]BadRequestViolation, len(pb.FieldViolations))
	for i, v := range pb.FieldViolations {
		violations[i] = BadRequestViolation{Field: v.Field, Description: v.Description}
	}
	return violations
}

// PreconditionsViolations returns a list of precondition violations. If the error details do not contain precondition
// violations, it returns nil.
func (xerr *Error) PreconditionViolations() []PreconditionViolation {
	pb, err := xerr.findPreconditionFailure()
	if errors.Is(err, errNotFound) {
		return nil
	}
	violations := make([]PreconditionViolation, len(pb.Violations))
	for i, v := range pb.Violations {
		violations[i] = PreconditionViolation{Description: v.Description, Subject: v.Subject, Typ: v.Type}
	}
	return violations
}

// ErrorInfo describes the cause of the error with structured details.
//
// Example of an error when contacting the "pubsub.googleapis.com" API when it
// is not enabled:
//
//	{ "reason": "API_DISABLED",
//	  "domain": "googleapis.com",
//	  "metadata": {
//	    "resource": "projects/123",
//	    "service": "pubsub.googleapis.com"
//	  }
//	}
//
// This response indicates that the pubsub.googleapis.com API is not enabled.
// Example of an error that is returned when attempting to create a Spanner
// instance in a region that is out of stock:
//
//	{ "reason": "STOCKOUT",
//	  "domain": "spanner.googleapis.com",
//	  "metadata": {
//	    "availableRegions": "us-central1,us-east2"
//	  }
//	}
type ErrorInfo struct {
	// Domain is the logical grouping to which the "reason" belongs. The error domain is typically the registered
	// service name of the tool or product that generated the error. The domain must be a globally unique value.
	Domain string
	// Reason is a short snake_case description of why the error occurred. Error reasons are unique within a particular
	// domain of errors. The application should define an enum of error reasons.
	//
	// The reason should have these properties:
	//	- Be meaningful enough for a human reader to understand what the reason refers to.
	//	- Be unique and consumable by machine actors for automation.
	//	- Example: CPU_AVAILABILITY
	//	- Distill your error message into its simplest form. For example, the reason string could be one of the
	//	  following text examples in UPPER_SNAKE_CASE: UNAVAILABLE, NO_STOCK, CHECKED_OUT, AVAILABILITY_ERROR, if your
	//	  error message is:
	//	  The Book, "The Great Gatsby", is unavailable at the Library, "Garfield East". It is expected to be available
	//	  again on 2199-05-13.
	Reason string
	// Metadata is additional structured details about this error, which should provide important context for clients
	// to identify resolution steps. Keys should be in lower camel-case, and be limited to 64 characters in length.
	// When identifying the current value of an exceeded limit, the units should be contained in the key, not the value.
	//
	// Example: {"vmType": "e2-medium", "attachment": "local-ssd=3,nvidia-t4=2", "zone": us-east1-a"}
	Metadata map[string]string
}

// ErrorInfo returns the error info details. If the error details do not contain error info details, it returns an
// invalid optional.
func (xerr *Error) ErrorInfo() Optional[ErrorInfo] {
	pb, err := xerr.findErrorInfo()
	if errors.Is(err, errNotFound) {
		return newInvalidOptional[ErrorInfo]()
	}
	return newValidOptional(ErrorInfo{Domain: pb.Domain, Reason: pb.Reason, Metadata: pb.Metadata})
}

// Describes additional debugging info.
type DebugInfo struct {
	// Additional debugging information provided by the server.
	Detail string
	// The stack trace entries indicating where the error occurred.
	StackEntries []string
}

// DebugInfo returns the error info details. If the error details do not contain error info details, it returns an
// invalid optional.
func (xerr *Error) DebugInfo() Optional[DebugInfo] {
	pb, err := xerr.findDebugInfo()
	if errors.Is(err, errNotFound) {
		return newInvalidOptional[DebugInfo]()
	}
	return newValidOptional(DebugInfo{Detail: pb.Detail, StackEntries: pb.StackEntries})
}

// ResourceInfos returns a list of resource info details. If the error details do not contain resource info details, it
// returns nil.
func (xerr *Error) ResourceInfos() []ResourceInfo {
	pb, err := xerr.findResourceInfos()
	if errors.Is(err, errNotFound) {
		return nil
	}
	infos := make([]ResourceInfo, len(pb))
	for i, v := range pb {
		infos[i] = ResourceInfo{Description: v.Description,
			ResourceName: v.ResourceName,
			ResourceType: v.ResourceType,
			Owner:        v.Owner,
		}
	}
	return infos
}

// QuotaViolations returns a list of quota violations. If the error details do not contain quota violations, it returns
// nil.
func (xerr *Error) QuotaViolations() []QuotaViolation {
	pb, err := xerr.findQuotaFailure()
	if errors.Is(err, errNotFound) {
		return nil
	}
	violations := make([]QuotaViolation, len(pb.Violations))
	for i, v := range pb.Violations {
		violations[i] = QuotaViolation{Subject: v.Subject, Description: v.Description}
	}
	return violations
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

// ShowDetails marks the error as having shown details. This is the inverse of HideDetails.
func (xerr *Error) ShowDetails() *Error {
	xerr.detailsHidden = false
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

// StatusProto returns the status proto contained in the error.
func (xerr *Error) StatusProto() *spb.Status {
	return xerr.status.Proto()
}

func (xerr *Error) StatusCode() codes.Code {
	return xerr.status.Code()
}

func (xerr *Error) StatusMessage() string {
	return xerr.status.Message()
}

// IsDomainError compares the error with the provided domain-specific error details (the domain and reason).
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
//		 xerr := xgrpc.ErrorFrom(err)
//		 if xerr.IsDomainError(othersystem.Domain, othersystem.NO_STOCK) {
//			 requestMoreStock() // decision based on the error type
//		 }
//	 }
func (xerr *Error) IsDomainError(domain, reason string) bool {
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

// MarshalJSON marshals the error to JSON. This is only useful for testing purposes to be able to generate golden
// files to be able to inspect the error in a human-readable format.
func (xerr *Error) MarshalJSON() ([]byte, error) {
	type marshallable struct {
		LogLevel      LogLevel    `json:"logLevel"`
		Status        *spb.Status `json:"status"`
		DetailsHidden bool        `json:"detailsHidden"`
		RuntimeState  []Var       `json:"runtimeState"`
	}
	err := marshallable{
		LogLevel:      xerr.logLevel,
		Status:        xerr.status.Proto(),
		DetailsHidden: xerr.detailsHidden,
		RuntimeState:  xerr.runtimeState,
	}
	return json.Marshal(err)
}

// WrappedError is a model that makes it easy to add more context to an error as it is passed up the call stack.
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

// From returns an Error instance from an error. It's meant to be used in your application, at the place in the code
// where the error is logged.
//
// If the error is not an Error instance, then it is an unexpected error and should be logged, so it can be discovered
// that there's code where the error isn't correctly handled.
func From(err error) *Error {
	var xerr *Error
	if !errors.As(err, &xerr) {
		return &Error{
			logLevel: LogLevelError,
			status:   *status.New(codes.Unknown, err.Error()),
		}
	}
	return xerr
}
