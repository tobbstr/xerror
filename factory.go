package xerror

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrFailedToAddErrorDetails = errors.New("failed to add error details")

const (
	msgInvalidArg           = "one request arguments was invalid"
	msgInvalidArgs          = "one or more request arguments were invalid"
	msgPreconditionFailure  = "one request precondition failed"
	msgPreconditionFailures = "one or more request preconditions failed"
	msgOutOfRange           = "one request argument was out of range"
	msgOutOfRangeErrors     = "one or more request arguments were out of range"
)

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

// BadRequestViolation is a message type used to describe a single bad request field.
type BadRequestViolation struct {
	// Field is a path that leads to a field in the request body. The value will be a
	// sequence of dot-separated identifiers that identify a protocol buffer
	// field.
	//
	// Consider the following:
	//
	//	message CreateContactRequest {
	//	  message EmailAddress {
	//	    enum Type {
	//	      TYPE_UNSPECIFIED = 0;
	//	      HOME = 1;
	//	      WORK = 2;
	//	    }
	//
	//	    optional string email = 1;
	//	    repeated EmailType type = 2;
	//	  }
	//
	//	  string full_name = 1;
	//	  repeated EmailAddress email_addresses = 2;
	//	}
	//
	// In this example, in proto `field` could take one of the following values:
	//
	//   - `full_name` for a violation in the `full_name` value
	//   - `email_addresses[1].email` for a violation in the `email` field of the
	//     first `email_addresses` message
	//   - `email_addresses[3].type[2]` for a violation in the second `type`
	//     value in the third `email_addresses` message.
	//
	// In JSON, the same values are represented as:
	//
	//   - `fullName` for a violation in the `fullName` value
	//   - `emailAddresses[1].email` for a violation in the `email` field of the
	//     first `emailAddresses` message
	//   - `emailAddresses[3].type[2]` for a violation in the second `type`
	//     value in the third `emailAddresses` message.
	Field string
	// Description is a description of why the request element is bad.
	Description string
}

func (f factory) newInvalidArgument(field, description string) *Error {
	// TODO(tobbstr): Add a function that accepts a the request object field and then it returns the field name.
	// Ex. Instead of the user having to construct the field name such as "person.ownedDogs[1].name", they can
	// pass the object and the function returns the field name.
	return f.newBadRequest(msgInvalidArg, BadRequestViolation{Field: field, Description: description})
}

func (f factory) newInvalidArgumentErrors(violations []BadRequestViolation) *Error {
	return f.newBatchBadRequest(msgInvalidArgs, violations)
}

// PreconditionViolation is a message type used to describe a single precondition failure.
type PreconditionViolation struct {
	// Description is a description of how the precondition failed. Developers can use this
	// description to understand how to fix the failure.
	//
	// For example: "Terms of service not accepted".
	Description string
	// Subject is the subject, relative to the type, that failed.
	// For example, "google.com/cloud" relative to the "TOS" type would indicate
	// which terms of service is being referenced.
	Subject string
	// Typ is the type of PreconditionFailure. It's recommended using a service-specific
	// enum type to define the supported precondition violation subjects. For
	// example, "TOS" for "Terms of Service violation".
	Typ string
}

func (f factory) newPreconditionFailure(subject, typ, description string) *Error {
	e := &Error{
		status:   *status.New(codes.FailedPrecondition, msgPreconditionFailure),
		logLevel: LogLevelWarn,
	}
	_ = e.AddPreconditionViolations([]PreconditionViolation{{Description: description, Subject: subject, Typ: typ}})
	return e
}

func (f factory) newPreconditionFailures(violations []PreconditionViolation) *Error {
	e := &Error{
		status:   *status.New(codes.FailedPrecondition, msgPreconditionFailures),
		logLevel: LogLevelWarn,
	}

	_ = e.AddPreconditionViolations(violations)
	return e
}

func (f factory) newOutOfRangeError(field, description string) *Error {
	e := &Error{
		status:   *status.New(codes.OutOfRange, msgOutOfRange),
		logLevel: LogLevelInfo,
	}
	_ = e.AddBadRequestViolations([]BadRequestViolation{{Field: field, Description: description}})
	return e
}

func (f factory) newOutOfRangeErrors(violations []BadRequestViolation) *Error {
	e := &Error{
		status:   *status.New(codes.OutOfRange, msgOutOfRangeErrors),
		logLevel: LogLevelInfo,
	}
	_ = e.AddBadRequestViolations(violations)
	return e
}

type ErrorInfoOptions struct {
	// Error is the error that occurred.
	Error error
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
	return f.newErrorInfoError(codes.Unauthenticated, LogLevelInfo, opts)
}

func (f factory) newPermissionDeniedError(opts ErrorInfoOptions) *Error { // nolint:unparam
	e := f.newErrorInfoError(codes.PermissionDenied, LogLevelInfo, opts)
	return e
}

// ResourceInfo describes the resource that is being accessed.
type ResourceInfo struct {
	// Description describes what error is encountered when accessing this resource.
	// For example, updating a cloud project may require the `writer` permission
	// on the developer console project.
	Description string
	// ResourceName is the name of the resource being accessed.  For example, a shared calendar
	// name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current
	// error is
	// [google.rpc.Code.PERMISSION_DENIED][google.rpc.Code.PERMISSION_DENIED].
	ResourceName string
	// ResourceType is a name for the type of resource being accessed, e.g. "sql table",
	// "cloud storage bucket", "file", "Google calendar"; or the type URL
	// of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic".
	ResourceType string
	// Owner is the owner of the resource (optional).
	// For example, "user:<owner email>" or "project:<Google developer project
	// id>".
	Owner string
}

func (_ factory) newNotFound(info ResourceInfo) *Error {
	const msg = "requested resource not found"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: LogLevelInfo,
	}

	_ = e.AddResourceInfos([]ResourceInfo{info})
	return e
}

func (_ factory) newBatchNotFound(infos []ResourceInfo) *Error {
	const msg = "requested resources not found"
	e := &Error{
		status:   *status.New(codes.NotFound, msg),
		logLevel: LogLevelInfo,
	}

	_ = e.AddResourceInfos(infos)
	return e
}

func (f factory) newAborted(opts ErrorInfoOptions) *Error {
	return f.newErrorInfoError(codes.Aborted, LogLevelWarn, opts)
}

func (f factory) newAlreadyExists(info ResourceInfo) *Error {
	const msg = "resource already exists"
	e := &Error{
		status:   *status.New(codes.AlreadyExists, msg),
		logLevel: LogLevelInfo,
	}

	_ = e.AddResourceInfos([]ResourceInfo{info})
	return e
}

func (_ factory) newAlreadyExistsBatch(infos []ResourceInfo) *Error {
	const msg = "resources already exist"
	e := &Error{
		status:   *status.New(codes.AlreadyExists, msg),
		logLevel: LogLevelInfo,
	}
	_ = e.AddResourceInfos(infos)
	return e
}

func (_ factory) newQuotaFailure(subject, description string) *Error {
	e := &Error{
		status:   *status.New(codes.ResourceExhausted, "the request cannot be completed because the quota has been exhausted"),
		logLevel: LogLevelInfo,
	}

	_ = e.AddQuotaViolations([]QuotaViolation{{Subject: subject, Description: description}})
	return e
}

func (_ factory) newQuotaFailureBatch(violations []QuotaViolation) *Error {
	e := &Error{
		status:   *status.New(codes.ResourceExhausted, "the request cannot be completed because the quota has been exhausted"),
		logLevel: LogLevelInfo,
	}

	_ = e.AddQuotaViolations(violations)
	return e
}

func (_ factory) newResourceExhausted(opts ErrorInfoOptions) *Error {
	return maker.newErrorInfoError(codes.ResourceExhausted, LogLevelWarn, opts)
}

func (_ factory) newCancelledError() *Error {
	const msg = "request cancelled by the client"
	e := &Error{
		status: *status.New(codes.Canceled, msg),
	}
	return e
}

func (f factory) newServerDataLoss(err error) *Error {
	var msg string
	if err == nil {
		msg = "server data loss"
	} else {
		msg = err.Error()
	}
	return f.newErrorWithDetailsHidden(codes.DataLoss, msg, LogLevelError)
}

func (_ factory) newRequestDataLoss(opts ErrorInfoOptions) *Error {
	return maker.newErrorInfoError(codes.DataLoss, LogLevelInfo, opts)
}

func (f factory) newUnknown(err error) *Error {
	var msg string
	if err == nil {
		msg = "something unknown happened"
	} else {
		msg = err.Error()
	}
	return f.newErrorWithDetailsHidden(codes.Unknown, msg, LogLevelError)
}

func (f factory) newInternalError(err error) *Error {
	var msg string
	if err == nil {
		msg = "an internal server error happened"
	} else {
		msg = err.Error()
	}
	return f.newErrorWithDetailsHidden(codes.Internal, msg, LogLevelError)
}

func (f factory) newNotImplemented() *Error {
	const msg = "not implemented"
	e := &Error{
		status:   *status.New(codes.Unimplemented, msg),
		logLevel: LogLevelInfo,
	}
	return e
}

func (f factory) newUnavailable(err error) *Error {
	var msg string
	if err == nil {
		msg = "the operation is currently unavailable"
	} else {
		msg = err.Error()
	}
	return f.newErrorWithDetailsHidden(codes.Unavailable, msg, LogLevelInfo)
}

func (f factory) newDeadlineExceeded() *Error {
	return f.newErrorWithDetailsHidden(
		codes.DeadlineExceeded,
		"the operation timed out (it might have succeeded though)",
		LogLevelWarn,
	)
}

/* ------------------------- Factory helper methods ------------------------- */

func (_ factory) newBadRequest(msg string, violation BadRequestViolation) *Error {
	e := &Error{
		status:   *status.New(codes.InvalidArgument, msg),
		logLevel: LogLevelInfo,
	}

	_ = e.AddBadRequestViolations([]BadRequestViolation{violation})
	return e
}

func (_ factory) newBatchBadRequest(msg string, violations []BadRequestViolation) *Error {
	e := &Error{
		status:   *status.New(codes.InvalidArgument, msg),
		logLevel: LogLevelInfo,
	}
	_ = e.AddBadRequestViolations(violations)
	return e
}

func (f factory) newErrorInfoError(code codes.Code, logLevel LogLevel, opts ErrorInfoOptions) *Error {
	if opts.Error == nil {
		return nil
	}
	e := &Error{
		status:   *status.New(code, opts.Error.Error()),
		logLevel: logLevel,
	}
	_ = e.SetErrorInfo(f.domain, opts.Reason, opts.Metadata)
	return e
}

func (_ factory) newErrorWithDetailsHidden(code codes.Code, msg string, logLevel LogLevel) *Error {
	var lvl LogLevel
	switch logLevel {
	case LogLevelUnspecified:
		lvl = LogLevelError
	default:
		lvl = logLevel
	}
	return &Error{
		status:        *status.New(code, msg),
		logLevel:      lvl,
		detailsHidden: true,
	}
}
