package xgrpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tobbstr/golden"
	"github.com/tobbstr/xerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorFrom(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *xerror.Error
	}{
		{
			name: "err is of type status.Status",
			args: args{err: status.New(codes.Canceled, "request cancelled by the client").Err()},
			want: xerror.NewCancelled(),
		},
		{
			name: "err is not of type status.Status",
			args: args{err: errors.New("some error")},
			want: xerror.NewUnknown(errors.New("some error")).SetLogLevel(xerror.LogLevelUnspecified),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* ---------------------------------- When ---------------------------------- */
			got := ErrorFrom(tt.args.err)

			/* ---------------------------------- Then ---------------------------------- */
			require.Equal(t, tt.want.StatusCode(), got.StatusCode())
			require.Equal(t, tt.want.StatusMessage(), got.StatusMessage())
			require.Equal(t, tt.want.LogLevel(), got.LogLevel())
		})
	}
}

func TestUnaryXErrorInterceptor(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     any
		info    *grpc.UnaryServerInfo
		handler grpc.UnaryHandler
	}
	type given struct {
		handler grpc.UnaryHandler
	}
	type want struct {
		value string
		err   string
	}
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "handler returns no error",
			given: given{
				handler: func(ctx context.Context, req any) (any, error) {
					return "whatever", nil
				},
			},
			want: want{
				value: "testdata/unary_xerror_interceptor/no_error.value.json",
				err:   "testdata/unary_xerror_interceptor/no_error.err.json",
			},
		},
		{ // The sensitive details should be present in the error
			name: "xerror without hidden details",
			given: given{
				handler: func(ctx context.Context, req any) (any, error) {
					return nil, xerror.NewCancelled().
						SetDebugInfo("this is a debug message", []string{"line 1", "line 2"}).
						SetErrorInfo("this is an error message", "this is a reason", map[string]any{"key": "value"})
				},
			},
			want: want{
				value: "testdata/unary_xerror_interceptor/no_hidden_details.value.json",
				err:   "testdata/unary_xerror_interceptor/no_hidden_details.err.json",
			},
		},
		{ // The sensitive details should be absent in the error
			name: "xerror with hidden details",
			given: given{
				handler: func(ctx context.Context, req any) (any, error) {
					return nil, xerror.NewCancelled().
						SetDebugInfo("this is a debug message", []string{"line 1", "line 2"}).
						SetErrorInfo("this is an error message", "this is a reason", map[string]any{"key": "value"}).
						HideDetails() // This call should hide the details
				},
			},
			want: want{
				value: "testdata/unary_xerror_interceptor/hidden_details.value.json",
				err:   "testdata/unary_xerror_interceptor/hidden_details.err.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* ---------------------------------- Given --------------------------------- */
			args := args{ctx: context.Background(), handler: tt.given.handler}

			/* ---------------------------------- When ---------------------------------- */
			got, err := UnaryXErrorInterceptor(args.ctx, args.req, args.info, args.handler)

			/* ---------------------------------- Then ---------------------------------- */
			// Assert the returned value
			golden.JSON(t, tt.want.value, got)

			// Assert the returned error
			xerr := ErrorFrom(err)
			golden.JSON(t, tt.want.err, xerr)
		})
	}
}
