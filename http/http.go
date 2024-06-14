package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tobbstr/xerror"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
)

// ErrorResponse represents the structure of the error response
type errorResponse struct {
	Error errorDetails `json:"error"`
}

// ErrorDetails represents the structure of the error details
type errorDetails struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Status  string            `json:"status"`
	Details []json.RawMessage `json:"details,omitempty"`
}

// RespondFailed returns a failed response to the client. It expects err to be of type *xerror.Error.
// If so, the returned error model is the Google Cloud APIs error model as declared in: https://google.aip.dev/193#error-response
//
// Otherwise, the response is a generic 500 Internal Server Error.
func RespondFailed(w http.ResponseWriter, err error) {
	var xerr *xerror.Error
	if !errors.As(err, &xerr) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("non-xerror received"))
		return
	}

	if xerr.IsDetailsHidden() {
		_ = xerr.RemoveSensitiveDetails()
	}

	writeError(w, xerr.StatusProto(), xerr.StatusCode(), xerr.StatusMessage())
}

func writeError(w http.ResponseWriter, st *spb.Status, code codes.Code, message string) {
	rawJSONDetails := make([]json.RawMessage, len(st.Details))
	for i, detail := range st.Details {
		b, err := protojson.Marshal(detail)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("proto marshalling detail"))
			return
		}
		rawJSONDetails[i] = b
	}

	resp := errorResponse{
		Error: errorDetails{
			Code:    int(code),
			Message: message,
			Status:  upperSnakeCaseFrom(code.String()),
			Details: rawJSONDetails,
		},
	}

	b, err := json.Marshal(&resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed to marshal error"))
		return
	}
	w.WriteHeader(runtime.HTTPStatusFromCode(code))
	_, _ = w.Write(b)
}
