package grpc

// type Error struct {
// 	xerror.Root
// 	status.Status
// }

// func (e *Error) Error() string {
// 	return e.Status.Err().Error()
// }

// // WithRoot sets the Error's Root. If called multiple times, each call will overwrite the previous Root.
// // It is useful when a root error is created in a part of the code which is shared between gRPC and HTTP handlers.
// // It allows the root error to be added to the gRPC error, so they can be passed up the call chain together.
// func (e *Error) WithRoot(r *xerror.Root) *Error {
// 	if r == nil {
// 		return e
// 	}
// 	e.Root = *r
// 	return e
// }
