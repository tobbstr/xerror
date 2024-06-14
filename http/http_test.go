package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tobbstr/golden"
	"github.com/tobbstr/xerror"
)

func TestRespondFailed(t *testing.T) {
	xerror.Init("myservice.example.com")

	type args struct {
		w   http.ResponseWriter
		err error
	}
	type given struct {
		err error
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			name: "invalid argument",
			given: given{
				err: xerror.NewInvalidArgument(xerror.BadRequestViolationOptions{
					Violation: xerror.BadRequestViolation{Field: "age", Description: "must be greater than 0"},
					LogLevel:  xerror.LogLevelError,
				}),
			},
			want: want{
				code: http.StatusBadRequest,
				body: "testdata/respond_failed/invalid_arg.json",
			},
		},
		{
			name: "invalid arguments",
			given: given{
				err: xerror.NewInvalidArguments(xerror.BadRequestViolationsOptions{
					Violations: []xerror.BadRequestViolation{
						{Field: "age", Description: "must be greater than 0"},
						{Field: "name", Description: "cannot be empty"},
					},
				}),
			},
			want: want{
				code: http.StatusBadRequest,
				body: "testdata/respond_failed/invalid_args.json",
			},
		},
		{
			name: "precondition failure",
			given: given{
				err: xerror.NewPreconditionFailure(xerror.PreconditionFailureOptions{
					Violation: xerror.PreconditionViolation{
						Description: "user could not be updated because the user was changed since it was read",
						Subject:     "example.com/v1/users/123",
						Typ:         "ErrVersionMismatch",
					},
				}),
			},
			want: want{
				code: http.StatusBadRequest,
				body: "testdata/respond_failed/precondition_failure.json",
			},
		},
		{
			name: "out of range",
			given: given{
				err: xerror.NewOutOfRange(xerror.BadRequestViolationOptions{
					Violation: xerror.BadRequestViolation{Field: "age", Description: "must be between 18 and 65"},
				}),
			},
			want: want{
				code: http.StatusBadRequest,
				body: "testdata/respond_failed/out_of_range.json",
			},
		},
		{
			name: "unauthenticated",
			given: given{
				err: xerror.NewUnauthenticated(xerror.ErrorInfoOptions{
					Error:  errors.New("failed to parse JWT token"),
					Reason: "INVALID_TOKEN",
					Metadata: map[string]any{
						"token": "JWT",
						"issue": "The length of the provided token is too short.",
					},
				}),
			},
			want: want{
				code: http.StatusUnauthorized,
				body: "testdata/respond_failed/unauthenticated.json",
			},
		},
		{
			name: "permission denied",
			given: given{
				err: xerror.NewPermissionDenied(xerror.ErrorInfoOptions{
					Error:  errors.New("user does not have permission to access the resource"),
					Reason: "API_DISABLED",
					Metadata: map[string]any{
						"resource": "projects/123",
						"service":  "pubsub.googleapis.com",
					},
				}),
			},
			want: want{
				code: http.StatusForbidden,
				body: "testdata/respond_failed/permission_denied.json",
			},
		},
		{
			name: "not found (single)",
			given: given{
				err: xerror.NewNotFound(xerror.NotFoundOptions{
					ResourceInfo: xerror.ResourceInfo{
						Description:  "resource not found",
						ResourceName: "example.v1.User",
						ResourceType: "User",
					},
				}),
			},
			want: want{
				code: http.StatusNotFound,
				body: "testdata/respond_failed/not_found_single.json",
			},
		},
		{
			name: "not found (multi)",
			given: given{
				err: xerror.NewNotFoundBulk(xerror.NotFoundBulkOptions{
					ResourceInfos: []xerror.ResourceInfo{
						{
							Description:  "resource not found",
							ResourceName: "projects/12345/buckets/MyErrors",
							ResourceType: "storage.v1.Bucket",
						},
						{
							Description:  "resource not found",
							ResourceName: "projects/12345/iam/MyUser",
							ResourceType: "iam.v1.User",
						},
					},
				}),
			},
			want: want{
				code: http.StatusNotFound,
				body: "testdata/respond_failed/not_found_multi.json",
			},
		},
		{
			name: "aborted",
			given: given{
				err: xerror.NewAborted(xerror.ErrorInfoOptions{
					Error:  errors.New("optimistic concurrency control conflict: resource revision mismatch"),
					Reason: "VERSION_MISMATCH",
					Metadata: map[string]any{
						"resource": "projects/123",
						"service":  "pubsub.googleapis.com",
					},
				}),
			},
			want: want{
				code: http.StatusConflict,
				body: "testdata/respond_failed/aborted.json",
			},
		},
		{
			name: "already exists (single)",
			given: given{
				err: xerror.NewAlreadyExists(xerror.AlreadyExistsOptions{
					ResourceInfo: xerror.ResourceInfo{
						Description:  "resource already exists",
						ResourceName: "projects/12345/buckets/MyErrors",
						ResourceType: "storage.v1.Bucket",
					},
				}),
			},
			want: want{
				code: http.StatusConflict,
				body: "testdata/respond_failed/already_exists.json",
			},
		},
		{
			name: "already exists (multi)",
			given: given{
				err: xerror.NewAlreadyExistsBulk(xerror.AlreadyExistsBulkOptions{
					ResourceInfos: []xerror.ResourceInfo{
						{
							Description:  "resource already exists",
							ResourceName: "projects/12345/buckets/MyErrors",
							ResourceType: "storage.v1.Bucket",
						},
						{
							Description:  "resource already exists",
							ResourceName: "projects/12345/iam/MyUser",
							ResourceType: "iam.v1.User",
						},
					},
				}),
			},
			want: want{
				code: http.StatusConflict,
				body: "testdata/respond_failed/already_exists_multi.json",
			},
		},
		{
			name: "resource exhausted",
			given: given{
				err: xerror.NewResourceExhausted(xerror.ResourceExhaustedOptions{
					Error: errors.New("the request for this project exceeds the available quota."),
					QuotaViolation: xerror.QuotaViolation{
						Subject:     "projects/123",
						Description: "the maximum number of instances for this project has been reached",
					},
				}),
			},
			want: want{
				code: http.StatusTooManyRequests,
				body: "testdata/respond_failed/resource_exhausted.json",
			},
		},
		{
			name: "canceled",
			given: given{
				err: xerror.NewCanceled(xerror.LogLevelDebug),
			},
			want: want{
				code: 499,
				body: "testdata/respond_failed/canceled.json",
			},
		},
		{
			name: "data loss",
			given: given{
				err: xerror.NewDataLoss(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("unrecoverable data loss or corruption"),
				}),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "testdata/respond_failed/data_loss.json",
			},
		},
		{
			name: "unknown",
			given: given{
				err: xerror.NewUnknown(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("something unknown happened"),
				}),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "testdata/respond_failed/unknown.json",
			},
		},
		{
			name: "internal",
			given: given{
				err: xerror.NewInternal(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("internal server error"),
				}),
			},
			want: want{
				code: http.StatusInternalServerError,
				body: "testdata/respond_failed/internal.json",
			},
		},
		{
			name: "not implemented",
			given: given{
				err: xerror.NewNotImplemented(xerror.LogLevelInfo),
			},
			want: want{
				code: http.StatusNotImplemented,
				body: "testdata/respond_failed/not_implemented.json",
			},
		},
		{
			name: "unavailable",
			given: given{
				err: xerror.NewUnavailable(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("service is currently unavailable"),
				}),
			},
			want: want{
				code: http.StatusServiceUnavailable,
				body: "testdata/respond_failed/unavailable.json",
			},
		},
		{
			name: "deadline exceeded",
			given: given{
				err: xerror.NewDeadlineExceeded(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("request timed out"),
				}),
			},
			want: want{
				code: http.StatusGatewayTimeout,
				body: "testdata/respond_failed/deadline_exceeded.json",
			},
		},
		{
			name: "hide details",
			given: given{
				err: xerror.NewDeadlineExceeded(xerror.ErrorWithHiddenDetailsOptions{
					Error: errors.New("request timed out"),
				}).
					SetDebugInfo("this is a debug message", []string{"line 1", "line 2"}).
					SetErrorInfo("this is an error message", "this is a reason", map[string]any{"key": "value"}).
					HideDetails(),
			},
			want: want{
				code: http.StatusGatewayTimeout,
				body: "testdata/respond_failed/hide_details.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			/* ---------------------------------- Given --------------------------------- */
			respRecorder := httptest.NewRecorder()

			args := args{w: respRecorder, err: tt.given.err}

			/* ---------------------------------- When ---------------------------------- */
			RespondFailed(args.w, args.err)

			/* ---------------------------------- Then ---------------------------------- */
			res := respRecorder.Result()

			// Assert the response
			require := require.New(t)
			require.Equal(tt.want.code, res.StatusCode)
			body := readBody(t, res.Body)
			var got map[string]any
			require.NoError(json.Unmarshal(body, &got))
			golden.JSON(t, tt.want.body, got)
		})
	}
}

func readBody(t *testing.T, body io.ReadCloser) []byte {
	t.Helper()
	b, err := io.ReadAll(body)
	require.NoError(t, err, "failed to read response body")
	defer body.Close()
	return b
}
