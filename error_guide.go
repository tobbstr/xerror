package xerror

type errorGuide struct{}

// ErrorGuide implements a decision tree that helps developers choose the right error type for their use case.
// Errors are not supposed to be created using this guide, although it is possible. Instead, use the guide
// to find the right error type and then create the error using the appropriate factory function.
//
// IMPORTANT! Each method has a comment that explains when to use the error type. Please read the comments carefully
// before choosing an error type.
func ErrorGuide() errorGuide {
	return errorGuide{}
}

// ProblemWithRequest is used when the server encounters an issue with the client's request. For example, when the
// server cannot process the request due to invalid input, missing data, or other issues with the request itself.
func (errorGuide) ProblemWithRequest() requestIssue {
	return requestIssue{}
}

// ProblemWithServer is used when the server encounters an issue that prevents it from processing the client's request.
// For example, when the server is unable to process the request due to an internal error, a temporary issue, or other
// problems with the server itself.
func (errorGuide) ProblemWithServer() serverIssue {
	return serverIssue{}
}

type serverIssue struct{}

type requestIssue struct{}

// Cancelled is used when a request is cancelled by the client before the server has completed processing it.
// Suppose a client sends a request to a server to perform a long-running computation. After some time, the client
// decides to cancel the request because it no longer needs the result. The client sends a cancellation request to
// the server, and the server should then respond with a Cancelled error.
//
// This case is a "CANCELLED" error.
func (requestIssue) Cancelled() func() *Error {
	return maker.newCancelledError
}

// InvalidArgument is used when a request is rejected due to invalid input.
func (requestIssue) InvalidArgument() *invalidArgIssue {
	return &invalidArgIssue{}
}

type invalidArgIssue struct{}

// Other is the default invalid argument error type. It must be used when the error does not fit any of the other
// specialized invalid argument types.
//
// For example:
//  1. When a user provides an invalid value for an email address or phone number.
//
// This case is an "INVALID_ARGUMENT" error.
func (invalidArgIssue) Other() func(field, description string) *Error {
	return maker.newInvalidArgument
}

// OutOfRange is a specialized type of invalid argument that occurs when a value is outside the acceptable range.
//
// For example:
//  1. When a user attempts to set a value that is outside the acceptable range of values for a given field.
//     Such as when the argument is a page number for pagination, and the value is greater than the total number of
//     pages.
//
// This case is an "OUT_OF_RANGE" error.
func (invalidArgIssue) OutOfRange() func(field, description string) *Error {
	return maker.newOutOfRangeError
}

// NotFound is a specialized type of invalid argument that occurs when a requested resource cannot be found.
//
// For example:
//  1. When a user attempts to access a resource that does not exist, such as a non-existent file or record. Such as
//     when a user makes a request to get a specific resource by id, but the resource with that id does not exist.
//
// This case is a "NOT_FOUND" error.
func (invalidArgIssue) NotFound() func(info ResourceInfo) *Error {
	return maker.newNotFound
}

// DataLoss is a specialized type of invalid argument that occurs when the integrity of data is compromised.
//
// For example:
//  1. When an application fails to serialize or deserialize request data correctly, resulting in data loss
//     during transmission.
//  2. When a system verifies the integrity of request data using checksums or hashes and finds
//     that the data does not match its expected state, indicating possible data loss or corruption.
//
// This case is a "DATA_LOSS" error.
func (invalidArgIssue) DataLoss() func(opts ErrorInfoOptions) *Error {
	return maker.newRequestDataLoss
}

// PermissionDenied is used when a user's identity has been verified (authenticated), but the user does not have the
// necessary permissions to perform the requested action.
//
// For example:
//  1. When a user attempts to access a resource that they do not have permission to access, such as a file or record
//     that is restricted to certain users.
//
// This case is a "PERMISSION_DENIED" error.
func (requestIssue) PermissionDenied() func(opts ErrorInfoOptions) *Error {
	return maker.newPermissionDeniedError
}

// Unauthenticated is used when the check for a user's identity fails. For example, when a user attempts to access a
// resource without providing the necessary credentials, or when the user's credentials are invalid or expired.
//
// This case is an "UNAUTHENTICATED" error.
func (requestIssue) Unauthenticated() func(opts ErrorInfoOptions) *Error {
	return maker.newUnauthenticatedError
}

// ServerDataLoss is used when the server encounters an issue that results in data loss.
//
// For example:
//  1. When a system detects that the data stored in a database, file system, or any storage medium has been corrupted
//     and cannot be recovered.
//  2. When an application fails to serialize or deserialize data (not received by the caller) correctly.
//  3. When a system verifies the integrity of data (not received from the caller) using checksums or hashes and
//     finds that the data does not match its expected state, indicating possible data loss or corruption.
//  4. When an attempt to recover data from a backup or a redundant system fails to restore all data correctly,
//     resulting in partial or complete data loss.
//
// This case is a "DATA_LOSS" error.
func (serverIssue) ServerDataLoss() func(err error) *Error {
	return maker.newServerDataLoss
}

// PreconditionFailed is used when a request fails because a precondition for the operation was not met.
// A precondition is a condition that must be true before an operation can be executed.
func (serverIssue) PreconditionFailed() precondFailureIssue {
	return precondFailureIssue{}
}

type precondFailureIssue struct{}

// Other is used when the other precondition failure error types do not apply.
//
// For example:
//  1. When an operation fails because it's only allowed for users between the ages of 18 and 65.
//  2. When an operation fails because the user has not agreed to the terms and conditions.
//
// This case is a "FAILED_PRECONDITION" error.
func (precondFailureIssue) Other() func(subject, typ, description string) *Error {
	return maker.newPreconditionFailure
}

// Aborted is a specialized form of precondition failure and is used to indicate that an operation was aborted,
// typically due to a concurrency issue. This code is used when the operation needs to be retried (at a higher level)
// because it was not completed successfully.
//
// For example:
//  1. When two transactions are attempting to update the same data simultaneously, causing a conflict.
//  2. When an optimistic locking mechanism detects that the data has been modified by another process since it was
//     last read.
//  3. When a distributed transaction fails due to an issue with the transaction coordinator or one of the
//     participating nodes.
//  4. When a concurrency control mechanism (like a semaphore or lock) cannot be acquired.
//
// This case is an "ABORTED" error.
func (precondFailureIssue) Aborted() func(opts ErrorInfoOptions) *Error {
	return maker.newAborted
}

// AlreadyExists is a specialized form of precondition failure and is used when an attempt to create a resource fails
// because the resource already exists. This status code is appropriate in scenarios where the client's request
// attempts to create a new entity, but the entity is already present in the system.
//
// For example:
//  1. When a user tries to register with a username that is already taken.
//  2. When a user tries to upload a file with a name that already exists in the directory.
//  3. When a user tries to register with an email address that is already associated with another account.
//  4. When an attempt is made to insert a record into a database with a primary key that already exists.
//
// This case is an "ALREADY_EXISTS" error.
func (precondFailureIssue) AlreadyExists() func(info ResourceInfo) *Error {
	return maker.newAlreadyExists
}

// ResourceExhausted is a specialized form of precondition failure and is used when a resource has been exhausted,
// meaning the server cannot complete the request due to a lack of resources
func (precondFailureIssue) ResourceExhausted() resourceExhaustedIssue {
	return resourceExhaustedIssue{}
}

type resourceExhaustedIssue struct{}

// Other is used when the other resource exhausted error type does not apply.
//
// For example:
//  1. When the server cannot process a request due to insufficient memory or storage.
//
// This case is a "RESOURCE_EXHAUSTED" error.
func (resourceExhaustedIssue) Other() func(opts ErrorInfoOptions) *Error {
	return maker.newResourceExhausted
}

// QuotaFailure is a specialized form of resource exhausted error and is used when an alloted quota or limit
// (e.g., rate limit) has been exceeded.
//
// For example:
//  1. When a client exceeds its allotted quota for API requests within a given time period.
//  2. When a client makes too many requests in a short period of time, exceeding the rate limit.
//  3. When a client exceeds the number of allowed concurrent requests.
//
// This case is a "RESOURCE_EXHAUSTED" error.
func (resourceExhaustedIssue) QuotaFailure() func(subject, description string) *Error {
	return maker.newQuotaFailure
}

// Unknown is used for errors that are unknown or that do not fit any other standard error categories. This is a
// catch-all error code for unexpected or unforeseen errors.
//
// For example:
//  1. When a third-party library used by the server throws an error that cannot be mapped to a specific gRPC error
//     code.
//  2. When the server receives data in an unexpected format that it cannot process, and this condition does not fit
//     other specific error codes.
//
// This case is an "UNKNOWN" error.
func (serverIssue) Unknown() func(err error) *Error {
	return maker.newUnknown
}

// Internal is used when the server encounters an unexpected condition that prevents it from fulfilling the request.
// It typically implies that something went wrong on the server side, which is not the client's fault, and that
// doesn't fit any other error category.
//
// For example:
//  1. When the server cannot connect to the database due to internal issues, such as network problems or database
//     server downtime.
//  2. When the server's configuration is incorrect or incomplete, leading to an operational failure.
//  3. When a critical internal dependency (e.g., a microservice) fails to respond or returns an error, and the
//     returned error status must not be used, then a mapping to a generic INTERNAL error is required.
//
// This case is an "INTERNAL" error.
func (serverIssue) Internal() func(err error) *Error {
	return maker.newInternalError
}

// NotImplemented is used when an operation is not supported by the server. This can be due to the feature not being
// implemented yet, the server intentionally not supporting the feature, or the method being deprecated or removed.
//
// For example:
//  1. When the client calls a method that has been defined in the API but has not been implemented by the server.
//  2. When the client calls a method that has been deprecated and removed from the server's implementation.
//  3. When the client requests an operation that the server intentionally does not support.
//  4. When the client calls an API method that is not available in the current version of the server software.
//  5. When the server partially supports a feature but the specific requested method or capability is not yet
//     implemented.
//
// This case is a "NOT_IMPLEMENTED" error.
func (serverIssue) NotImplemented() func() *Error {
	return maker.newNotImplemented
}

// Unavailable is used when the whole server is currently unavailable, not just the requested operation.
// This typically means that the server is temporarily unable to handle the request due to reasons like
// server overload, maintenance, network issues, or a dependency service being down.
//
// For example:
//  1. When the server is temporarily unable to handle the request due to high load.
//  2. When the server is down for scheduled maintenance.
//  3. When there are network issues preventing the server from being reachable.
//  4. When a critical dependency service is down and the server cannot fulfill the request.
//  5. When the server is temporarily unavailable because it is restarting.
//
// This case is an "UNAVAILABLE" error.
func (serverIssue) Unavailable() func(err error) *Error {
	return maker.newUnavailable
}

// DeadlineExceeded is used when the request took too long to complete and has exceeded the time allocated for it.
// This can happen due to various reasons, such as long-running operations, network latency, or server performance
// issues.
//
// For example:
//  1. When a database query takes longer than the allocated timeout to complete.
//  2. When a call to an external API takes longer than the maximum allowed time.
//  3. When data processing takes longer than the client is willing to wait.
//  4. When network latency causes the request to exceed its deadline.
//  5. When a batch processing job takes longer than the specified time limit.
//
// This case is a "DEADLINE_EXCEEDED" error.
func (serverIssue) DeadlineExceeded() func() *Error {
	return maker.newDeadlineExceeded
}
