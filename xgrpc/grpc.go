package xgrpc

import (
	"context"
	"errors"

	"github.com/tobbstr/xerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryXErrorInterceptor is a gRPC server unary interceptor that unwraps the XError and returns the wrapped
// error status. It also removes sensitive details from errors if they are marked as hidden.
//
// This interceptor must be used by gRPC servers if they are returning xerrors.
func UnaryXErrorInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Call the handler
	resp, err := handler(ctx, req)
	if err != nil {
		var xerr *xerror.Error
		if !errors.As(err, &xerr) {
			return resp, err
		}
		if xerr.IsDetailsHidden() {
			_ = xerr.RemoveSensitiveDetails()
		}
		return resp, xerr.Status().Err()
	}
	return resp, err
}

// ErrorFrom is a convenience function that creates a new xerror from a gRPC error.
//
// Ex.
//
//	err := othersystempb.SomeMethod(ctx, req)
//	if err != nil {
//	  return grpc.ErrorFrom(err).AddVar("requested_id", req.Id)
//	}
func ErrorFrom(err error) *xerror.Error {
	if err == nil {
		return nil
	}
	return new(xerror.Error).SetStatus(status.Convert(err))
}
