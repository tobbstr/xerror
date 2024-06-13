package xerror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFailedToAddErrorDetails = errors.New("failed to add error details")

var maker = factory{}

type factory struct {
	domain string
}

// Init initializes the package. It must called once, before creating any errors, and only be called at application
// startup-time. It is NOT thread-safe.
//
// The domain is the logical grouping to which the "reason" belongs. See the Reason field in the unexported errorInfo
// struct for more information about the "reason". The error domain is typically the registered service
// name of the tool or product that generated the error. The domain must be a globally unique value.
//   - Example: pubsub.googleapis.com
func Init(domain string) {
	maker = factory{domain: domain}
}

type BadRequestViolation struct {
	Field, Description string
}

type BadRequestViolationOptions struct {
	Violation BadRequestViolation
	LogLevel  LogLevel
}

func (f factory) newInvalidArgument(opts BadRequestViolationOptions) *Error {
	const msg = "one or more request arguments were invalid"
	return f.newBadRequest(msg, opts)
}

type BadRequestViolationsOptions struct {
	Violations []BadRequestViolation
	LogLevel   LogLevel
}

func (f factory) newInvalidArgumentErrors(opts BadRequestViolationsOptions) *Error {
	const msg = "one or more request arguments were invalid"
	return f.newBadRequestBulk(msg, opts)
}

type PreconditionViolation struct {
	Description, Subject, Typ string
}

type PreconditionFailureOptions struct {
	Violation PreconditionViolation
	LogLevel  LogLevel
}

func (f factory) newPreconditionFailure(opts PreconditionFailureOptions) *Error {
	const msg = "one or more request preconditions failed"
	e := &Error{
		status: *status.New(codes.FailedPrecondition, msg),
	}
	e.AddPreconditionViolations([]PreconditionViolation{opts.Violation})
	return e
}

type PreconditionFailuresOptions struct {
	Violations []PreconditionViolation
	LogLevel   LogLevel
}

func (f factory) newPreconditionFailures(opts PreconditionFailuresOptions) *Error {
	const msg = "one or more request preconditions failed"
	e := &Error{
		status: *status.New(codes.FailedPrecondition, msg),
	}

	e.AddPreconditionViolations(opts.Violations)
	return e
}

func (f factory) newOutOfRangeError(opts BadRequestViolationOptions) *Error {
	const msg = "one or more request arguments were out of range"
	return f.newBadRequest(msg, opts)
}

func (f factory) newOutOfRangeErrors(opts BadRequestViolationsOptions) *Error {
	const msg = "one or more request arguments were out of range"
	return f.newBadRequestBulk(msg, opts)
}

type ErrorInfoOptions struct {
	// Error is the error that occurred.
	Error error
	// Loglevel is the level of logging that should be used for this error.
	LogLevel LogLevel
	// Reason is a short snake_case description of why the error occurred. Error reasons are unique within a particular
	// domain of errors. The application should define an enum of error reasons.
	//
	// The reason should have these properties:
	//  - Be meaningful enough for a human reader to understand what the reason refers to.
	//  - Be unique and consumable by machine actors for automation.
	//  - Example: CPU_AVAILABILITY
	//  - Distill your error message into its simplest form. For example, the reason string could be one of the
	//    following text examples in UPPER_SNAKE_CASE: UNAVAILABLE, NO_STOCK, CHECKED_OUT, AVAILABILITY_ERROR, if your
	//    error message is:
	//    The Book, "The Great Gatsby", is unavailable at the Library, "Garfield East". It is expected to be available
	//    again on 2199-05-13.
	Reason string
	// Metadata is additional structured details about this error, which should provide important context for clients
	// to identify resolution steps. Keys should be in lower camel-case, and be limited to 64 characters in length.
	// When identifying the current value of an exceeded limit, the units should be contained in the key, not the value.
	//
	// Example: {"vmType": "e2-medium", "attachment": "local-ssd=3,nvidia-t4=2", "zone": us-east1-a"}
	Metadata map[string]any
}

func (f factory) newUnauthenticatedError(opts ErrorInfoOptions) *Error {
	return f.newErrorInfoError(codes.Unauthenticated, opts)
}

func (f factory) newPermissionDeniedError(opts ErrorInfoOptions) *Error { // nolint:unparam
	e := f.newErrorInfoError(codes.PermissionDenied, opts)
	return e
}

type ResourceInfo struct {
	Description, ResourceName, ResourceType, Owner string
}

type NotFoundOptions struct {
	ResourceInfo ResourceInfo
	LogLevel     LogLevel
}

func (_ factory) newNotFound(opts NotFoundOptions) *Error {
	const msg = "requested resources not found"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: opts.LogLevel,
	}

	e.AddResourceInfos([]ResourceInfo{opts.ResourceInfo})
	return e
}

type NotFoundBulkOptions struct {
	ResourceInfos []ResourceInfo
	LogLevel      LogLevel
}

func (_ factory) newNotFoundBulk(opts NotFoundBulkOptions) *Error {
	const msg = "requested resources not found"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: opts.LogLevel,
	}

	e.AddResourceInfos(opts.ResourceInfos)
	return e
}

func (f factory) newAborted(opts ErrorInfoOptions) *Error {
	return f.newErrorInfoError(codes.Aborted, opts)
}

type AlreadyExistsOptions struct {
	ResourceInfo ResourceInfo
	LogLevel     LogLevel
}

func (f factory) newAlreadyExists(opts AlreadyExistsOptions) *Error {
	const msg = "resource already exists"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: opts.LogLevel,
	}

	e.AddResourceInfos([]ResourceInfo{opts.ResourceInfo})
	return e
}

type AlreadyExistsBulkOptions struct {
	ResourceInfos []ResourceInfo
	LogLevel      LogLevel
}

func (_ factory) newAlreadyExistsBulk(opts AlreadyExistsBulkOptions) *Error {
	const msg = "resources already exist"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: opts.LogLevel,
	}
	e.AddResourceInfos(opts.ResourceInfos)
	return e
}

type ResourceExhaustedOptions struct {
	Error          error
	QuotaViolation QuotaViolation
	LogLevel       LogLevel
}

func (_ factory) newResourceExhausted(opts ResourceExhaustedOptions) *Error {
	e := &Error{
		status:   *status.New(codes.ResourceExhausted, opts.Error.Error()),
		logLevel: opts.LogLevel,
	}

	e.AddQuotaViolations([]QuotaViolation{opts.QuotaViolation})
	return e
}

type ResourcesExhaustedOptions struct {
	Error           error
	QuotaViolations []QuotaViolation
	LogLevel        LogLevel
}

func (_ factory) newResourcesExhausted(opts ResourcesExhaustedOptions) *Error {
	e := &Error{
		status:   *status.New(codes.ResourceExhausted, opts.Error.Error()),
		logLevel: opts.LogLevel,
	}

	e.AddQuotaViolations(opts.QuotaViolations)
	return e
}

func (_ factory) newCanceledError(logLevel LogLevel) *Error {
	const msg = "request canceled by the client"
	e := &Error{
		status:   *status.New(codes.Canceled, msg),
		logLevel: logLevel,
	}
	return e
}

type ErrorWithHiddenDetailsOptions struct {
	Error    error
	LogLevel LogLevel
}

func (f factory) newDataLoss(opts ErrorWithHiddenDetailsOptions) *Error {
	return f.newErrorWithHiddenDetails(codes.DataLoss, opts)
}

func (f factory) newUnknown(opts ErrorWithHiddenDetailsOptions) *Error {
	return f.newErrorWithHiddenDetails(codes.Unknown, opts)
}

func (f factory) newInternalError(opts ErrorWithHiddenDetailsOptions) *Error {
	return f.newErrorWithHiddenDetails(codes.Internal, opts)
}

func (f factory) newNotImplemented(logLevel LogLevel) *Error {
	const msg = "not implemented"
	e := &Error{
		status:   *status.New(codes.Unimplemented, msg),
		logLevel: logLevel,
	}
	return e
}

func (f factory) newUnavailable(opts ErrorWithHiddenDetailsOptions) *Error {
	return f.newErrorWithHiddenDetails(codes.Unavailable, opts)
}

func (f factory) newDeadlineExceeded(opts ErrorWithHiddenDetailsOptions) *Error {
	return f.newErrorWithHiddenDetails(codes.DeadlineExceeded, opts)
}

/* ------------------------- Factory helper methods ------------------------- */

func (_ factory) newBadRequest(msg string, opts BadRequestViolationOptions) *Error {
	e := &Error{
		status:   *status.New(codes.InvalidArgument, msg),
		logLevel: opts.LogLevel,
	}

	e.AddBadRequestViolations([]BadRequestViolation{opts.Violation})
	return e
}

func (_ factory) newBadRequestBulk(msg string, opts BadRequestViolationsOptions) *Error {
	e := &Error{
		status:   *status.New(codes.InvalidArgument, msg),
		logLevel: opts.LogLevel,
	}

	e.AddBadRequestViolations(opts.Violations)
	return e
}

func (f factory) newErrorInfoError(code codes.Code, opts ErrorInfoOptions) *Error {
	e := &Error{
		status:   *status.New(code, opts.Error.Error()),
		logLevel: opts.LogLevel,
	}
	e.SetErrorInfo(f.domain, opts.Reason, opts.Metadata)
	return e
}

func (_ factory) newErrorWithHiddenDetails(code codes.Code, opts ErrorWithHiddenDetailsOptions) *Error {
	return &Error{
		status:        *status.New(code, opts.Error.Error()),
		logLevel:      opts.LogLevel,
		detailsHidden: true,
	}
}
