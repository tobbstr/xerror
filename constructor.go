package xerror

/* -------------------------------------------------------------------------- */
/*                          Server-initialized errors                         */
/* -------------------------------------------------------------------------- */

// NewInvalidArgument creates a new InvalidArgument error.
//
// Parameters:
//
//	Field:
//		A path that leads to a field in the request body. The value will be a
//		sequence of dot-separated identifiers that identify a protocol buffer
//		field.
//
//		Consider the following:
//
//		message CreateContactRequest {
//		  message EmailAddress {
//		    enum Type {
//		      TYPE_UNSPECIFIED = 0;
//		      HOME = 1;
//		      WORK = 2;
//		    }
//
//		    optional string email = 1;
//		    repeated EmailType type = 2;
//		  }
//
//		  string full_name = 1;
//		  repeated EmailAddress email_addresses = 2;
//		}
//
//		In this example, in proto, `field` could take one of the following values:
//
//		  - `full_name` for a violation in the `full_name` value
//		  - `email_addresses[1].email` for a violation in the `email` field of the
//		    first `email_addresses` message
//		  - `email_addresses[3].type[2]` for a violation in the second `type`
//		    value in the third `email_addresses` message.
//
//		In JSON, the same values are represented as:
//
//		  - `fullName` for a violation in the `fullName` value
//		  - `emailAddresses[1].email` for a violation in the `email` field of the
//		    first `emailAddresses` message
//		  - `emailAddresses[3].type[2]` for a violation in the second `type`
//		    value in the third `emailAddresses` message.
//
//	Description:
//		A description of why the request element is bad. See the below examples:
//
//		  - `fullName`: "The full name of the contact. It should include both first
//		    and last names. Example: `John Doe`".
//		  - `emailAddresses[3].type[2]`: "The type of the email address. It can be HOME,
//		    WORK, or unspecified. Example: [HOME, WORK]".
//
// For when to use this error type, see the ErrorGuide function for more information.
func NewInvalidArgument(field, description string) *Error {
	return maker.newInvalidArgument(field, description)
}

// NewInvalidArgumentBatch creates a new InvalidArgument error. This is the batch version that adds multiple field
// violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewInvalidArgumentBatch(violations []BadRequestViolation) *Error {
	return maker.newInvalidArgumentErrors(violations)
}

// NewFailedPrecondition creates a new FailedPrecondition error.
//
// Parameters:
//
//	Description:
//		Is a description of how the precondition failed. Developers can use this
//		description to understand how to fix the failure.  For example, "Terms of service not accepted".
//
//	Subject:
//		Is the subject, relative to the type, that failed.
//		For example, "google.com/cloud" relative to the "TOS" type would indicate
//		which terms of service is being referenced.
//
//	Typ:
//		Is the type of PreconditionFailure. It is recommended to use a service-specific
//		enum type to define the supported precondition violation subjects. For
//		example, "TOS" for "Terms of Service violation".
//
// For when to use this error type, see the ErrorGuide function for more information.
func NewPreconditionFailure(subject, typ, description string) *Error {
	return maker.newPreconditionFailure(subject, typ, description)
}

// NewFailedPreconditionBatch creates a new FailedPrecondition error. This is the batch version that adds multiple
// precondition violations.
func NewPreconditionFailureBatch(violations []PreconditionViolation) *Error {
	return maker.newPreconditionFailures(violations)
}

// NewOutOfRange creates a new OutOfRange error.
//
// Parameters: see NewInvalidArgument
//
// For when to use this, see the ErrorGuide function for more information.
func NewOutOfRange(field, description string) *Error {
	return maker.newOutOfRangeError(field, description)
}

// NewOutOfRangeBatch creates a new OutOfRange error. This is the batch version that adds multiple field violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewOutOfRangeBatch(violations []BadRequestViolation) *Error {
	return maker.newOutOfRangeErrors(violations)
}

// NewUnauthenticated creates a new Unauthenticated error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewUnauthenticated(opts ErrorInfoOptions) *Error {
	return maker.newUnauthenticatedError(opts)
}

// NewPermissionDenied creates a new PermissionDenied error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewPermissionDenied(opts ErrorInfoOptions) *Error {
	return maker.newPermissionDeniedError(opts)
}

// NewNotFound creates a new NotFound error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewNotFound(info ResourceInfo) *Error {
	return maker.newNotFound(info)
}

// NewNotFoundBatch creates a new NotFound error. This is the batch version that adds information about multiple
// resources that were not found.
//
// For when to use this, see the ErrorGuide function for more information.
func NewNotFoundBatch(infos []ResourceInfo) *Error {
	return maker.newBatchNotFound(infos)
}

// NewAborted creates a new Aborted error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewAborted(opts ErrorInfoOptions) *Error {
	return maker.newAborted(opts)
}

// NewAlreadyExists creates a new AlreadyExists error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewAlreadyExists(info ResourceInfo) *Error {
	return maker.newAlreadyExists(info)
}

// NewAlreadyExistsBatch creates a new AlreadyExists error. This is the batch version that adds information about
// multiple resources that already exist.
//
// For when to use this, see the ErrorGuide function for more information.
func NewAlreadyExistsBatch(infos []ResourceInfo) *Error {
	return maker.newAlreadyExistsBatch(infos)
}

// NewQuotaFailure creates a new QuotaFailure error, which is a specialized version of a resource exhausted error.
//
// Parameters:
//
//	Subject:
//		The subject on which the quota check failed.
//		For example, "clientip:<ip address of client>" or "project:<Google
//		developer project id>".
//	Description:
//		Description of how the quota check failed. Clients can use this
//		description to find more about the quota configuration in the service's
//		public documentation.
//		For example: "Service disabled" or "Daily Limit for read operations
//		exceeded".
//
// For when to use this, see the ErrorGuide function for more information.
func NewQuotaFailure(subject, description string) *Error {
	return maker.newQuotaFailure(subject, description)
}

// NewQuotaFailureBatch creates a new QuotaFailure error. This is the batch version that adds multiple quota violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewQuotaFailureBatch(violations []QuotaViolation) *Error {
	return maker.newQuotaFailureBatch(violations)
}

// NewResourceExhausted creates a new ResourceExhausted error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewResourceExhausted(opts ErrorInfoOptions) *Error {
	return maker.newResourceExhausted(opts)
}

// NewCancelled creates a new Cancelled error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewCancelled() *Error {
	return maker.newCancelledError()
}

// NewServerDataLoss creates a new DataLoss error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewServerDataLoss(err error) *Error {
	return maker.newServerDataLoss(err)
}

// NewRequestDataLoss creates a new DataLoss error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewRequestDataLoss(opts ErrorInfoOptions) *Error {
	return maker.newRequestDataLoss(opts)
}

// NewUnknown creates a new Unknown error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewUnknown(err error) *Error {
	return maker.newUnknown(err)
}

// NewInternal creates a new Internal error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewInternal(err error) *Error {
	return maker.newInternalError(err)
}

// NewNotImplemented creates a new NotImplemented error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewNotImplemented() *Error {
	return maker.newNotImplemented()
}

// NewUnavailable creates a new Unavailable error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewUnavailable(err error) *Error {
	return maker.newUnavailable(err)
}

// NewDeadlineExceeded creates a new DeadlineExceeded error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewDeadlineExceeded() *Error {
	return maker.newDeadlineExceeded()
}
