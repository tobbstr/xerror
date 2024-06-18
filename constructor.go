package xerror

/* -------------------------------------------------------------------------- */
/*                          Server-initialized errors                         */
/* -------------------------------------------------------------------------- */

// NewInvalidArgument creates a new InvalidArgument error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewInvalidArgument(opts BadRequestOptions) *Error {
	return maker.newInvalidArgument(opts)
}

// NewInvalidArgumentBatch creates a new InvalidArgument error. This is the batch version that adds multiple field
// violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewInvalidArgumentBatch(opts BadRequestBatchOptions) *Error {
	return maker.newInvalidArgumentErrors(opts)
}

// NewFailedPrecondition creates a new FailedPrecondition error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewPreconditionFailure(opts PreconditionFailureOptions) *Error {
	return maker.newPreconditionFailure(opts)
}

// NewFailedPreconditionBatch creates a new FailedPrecondition error. This is the batch version that adds multiple
// precondition violations.
func NewPreconditionFailureBatch(opts PreconditionFailureBatchOptions) *Error {
	return maker.newPreconditionFailures(opts)
}

// NewOutOfRange creates a new OutOfRange error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewOutOfRange(opts BadRequestOptions) *Error {
	return maker.newOutOfRangeError(opts)
}

// NewOutOfRangeBatch creates a new OutOfRange error. This is the batch version that adds multiple field violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewOutOfRangeBatch(opts BadRequestBatchOptions) *Error {
	return maker.newOutOfRangeErrors(opts)
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
func NewNotFound(opts NotFoundOptions) *Error {
	return maker.newNotFound(opts)
}

// NewNotFoundBatch creates a new NotFound error. This is the batch version that adds information about multiple
// resources that were not found.
//
// For when to use this, see the ErrorGuide function for more information.
func NewNotFoundBatch(opts NotFoundBatchOptions) *Error {
	return maker.newBatchNotFound(opts)
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
func NewAlreadyExists(opts AlreadyExistsOptions) *Error {
	return maker.newAlreadyExists(opts)
}

// NewAlreadyExistsBatch creates a new AlreadyExists error. This is the batch version that adds information about
// multiple resources that already exist.
//
// For when to use this, see the ErrorGuide function for more information.
func NewAlreadyExistsBatch(opts AlreadyExistsBatchOptions) *Error {
	return maker.newAlreadyExistsBatch(opts)
}

// NewQuotaFailure creates a new QuotaFailure error, which is a specialized version of a resource exhausted error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewQuotaFailure(opts QuotaFailureOptions) *Error {
	return maker.newQuotaFailure(opts)
}

// NewQuotaFailureBatch creates a new QuotaFailure error. This is the batch version that adds multiple quota violations.
//
// For when to use this, see the ErrorGuide function for more information.
func NewQuotaFailureBatch(opts QuotaFailureBatchOptions) *Error {
	return maker.newQuotaFailureBatch(opts)
}

// NewResourceExhausted creates a new ResourceExhausted error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewResourceExhausted(opts ErrorInfoOptions) *Error {
	return maker.newResourceExhausted(opts)
}

// NewCanceled creates a new Canceled error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewCanceled(logLevel LogLevel) *Error {
	return maker.newCanceledError(logLevel)
}

// NewServerDataLoss creates a new DataLoss error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewServerDataLoss(opts SimpleOptions) *Error {
	return maker.newServerDataLoss(opts)
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
func NewUnknown(opts SimpleOptions) *Error {
	return maker.newUnknown(opts)
}

// NewInternal creates a new Internal error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewInternal(opts SimpleOptions) *Error {
	return maker.newInternalError(opts)
}

// NewNotImplemented creates a new NotImplemented error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewNotImplemented(logLevel LogLevel) *Error {
	return maker.newNotImplemented(logLevel)
}

// NewUnavailable creates a new Unavailable error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewUnavailable(opts SimpleOptions) *Error {
	return maker.newUnavailable(opts)
}

// NewDeadlineExceeded creates a new DeadlineExceeded error.
//
// For when to use this, see the ErrorGuide function for more information.
func NewDeadlineExceeded(opts SimpleOptions) *Error {
	return maker.newDeadlineExceeded(opts)
}
