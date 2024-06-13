package grpc

import (
	"context"
	"errors"

	"github.com/tobbstr/xerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryDetailsRemoverInterceptor is a gRPC server unary interceptor that removes sensitive details from errors if
// they are marked as hidden.
//
// This interceptor should be used in gRPC servers that return errors to external clients that are not trusted.
// For example, if the server is a public API. If the server is an internal service that is
// only called by other internal services, then it is recommended that it is not used.
func UnaryDetailsRemoverInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Call the handler
	resp, err := handler(ctx, req)
	if err != nil {
		var e *xerror.Error
		if !errors.As(err, &e) || !e.IsDetailsHidden() {
			return resp, err
		}
		_ = e.RemoveSensitiveDetails()
		return resp, err
	}
	return resp, nil
}

// XErrorFrom is a convenience function that creates a new Error from a gRPC error.
//
// Ex.
//
//	err := othersystempb.SomeMethod(ctx, req)
//	if err != nil {
//	  return grpc.XErrorFrom(err).AddVar("requested_id", req.Id)
//	}
func XErrorFrom(err error) *xerror.Error {
	return new(xerror.Error).SetStatus(status.Convert(err))
}
