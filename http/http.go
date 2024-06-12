package http

// Alias is used to avoid name collision with the Error method on the Error type.
// type httpError = gen.Error

// type Error struct {
// 	httpError
// 	xerror.Root
// }

// func (e *Error) Error() string {
// 	return e.httpError.Error()
// }

// // WithRoot sets the Error's Root. If called multiple times, each call will overwrite the previous Root.
// // It is useful when a root error is created in a part of the code which is shared between gRPC and HTTP handlers.
// // It allows the root error to be added to the HTTP error, so they can be passed up the call chain together.
// func (e *Error) WithRoot(r *xerror.Root) *Error {
// 	if r == nil {
// 		return e
// 	}
// 	e.Root = *r
// 	return e
// }

// func RespondFailed(w http.ResponseWriter, err error) {
// 	var supportedError *Error
// 	if !errors.As(err, &supportedError) {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		writeError(w, &NewUnknownError(fmt.Sprintf("want error of type Error, got %T", err)).httpError)
// 		return
// 	}

// 	writeError(w, &supportedError.httpError)
// }

// func writeError(w http.ResponseWriter, err *httpError) {
// 	w.WriteHeader(err.Code)
// 	if err := json.NewEncoder(w).Encode(err); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		_, _ = w.Write([]byte("failed to write error response"))
// 	}
// }
