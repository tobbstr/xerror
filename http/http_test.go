package http

// func TestRespondFailed(t *testing.T) {
// 	type args struct {
// 		err error
// 	}
// 	type given struct {
// 		args args
// 	}
// 	type want struct {
// 		code int
// 		body string
// 	}
// 	tests := []struct {
// 		name  string
// 		given given
// 		want  want
// 	}{
// 		{
// 			name: "canceled",
// 			given: given{
// 				args: args{
// 					err: NewCanceledError(),
// 				},
// 			},
// 			want: want{
// 				code: 499,
// 				body: "testdata/respond_failed/canceled.json",
// 			},
// 		},
// 		{
// 			name: "deadline exceeded",
// 			given: given{
// 				args: args{
// 					err: NewDeadlineExceeded(),
// 				},
// 			},
// 			want: want{
// 				code: 408,
// 				body: "testdata/respond_failed/deadline_exceeded.json",
// 			},
// 		},
// 		{
// 			name: "internal",
// 			given: given{
// 				args: args{
// 					err: NewInternalError(context.Background()),
// 				},
// 			},
// 			want: want{
// 				code: 500,
// 				body: "testdata/respond_failed/internal.json",
// 			},
// 		},
// 		{
// 			name: "invalid argument",
// 			given: given{
// 				args: args{
// 					err: NewInvalidArgumentError("age", "must be greater than 0"),
// 				},
// 			},
// 			want: want{
// 				code: 400,
// 				body: "testdata/respond_failed/invalid_arg.json",
// 			},
// 		},
// 		{
// 			name: "invalid arguments",
// 			given: given{
// 				args: args{
// 					err: NewInvalidArgumentErrors([]FieldViolation{
// 						{Field: "age", Description: "must be greater than 0"},
// 						{Field: "name", Description: "cannot be empty"},
// 					}),
// 				},
// 			},
// 			want: want{
// 				code: 400,
// 				body: "testdata/respond_failed/invalid_args.json",
// 			},
// 		},
// 		{
// 			name: "model binding",
// 			given: given{
// 				args: args{
// 					err: NewModelBindingError(errors.New("failed to bind model")),
// 				},
// 			},
// 			want: want{
// 				code: 400,
// 				body: "testdata/respond_failed/model_binding.json",
// 			},
// 		},
// 		{
// 			name: "not found (single)",
// 			given: given{
// 				args: args{
// 					err: NewNotFoundError(
// 						"resource not found",
// 						"User",
// 						"example.v1.User",
// 					),
// 				},
// 			},
// 			want: want{
// 				code: 404,
// 				body: "testdata/respond_failed/not_found_single.json",
// 			},
// 		},
// 		{
// 			name: "not found (multi)",
// 			given: given{
// 				args: args{
// 					err: NewNotFoundErrors([]ResourceInfo{
// 						{
// 							Description:  "resource not found",
// 							ResourceName: "projects/12345/buckets/MyErrors",
// 							ResourceType: "storage.v1.Bucket",
// 						},
// 						{
// 							Description:  "resource not found",
// 							ResourceName: "projects/12345/iam/MyUser",
// 							ResourceType: "iam.v1.User",
// 						},
// 					}),
// 				},
// 			},
// 			want: want{
// 				code: 404,
// 				body: "testdata/respond_failed/not_found_multi.json",
// 			},
// 		},
// 		{
// 			name: "permission denied",
// 			given: given{
// 				args: args{
// 					err: NewPermissionDeniedError(
// 						"googleapis.com",
// 						"API_DISABLED",
// 						map[string]string{
// 							"resource": "projects/123",
// 							"service":  "pubsub.googleapis.com",
// 						},
// 					),
// 				},
// 			},
// 			want: want{
// 				code: 403,
// 				body: "testdata/respond_failed/permission_denied.json",
// 			},
// 		},
// 		{
// 			name: "precondition failure",
// 			given: given{
// 				args: args{
// 					err: NewPreconditionFailure(
// 						"user could not be updated because the user was changed since it was read",
// 						"example.com/v1/users/123",
// 						"ErrVersionMismatch",
// 					),
// 				},
// 			},
// 			want: want{
// 				code: 400,
// 				body: "testdata/respond_failed/precondition_failure.json",
// 			},
// 		},
// 		{
// 			name: "unauthenticated",
// 			given: given{
// 				args: args{
// 					err: NewUnauthenticatedError(
// 						errors.New("failed to parse JWT token"),
// 						"api.example.com",
// 					),
// 				},
// 			},
// 			want: want{
// 				code: 401,
// 				body: "testdata/respond_failed/unauthenticated.json",
// 			},
// 		},
// 		{
// 			name: "unknown",
// 			given: given{
// 				args: args{
// 					err: NewUnknownError("more information about the error"),
// 				},
// 			},
// 			want: want{
// 				code: 500,
// 				body: "testdata/respond_failed/unknown.json",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			/* ---------------------------------- Given --------------------------------- */
// 			respRecorder := httptest.NewRecorder()

// 			/* ---------------------------------- When ---------------------------------- */
// 			RespondFailed(respRecorder, tt.given.args.err)

// 			/* ---------------------------------- Then ---------------------------------- */
// 			res := respRecorder.Result()

// 			// Assert the response
// 			require := require.New(t)
// 			require.Equal(tt.want.code, res.StatusCode)
// 			body := readBody(t, res.Body)
// 			var got map[string]any
// 			require.NoError(json.Unmarshal(body, &got))
// 			golden.JSON(t, tt.want.body, got)
// 		})
// 	}
// }

// func TestError(t *testing.T) {
// 	/* ---------------------------------- Given --------------------------------- */
// 	vars := []xerror.Var{
// 		{Name: "var1", Value: "value1"},
// 		{Name: "var2", Value: "value2"},
// 	}

// 	/* ---------------------------------- When ---------------------------------- */
// 	err := NewInternalError(context.Background()).WithError(errors.New("original error")).WithDirectRetry().
// 		WithSeverity(xerror.LogLevelError).WithVars(vars...)

// 	/* ---------------------------------- Then ---------------------------------- */
// 	require.Equal(t, err.Err, errors.New("original error"))
// 	require.Equal(t, err.DirectlyRetryable, true)
// 	require.Equal(t, err.Severity, xerror.LogLevelError)
// 	require.Equal(t, err.RuntimeState, vars, "runtime state mismatch")
// }

// func readBody(t *testing.T, body io.ReadCloser) []byte {
// 	t.Helper()
// 	b, err := io.ReadAll(body)
// 	require.NoError(t, err, "failed to read response body")
// 	defer body.Close()
// 	return b
// }
