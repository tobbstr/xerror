package grpc

// func NewUnknownError(details string) *Error {
// 	st := status.New(codes.Unknown, details)
// 	return &Error{
// 		Root:   xerror.Root{Kind: xerror.KindUnknown},
// 		Status: *st,
// 	}
// }

// func NewDeadlineExceeded() *Error {
// 	st := status.New(codes.DeadlineExceeded, "")
// 	return &Error{
// 		Root:   xerror.Root{Kind: xerror.KindDeadlineExceeded},
// 		Status: *st,
// 	}
// }

// func NewInternalError(ctx context.Context) *Error {
// 	st := status.New(codes.Internal, "")
// 	e := &Error{
// 		Root:   xerror.Root{Kind: xerror.KindInternal},
// 		Status: *st,
// 	}
// 	return e
// }

// func NewInvalidArgumentError(field, desc string) *Error {
// 	// Create a BadRequest error detail with the field and description
// 	violations := []*errdetails.BadRequest_FieldViolation{{Field: field, Description: desc}}
// 	return NewInvalidArgumentErrors(violations)
// }

// func NewInvalidArgumentErrors(violations []*errdetails.BadRequest_FieldViolation) *Error {
// 	st := status.New(codes.InvalidArgument, "client provided invalid argument(s)")

// 	// Create the error detail
// 	badRequest := errdetails.BadRequest{FieldViolations: violations}

// 	// Attach the detail to the gRPC status
// 	stWithDetails, err := st.WithDetails(&badRequest)
// 	if err != nil {
// 		return &Error{
// 			Root:   xerror.Root{Kind: xerror.KindInvalidArgs},
// 			Status: *st,
// 		}
// 	}

// 	// Return the error with the detail attached
// 	return &Error{
// 		Root:   xerror.Root{Kind: xerror.KindInvalidArgs},
// 		Status: *stWithDetails,
// 	}
// }

// func NewPreconditionFailure(description, subject, typ string) *Error {
// 	st := status.New(codes.FailedPrecondition, "system state does not allow operation")

// 	// Create the error detail
// 	precondFailure := errdetails.PreconditionFailure{
// 		Violations: []*errdetails.PreconditionFailure_Violation{
// 			{Type: typ, Subject: subject, Description: description},
// 		},
// 	}

// 	// Attach the detail to the gRPC status
// 	stWithDetails, err := st.WithDetails(&precondFailure)
// 	if err != nil {
// 		return &Error{
// 			Root:   xerror.Root{Kind: xerror.KindInvalidArgs},
// 			Status: *st,
// 		}
// 	}

// 	// Return the error with the detail attached
// 	return &Error{
// 		Root:   xerror.Root{Kind: xerror.KindPreconditionViolation},
// 		Status: *stWithDetails,
// 	}
// }

// // func NewModelBindingError(err error) *Error {
// // 	httpError := gen.NewModelBindingError(err)
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindInvalidArgs},
// // 	}
// // 	return e
// // }

// // func NewUnauthenticatedError(err error, domain string) *Error {
// // 	httpError := gen.NewUnauthenticatedError(err, domain)
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindUnauthenticated},
// // 	}
// // 	return e
// // }

// // func NewNotFoundError(desc, rscName, rscType string) *Error {
// // 	httpError := gen.NewNotFoundError(desc, rscName, rscType)
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindNotFound},
// // 	}
// // 	return e
// // }

// // func NewNotFoundErrors(infos []ResourceInfo) *Error {
// // 	httpError := gen.NewNotFoundErrors(infos)
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindNotFound},
// // 	}
// // 	return e
// // }

// // func NewPermissionDeniedError(domain, reason string, metadata map[string]string) *Error { // nolint:unparam
// // 	httpError := gen.NewPermissionDeniedError(domain, reason, metadata)
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindPermissionDenied},
// // 	}
// // 	return e
// // }

// // func NewCanceledError() *Error {
// // 	httpError := gen.NewCanceledError()
// // 	e := &Error{
// // 		httpError: *httpError,
// // 		Root:      xerror.Root{Kind: xerror.KindCanceled},
// // 	}
// // 	return e
// // }
