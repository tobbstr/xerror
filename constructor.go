package xerror

/* -------------------------------------------------------------------------- */
/*                                Constructors                                */
/* -------------------------------------------------------------------------- */

func NewInvalidArgument(opts BadRequestViolationOptions) *Error {
	return maker.newInvalidArgument(opts)
}

func NewInvalidArguments(opts BadRequestViolationsOptions) *Error {
	return maker.newInvalidArgumentErrors(opts)
}

func NewPreconditionFailure(opts PreconditionFailureOptions) *Error {
	return maker.newPreconditionFailure(opts)
}

func NewPreconditionFailures(opts PreconditionFailuresOptions) *Error {
	return maker.newPreconditionFailures(opts)
}

func NewOutOfRange(opts BadRequestViolationOptions) *Error {
	return maker.newOutOfRangeError(opts)
}

func NewOutOfRangeBulk(opts BadRequestViolationsOptions) *Error {
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

func NewNotFoundBulk(opts NotFoundBulkOptions) *Error {
	return maker.newNotFoundBulk(opts)
}

func NewAborted(opts ErrorInfoOptions) *Error {
	return maker.newAborted(opts)
}

func NewAlreadyExists(opts AlreadyExistsOptions) *Error {
	return maker.newAlreadyExists(opts)
}

func NewAlreadyExistsBulk(opts AlreadyExistsBulkOptions) *Error {
	return maker.newAlreadyExistsBulk(opts)
}

func NewResourceExhausted(opts ResourceExhaustedOptions) *Error {
	return maker.newResourceExhausted(opts)
}

func NewCanceled(logLevel LogLevel) *Error {
	return maker.newCanceledError(logLevel)
}

func NewDataLoss(opts ErrorWithHiddenDetailsOptions) *Error {
	return maker.newDataLoss(opts)
}

func NewUnknown(opts ErrorWithHiddenDetailsOptions) *Error {
	return maker.newUnknown(opts)
}

func NewInternal(opts ErrorWithHiddenDetailsOptions) *Error {
	return maker.newInternalError(opts)
}

func NewNotImplemented(logLevel LogLevel) *Error {
	return maker.newNotImplemented(logLevel)
}

func NewUnavailable(opts ErrorWithHiddenDetailsOptions) *Error {
	return maker.newUnavailable(opts)
}

func NewDeadlineExceeded(opts ErrorWithHiddenDetailsOptions) *Error {
	return maker.newDeadlineExceeded(opts)
}
