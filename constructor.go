package xerror

/* -------------------------------------------------------------------------- */
/*                          Server-initialized errors                         */
/* -------------------------------------------------------------------------- */

func NewInvalidArgument(opts BadRequestOptions) *Error {
	return maker.newInvalidArgument(opts)
}

func NewBatchInvalidArgument(opts BadRequestBatchOptions) *Error {
	return maker.newInvalidArgumentErrors(opts)
}

func NewPreconditionFailure(opts PreconditionFailureOptions) *Error {
	return maker.newPreconditionFailure(opts)
}

func NewBatchPreconditionFailure(opts PreconditionFailureBatchOptions) *Error {
	return maker.newPreconditionFailures(opts)
}

func NewOutOfRange(opts BadRequestOptions) *Error {
	return maker.newOutOfRangeError(opts)
}

func NewBatchOutOfRange(opts BadRequestBatchOptions) *Error {
	return maker.newOutOfRangeErrors(opts)
}

func NewUnauthenticated(opts ErrorInfoOptions) *Error {
	return maker.newUnauthenticatedError(opts)
}

func NewPermissionDenied(opts ErrorInfoOptions) *Error {
	return maker.newPermissionDeniedError(opts)
}

func NewNotFound(opts NotFoundOptions) *Error {
	return maker.newNotFound(opts)
}

func NewBatchNotFound(opts NotFoundBatchOptions) *Error {
	return maker.newBatchNotFound(opts)
}

func NewAborted(opts ErrorInfoOptions) *Error {
	return maker.newAborted(opts)
}

func NewAlreadyExists(opts AlreadyExistsOptions) *Error {
	return maker.newAlreadyExists(opts)
}

func NewBatchAlreadyExists(opts AlreadyExistsBatchOptions) *Error {
	return maker.newBatchAlreadyExists(opts)
}

func NewResourceExhausted(opts ResourceExhaustedOptions) *Error {
	return maker.newResourceExhausted(opts)
}

func NewBatchResourceExhausted(opts ResourceExhaustedBatchOptions) *Error {
	return maker.newBatchResourceExhausted(opts)
}

func NewCanceled(logLevel LogLevel) *Error {
	return maker.newCanceledError(logLevel)
}

func NewDataLoss(opts SimpleOptions) *Error {
	return maker.newDataLoss(opts)
}

func NewUnknown(opts SimpleOptions) *Error {
	return maker.newUnknown(opts)
}

func NewInternal(opts SimpleOptions) *Error {
	return maker.newInternalError(opts)
}

func NewNotImplemented(logLevel LogLevel) *Error {
	return maker.newNotImplemented(logLevel)
}

func NewUnavailable(opts SimpleOptions) *Error {
	return maker.newUnavailable(opts)
}

func NewDeadlineExceeded(opts SimpleOptions) *Error {
	return maker.newDeadlineExceeded(opts)
}
