package http

// Aliases for the types in the generated code
// type (
// 	FieldViolation        = gen.FieldViolation
// 	PreconditionViolation = gen.PreconditionViolation
// 	ResourceInfo          = gen.ResourceInfo
// )

// func NewUnknownError(err error, details string) *Error {
// 	httpError := gen.NewUnknownError(details)
// 	return &Error{
// 		httpError: *httpError,
// 	}
// }

// func NewDeadlineExceeded() *Error {
// 	httpError := gen.NewDeadlineExceeded()
// 	return &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindDeadlineExceeded},
// 	}
// }

// func NewInternalError(ctx context.Context) *Error {
// 	httpError := gen.NewInternalError(ctx)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindInternal},
// 	}
// 	return e
// }

// func NewInvalidArgumentError(field, desc string) *Error {
// 	httpError := gen.NewInvalidArgumentError(field, desc)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindInvalidArgs},
// 	}
// 	return e
// }

// func NewInvalidArgumentErrors(violations []FieldViolation) *Error {
// 	httpError := gen.NewInvalidArgumentErrors(violations)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindInvalidArgs},
// 	}
// 	return e
// }

// func NewPreconditionFailure(description, subject, typ string) *Error {
// 	httpError := gen.NewPreconditionFailure(description, subject, typ)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindPreconditionViolation},
// 	}
// 	return e
// }

// func NewModelBindingError(err error) *Error {
// 	httpError := gen.NewModelBindingError(err)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindInvalidArgs},
// 	}
// 	return e
// }

// func NewUnauthenticatedError(err error, domain string) *Error {
// 	httpError := gen.NewUnauthenticatedError(err, domain)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindUnauthenticated},
// 	}
// 	return e
// }

// func NewNotFoundError(desc, rscName, rscType string) *Error {
// 	httpError := gen.NewNotFoundError(desc, rscName, rscType)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindNotFound},
// 	}
// 	return e
// }

// func NewNotFoundErrors(infos []ResourceInfo) *Error {
// 	httpError := gen.NewNotFoundErrors(infos)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindNotFound},
// 	}
// 	return e
// }

// func NewPermissionDeniedError(domain, reason string, metadata map[string]string) *Error { // nolint:unparam
// 	httpError := gen.NewPermissionDeniedError(domain, reason, metadata)
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindPermissionDenied},
// 	}
// 	return e
// }

// func NewCanceledError() *Error {
// 	httpError := gen.NewCanceledError()
// 	e := &Error{
// 		httpError: *httpError,
// 		Root:      xerror.Root{Kind: xerror.KindCanceled},
// 	}
// 	return e
// }
