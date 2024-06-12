package gen

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

var errNotFound = errors.New("something not found")

const ctxKeyRequestID = "x_requestID"

/* -------------------------------------------------------------------------- */
/*                            Custom Error methods                            */
/* -------------------------------------------------------------------------- */

// getCustomStatusString maps HTTP status codes to custom status strings
// func getCustomStatusString(statusCode int) string {
// 	switch statusCode {
// 	case http.StatusOK:
// 		return "OK"
// 	case http.StatusBadRequest:
// 		return "INVALID_ARGUMENT"
// 	case http.StatusUnauthorized:
// 		return "UNAUTHORIZED"
// 	case http.StatusForbidden:
// 		return "PERMISSION_DENIED"
// 	case http.StatusNotFound:
// 		return "NOT_FOUND"
// 	case http.StatusConflict:
// 		return "ALREADY_EXISTS"
// 	case http.StatusInternalServerError:
// 		return "INTERNAL_ERROR"
// 	default:
// 		return "UNKNOWN_STATUS"
// 	}
// }

func (e *Error) Error() string {
	// TODO
	return ""
}

//nolint:nonamedreturns
func (e *Error) findBadRequestDetails() (i int, details BadRequestDetails, err error) {
	if e.Details == nil {
		return 0, BadRequestDetails{}, errNotFound
	}
	var item Error_Details_Item
	for i, item = range *e.Details {
		details, err = item.AsBadRequestDetails()
		if err != nil {
			continue
		}
		if details.FieldViolations == nil {
			continue
		}
		return i, details, nil
	}
	return 0, BadRequestDetails{}, errNotFound
}

//nolint:nonamedreturns
func (e *Error) findResourceInfoDetails() (i int, details ResourceInfoDetails, err error) {
	if e.Details == nil {
		return 0, ResourceInfoDetails{}, nil
	}

	var item Error_Details_Item
	for i, item = range *e.Details {
		details, err = item.AsResourceInfoDetails()
		if err != nil {
			continue
		}
		if details.ResourceInfos == nil {
			continue
		}
		return i, details, nil
	}
	return 0, ResourceInfoDetails{}, errNotFound
}

func (e *Error) findDebugInfoDetailsItem() *Error_Details_Item {
	if e.Details == nil {
		return nil
	}
	for _, item := range *e.Details {
		detail, err := item.AsDebugInfoDetails()
		if err != nil {
			continue
		}
		if len(detail.Detail) == 0 && len(*detail.StackEntries) == 0 {
			continue
		}
		return &item
	}
	return nil
}

func (e *Error) containsRequestInfoDetails() bool {
	if e.Details == nil {
		return false
	}
	for _, item := range *e.Details {
		detail, err := item.AsRequestInfoDetails()
		if err != nil || len(detail.RequestId) == 0 {
			continue
		}
		return true
	}
	return false
}

// SetRequestInfoDetails sets the request ID in the error details. If the error details already contain a request ID,
// this method is a no-op.
func (e *Error) SetRequestInfoDetails(ctx context.Context) {
	reqID, exists := ctx.Value(ctxKeyRequestID).(string)
	if !exists {
		return
	}

	if e.containsRequestInfoDetails() {
		return
	}

	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}

	detail := new(Error_Details_Item) //nolint:nosnakecase
	_ = detail.FromRequestInfoDetails(RequestInfoDetails{RequestId: reqID, Type: TypeGoogleapisComgoogleRpcRequestInfo})
	appendedDetails := append(*e.Details, *detail)
	e.Details = &appendedDetails
}

// AddBadRequestDetail adds a bad request detail to the error details. If the error details already contain a bad
// request detail, the new field violation is appended to the existing ones.
func (e *Error) AddBadRequestDetail(field, desc string) error {
	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}

	i, existingItem, err := e.findBadRequestDetails()
	if !errors.Is(err, errNotFound) {
		existingItem.FieldViolations = append(
			existingItem.FieldViolations, FieldViolation{Description: desc, Field: field},
		)
		return (*e.Details)[i].FromBadRequestDetails(existingItem)
	}
	detail := BadRequestDetails{
		Type:            TypeGoogleapisComgoogleRpcBadRequest,
		FieldViolations: []FieldViolation{{Description: desc, Field: field}},
	}
	detailItem := new(Error_Details_Item) //nolint:nosnakecase
	if err := detailItem.FromBadRequestDetails(detail); err != nil {
		return err
	}
	appendedDetails := append(*e.Details, *detailItem)
	e.Details = &appendedDetails
	return nil
}

// AddBadRequestDetails adds a list of bad request details to the error details. If the error details already contain
// a bad request detail, the new details are appended to the existing one.
func (e *Error) AddBadRequestDetails(violations []FieldViolation) []error {
	var errs []error
	for _, v := range violations {
		if err := e.AddBadRequestDetail(v.Field, v.Description); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (e *Error) SetPreconditionViolations(violations []PreconditionViolation) {
	detail := new(Error_Details_Item) //nolint:nosnakecase
	_ = detail.FromPreconditionFailureDetails(PreconditionFailureDetails{
		Type:       TypeGoogleapisComgoogleRpcPreconditionFailure,
		Violations: violations,
	})

	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}
	appendedDetails := append(*e.Details, *detail)
	e.Details = &appendedDetails
}

func (e *Error) SetErrorInfoDetails(domain, reason string, metadata map[string]string) {
	detail := new(Error_Details_Item) //nolint:nosnakecase
	_ = detail.FromErrorInfoDetails(ErrorInfoDetails{
		Domain:   domain,
		Reason:   reason,
		Metadata: metadata,
		Type:     TypeGoogleapisComgoogleRpcErrorInfo,
	})

	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}
	appendedDetails := append(*e.Details, *detail)
	e.Details = &appendedDetails
}

// AddResourceInfoDetails adds a list of resource info details to the error details. If the error details already
// contain a resource info detail, the new details are appended to the existing one.
func (e *Error) AddResourceInfoDetails(resourceInfos []ResourceInfo) []error {
	var errs []error
	for _, info := range resourceInfos {
		if err := e.AddResourceInfoDetail(info.Description, info.ResourceName, info.ResourceType); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// AddResourceInfoDetail adds a resource info detail to the error details. If the error details already contain a
// resource info detail, the new ResourceInfo is appended to the existing ones.
func (e *Error) AddResourceInfoDetail(desc, resourceName, resourceType string) error {
	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}

	i, existingItem, err := e.findResourceInfoDetails()
	if !errors.Is(err, errNotFound) {
		existingItem.ResourceInfos = append(
			existingItem.ResourceInfos, ResourceInfo{Description: desc, ResourceName: resourceName, ResourceType: resourceType},
		)
		return (*e.Details)[i].FromResourceInfoDetails(existingItem)
	}
	detail := ResourceInfoDetails{
		Type: TypeGoogleapisComgoogleRpcResourceInfo,
		ResourceInfos: []ResourceInfo{{
			Description:  desc,
			ResourceName: resourceName,
			ResourceType: resourceType,
		}},
	}
	detailItem := new(Error_Details_Item) //nolint:nosnakecase
	if err := detailItem.FromResourceInfoDetails(detail); err != nil {
		return err
	}
	appendedDetails := append(*e.Details, *detailItem)
	e.Details = &appendedDetails
	return nil
}

// SetDebugInfoDetail sets debug info detail to the error details.
func (e *Error) SetDebugInfoDetail(detail string) {
	if detail == "" {
		return
	}

	if e.Details == nil {
		var panicPrevention []Error_Details_Item //nolint:nosnakecase
		e.Details = &panicPrevention
	}

	d := DebugInfoDetails{
		Type:   TypeGoogleapisComgoogleRpcDebugInfo,
		Detail: detail,
	}
	existingItem := e.findDebugInfoDetailsItem()
	if existingItem != nil {
		return
	}
	detailItem := new(Error_Details_Item) //nolint:nosnakecase
	_ = detailItem.FromDebugInfoDetails(d)
	appendedDetails := append(*e.Details, *detailItem)
	e.Details = &appendedDetails
}

/* -------------------------------------------------------------------------- */
/*                             Error constructors                             */
/* -------------------------------------------------------------------------- */

func NewUnknownError(details string) *Error {
	e := &Error{
		Code:    http.StatusInternalServerError,
		Message: "something unknown happened",
		Status:  codes.Unknown.String(),
	}
	e.SetDebugInfoDetail(details)
	return e
}

func NewDeadlineExceeded() *Error {
	return &Error{
		Code:    http.StatusRequestTimeout,
		Message: "context deadline exceeded",
		Status:  codes.DeadlineExceeded.String(),
	}
}

func NewInternalError(ctx context.Context) *Error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return NewDeadlineExceeded()
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return NewCanceledError()
	}
	e := &Error{
		Code:    http.StatusInternalServerError,
		Message: "something bad happened",
		Status:  codes.Internal.String(),
	}
	e.SetRequestInfoDetails(ctx)
	return e
}

func NewInvalidArgumentError(field, desc string) *Error {
	e := &Error{
		Code:    http.StatusBadRequest,
		Message: "one or more request arguments were invalid",
		Status:  codes.InvalidArgument.String(),
	}
	if err := e.AddBadRequestDetail(field, desc); err != nil {
		return NewUnknownError(err.Error())
	}
	return e
}

func NewInvalidArgumentErrors(violations []FieldViolation) *Error {
	e := &Error{
		Code:    http.StatusBadRequest,
		Message: "one or more request arguments were invalid",
		Status:  codes.InvalidArgument.String(),
	}
	e.AddBadRequestDetails(violations)
	return e
}

func NewPreconditionFailure(description, subject, typ string) *Error {
	e := &Error{
		Code:    http.StatusBadRequest,
		Message: description,
		Status:  codes.FailedPrecondition.String(),
	}
	e.SetPreconditionViolations([]PreconditionViolation{{
		Description: description, Subject: subject, Type: typ,
	}})
	return e
}

func NewModelBindingError(err error) *Error {
	// items := []Error_Details_Item{{
	// 	union: []byte(err.Error()),
	// }}
	e := &Error{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Status:  codes.InvalidArgument.String(),
		// Details: &items, //TODO(mahmood): due to using oneof Details contains a private property called union which results in json marshalling error
	}
	return e
}

func NewUnauthenticatedError(err error, domain string) *Error {
	e := &Error{
		Code:    http.StatusUnauthorized,
		Message: "request could not be authenticated",
		Status:  codes.Unauthenticated.String(),
	}
	e.SetErrorInfoDetails(domain, "UNAUTHENTICATED", map[string]string{
		"err": err.Error(),
	})
	return e
}

func NewNotFoundError(desc, rscName, rscType string) *Error {
	e := &Error{
		Code:    http.StatusNotFound,
		Message: "requested resource not found",
		Status:  codes.NotFound.String(),
	}
	e.AddResourceInfoDetails([]ResourceInfo{{
		Description: desc, ResourceName: rscName, ResourceType: rscType,
	}})
	return e
}

func NewNotFoundErrors(infos []ResourceInfo) *Error {
	e := &Error{
		Code:    http.StatusNotFound,
		Message: "requested resource not found",
		Status:  codes.NotFound.String(),
	}
	e.AddResourceInfoDetails(infos)
	return e
}

func NewPermissionDeniedError(domain, reason string, metadata map[string]string) *Error { // nolint:unparam
	e := &Error{
		Code:    http.StatusForbidden,
		Message: "permission denied",
		Status:  codes.Unauthenticated.String(),
	}
	e.SetErrorInfoDetails(domain, "PERMISSION_DENIED", metadata)
	return e
}

func NewCanceledError() *Error {
	e := &Error{
		// Sadly there's no defined HTTP status code 499, but according to the docs that's what we're supposed to
		// return when a request is cancelled: See: https://cloud.google.com/apis/design/errors#error_payloads
		// disable lint rule magic
		Code:    499, //nolint:mnd
		Message: "request canceled by the client",
		Status:  codes.Canceled.String(),
	}
	return e
}
