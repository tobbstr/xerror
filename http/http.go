package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/tobbstr/xerror"
	"github.com/tobbstr/xerror/http/internal/gen"
)

// Alias is used to avoid name collision with the Error method on the Error type.
type httpError = gen.Error

type Error struct {
	httpError
	xerror.Root
}

func (e *Error) Error() string {
	return e.httpError.Error()
}

func RespondFailed(w http.ResponseWriter, err error) {
	var supportedError *Error
	if !errors.As(err, &supportedError) {
		w.WriteHeader(http.StatusInternalServerError)
		writeError(w, &NewUnknownError(fmt.Sprintf("want error of type Error, got %T", err)).httpError)
		return
	}

	writeError(w, &supportedError.httpError)
}

func writeError(w http.ResponseWriter, err *httpError) {
	w.WriteHeader(err.Code)
	if err := json.NewEncoder(w).Encode(err); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to write error response"))
	}
}
