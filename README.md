# xerror

xerror is a powerful error handling library designed to enhance error management in applications that utilize both HTTP and gRPC APIs. With xerror, you can easily define and wrap errors with additional context, add runtime states for logging purposes, and seamlessly handle errors in both gRPC and HTTP environments. This library also enables you to hide sensitive data in error responses and provides flexible log level configuration.

By using xerror, you can streamline your error handling process and ensure consistent and compliant error responses. Whether you're building a microservice or a complex distributed system, xerror offers a comprehensive set of features to simplify error management and improve the reliability of your applications.

Don't let errors hinder the reliability and stability of your applications. Try xerror today and experience a new level of error handling sophistication.

The underlying error model used is the [Google Cloud APIs error model](https://google.aip.dev/193).

## What's included

- Simplifies error management in applications that utilize both HTTP and gRPC APIs.
- Provides additional context and runtime states for error handling and logging purposes.
- Enables hiding sensitive data in error responses.
- Offers flexible log level configuration.
- Ensures consistent and compliant error responses.
- Streamlines error handling process.
- Improves the reliability of applications.
- Provides a comprehensive set of features for error management.
- Aligns with the Google Cloud APIs error model.
- Enhances error handling sophistication.
- Prevents errors from hindering reliability and stability of applications.
- Makes it easy to classify errors with a builtin error guide (usable from your code)

## Usage

To get started, simply add a call to `xerror.Init()` in your main.go file and run `go mod tidy`. Then, you can leverage
the various functionalities of xerror to handle errors effectively and efficiently.

The `xerror.Init()` call is used to initialize the `xerror` package. It sets a global value called "domain" that is used when adding an [ErrorInfo](https://github.com/googleapis/googleapis/blob/master/google/rpc/error_details.proto#L51) detail to our error. The "domain" value should represent either the name of your service or its domain name. By setting the "domain" value globally, you don't have to specify it every time you add an ErrorInfo detail to your error. This simplifies the error handling process and ensures consistency throughout your application.

An example:

```go
// main.go

import "github.com/tobbstr/xerror"

func main() {
    xerror.Init("pubsub.googleapis.com") // replace this string with your service name or service domain as in the example
}
```

Then you're all set! ✅

See the next sections for how to use it for different purposes.

## XError properties

When working with xerror, you can take advantage of the following properties:

1. **Error Model**: xerror utilizes the [google.rpc.status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto) model to effectively return consistent errors from your application to API callers.
2. **Debug Info and Error Info**: You have the option to remove debug info and error info details from the error response sent to API callers.
3. **Runtime State**: Capture and log relevant information by storing key-value pairs in the runtime state. This allows for better error analysis and troubleshooting.
4. **Log Level**: Set the log level for your application, ensuring that errors are logged at the appropriate severity level.
5. **Retries**: Utilize the xerror API to easily determine if an error is retryable.

By leveraging these properties, you can enhance your error handling process and improve the overall reliability of your application.
Remember to always consider the specific needs and requirements of your application when utilizing xerror.


## Errors originating within your system

This section provides guidance on working with errors that occur within your system. Whether it's validating arguments or checking preconditions, you are responsible for creating the root error and classifying it correctly according to the [google.rpc.code](https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto) definitions.

### Error Constructors

The xerror library offers convenient constructors to initialize errors of different types. These types align with the [google.rpc.code](https://github.com/googleapis/googleapis/blob/master/google/rpc/code.proto) definitions, including:

- INVALID_ARGUMENT
- FAILED_PRECONDITION
- OUT_OF_RANGE
- UNAUTHENTICATED
- PERMISSION_DENIED
- NOT_FOUND
- ABORTED
- ALREADY_EXISTS
- RESOURCE_EXHAUSTED
- CANCELLED
- DATA_LOSS
- UNKNOWN
- INTERNAL
- NOT_IMPLEMENTED
- UNAVAILABLE
- DEADLINE_EXCEEDED

For example, to create an error of type `INVALID_ARGUMENT`, you can use the following constructor:

```go
// validating arguments ...

if request.Age < 65 {
    return xerror.NewInvalidArgument(BadRequestOptions{
        Violation: BadRequestViolation{Field: "age", Description: "This operation is reserved for people 65+"}
    })
}
```

### Visual Overview of Error Types

To help you understand the organization of error types, a visual overview is provided. Some error types are specific to problems with the request, while others are related to server issues. Additionally, certain error types are nested within others for more specialized scenarios. Here is a simplified representation:

```mermaid
flowchart LR
subgraph "Problems with the Request"
    CANCELED
    subgraph INVALID_ARGUMENT
        OUT_OF_RANGE
        NOT_FOUND
        DATA_LOSS
    end
    PERMISSION_DENIED
    UNAUTHENTICATED
end

subgraph "Problems with the Server"
    DATA_LOSS_2["DATA_LOSS"]
    subgraph FAILED_PRECONDITION
        ABORTED
        ALREADY_EXISTS
        RESOURCE_EXHAUSTED
    end
    UNKNOWN
    INTERNAL
    NOT_IMPLEMENTED
    UNAVAILABLE
    DEADLINE_EXCEEDED
end
```

### Error Guide

Don't let error classification complexity hinder your development. With xerror's error guide
(`xerror.ErrorGuide()`), confidently identify and classify errors, making your code more consistent, robust and maintainable. The xerror library provides a built-in error guide function that simplifies error classification in your code. This error guide assists you in accurately categorizing errors, ensuring proper handling and management. By leveraging this feature, you can streamline your error handling process and improve the reliability of your application.

While the error guide is a valuable resource, you can also take advantage of the convenient constructor functions such as `xerror.NewNotFound()` once you have a clear understanding of your requirements.

Feel free to explore the error guide and constructor functions to streamline your error handling process and ensure accurate error classification. Hopefully this information helps you effectively handle errors within your system.

## Error logging

When it comes to logging, as a developer you're interested in two things. First, to log the error and all the relevant
details. Second, to do that at a desired log level.

Ex. Suppose you're making a call to an external system and that it fails. It's a gRPC call and you want to define
an error that captures what went wrong and all the relevant details. This error is meant to be returned up the call
stack until it reaches a point where it should be logged.

```go
import "github.com/tobbstr/xerror/xgrpc"

// skipped for brevity ...

resp, err := otherService.SomeGrpcMethod(&SomeGrpcMethodRequest{
    Name: name, Age: age,
}); err != nil {
    // Initialize a new xerror and add the relevant context to it. Also set its log level
    return xgrpc.ErrorFrom(err).SetLogLevel(xerror.LogLevelError).AddVar("name", name).AddVar("age", age)
}
```

### Adding and using runtime state

Is useful when you want to capture the values of relevant variables when the error happened.

```go
err := someFailingFunction()
if err != nil {
    return xerror.NewInternal(xerror.SimpleOptions{Error: err}).
        AddVar("name", name). // Stores a single variable
        // Stores multiple variables
        AddVars([]xerror.Var{
            { Name: "name", Value: name},
            { Name: "age", Value: age},
        }...)
}
```

Let's say we have the following call stack:

```
entrypoint()                                <- Since it was returned in the previous function it's now here
├─ function_1()                             <- We initialize a new xerror here by calling xerror.NewInternal(...)
   ├─ externallibrary.FailingFunction()     <- Error happened here
```

Now that the error is the furthest it can go in the call stack, it's time to log the error. In the example below
we're using the [zap library](https://github.com/uber-go/zap).

```go
// entrypoint code
{
    // skipped for brevity ...

    err := function_1()
    if err != nil {
        // Log the error
        xerr := xerror.From(err)
        runtimeState := xerr.RuntimeState()
        zapFields := make([]zap.ZapField, len(runtimeState))
        for i, v := range runtimeState {
            zapFields[i] = zap.Any(v.Name, v.Value)
        }

        // We're using the log level of the error to determine the severity of the logged message
        switch xerr.LogLevel() {
        case xerror.LogLevelInfo:
            logger.Info("invoking function_1()", zapFields...)

        // the rest of the cases are skipped for brevity ...

        }
    }
}
```

### Setting the log level

Setting the log level is quite easy as demonstrated in the example below.

```go
return xerror.From(err).SetLogLevel(xerror.LogLevelWarn) // This sets a warning log level
```

## Errors originating from external systems

It happens that your application is returned an error from an external system such as when calling a gRPC
endpoint. In general, you should inspect the error and create a desired xerror using the [constructors](#error-constructors).
However, in the case of gRPC this library contains a package called `xgrpc` which has a convenience function for the
times you do not want to map the error into a particular xerror, but rather want to keep the status from the external
system. An example of this is given below:

```go
// Example: When you want to keep the returned error status from the external call
resp, err := grpcClient.DoThat()
if err != nil {
    return xgrpc.ErrorFrom(err). // this xerror keeps the error status from err
        AddVar("variable_name_1", valueOfVariable1) // adds context to the error that can be logged easily
}
```

### Advanced error handling

Sometimes it's not enough to just wrap a returned error in an xerror directly. Let's say you need to inspect the
returned error and handle different error types differently. For example, let's say in the case you're out of stock
of pencils, your service is responsible for ordering more of them. When your service makes a call to place an order
for pencils, then it can fail and the returned error may be of a certain type to signify that the order couldn't be
fulfilled due to being out of stock, in which case your service should make another call to order restocking of pencils.

Does Google Cloud APIs error model support this usecase? Yes, it does. The way it handles it is by having an
`ErrorInfo` detail included in the error status.

```json
{
  "error": {
    "code": 8,
    "message": "The order couldn't be fulfilled. The requested item is out of stock",
    "status": "RESOURCE_EXHAUSTED",
    "details": [
      {
        "@type": "type.googleapis.com/google.rpc.ErrorInfo",
        "reason": "OUT_OF_STOCK",
        "domain": "greatpencils.com",
        "metadata": {
          "service": "order.greatpencils.com"
        }
      }
    ]
  }
}
```

Neither the "error.message" nor the "error.status" is meant for these domain-specific errors. Instead it's the
reason field (the "error.details[0].reason" in this example) that is meant to be inspected. The reason only makes
sense in combination with the domain. The reason is domain-scoped. Two different services may have the exact same
reason value, but to distinguish between them the domain must also be taken into account when inspecting.

Luckily the xerror library provides a convenience method on the xerror for checking for domain-specific errors.

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    xerr := xgrpc.ErrorFrom(err)
    // The order service has an exported domain constant and an enum for the reason. These are used to check the
    // returned error type.
    if xerr.IsDomainError(orderpb.Domain, orderpb.ReasonOutOfStock) { 
        // handle the case when pencils are out of stock
        restockPencils()
    }
}
```

NOTE! In the example above the order service (which is an external service in this example) exports its domain (ex.
"order.greatpencils.com") and its domain-specific reasons (as enums) that we use for the check. See [Google's
enum definitions](https://github.com/googleapis/googleapis/blob/master/google/api/error_reason.proto) for inspiration.

## Retries

Whenever a function/method or external call etc., is made, it may fail for any number of reasons. For some of them
retrying is futile, but for some it could be worth an attempt or two.

[Retries can be categorised into two types](https://cloud.google.com/apis/design/errors#retrying_errors). One is the retry of the immediate call that failed. The other is the
retry at a higher level in the code. A higher level means futher up the call stack, which could mean a retry of
a whole transaction.

These two categories are supported by the Google Cloud APIs error model, by inspection of the error status' code.
This library reduces the cognitive load of developers by not requiring them to remember which code means that
retries could be attempted.

Two methods are provided that map to these retry categories:

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    xerr := xgrpc.ErrorFrom(err)
    if xerr.IsDirectlyRetryable() {
        // implementation skipped for brevity
    } else if xerr.IsRetryableAtHigherLevel() {
        // implementation skipped for brevity
    }
}
```

## Error propagation outside of your domain or bounded context

This section discusses when errors are to be returned to callers of your service. When that happens the first thing
to take into consideration is whether the caller can be considered trusted. Internal services belonging to the same
organsation is often considered trusted, but if they're on public networks such as the Internet, then maybe they
shouldn't be given too much trust.

### Trusted callers

It's completely up to your service what error to return to the caller.

### Untrusted callers

Untrusted callers come in two flavours:

1. Applications such as mobile apps, web apps etc., that belong to your organisation, but they're run on untrusted
networks. In that case only the network is untrusted.
2. The other flavour is completely external code that calls your service. In that case both the caller and the
network it's running on are untrusted.

The recommendation when dealing with untrusted callers, no matter which flavour, is not to include sensitive
information in the error such as stack traces. The xerror library provides a convenience method to do this.

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    return xgrpc.ErrorFrom(err). // constructs an xerror from the err
        // marks the xerror as having sensistive details so it will be removed before returning it to the caller
        HideDetails()
}
```

When dealing with untrusted callers of type (1) then the error may be propagated to the caller "as is" as in the
example above, but it's recommended to strip it of sensitive information. For untrusted callers of type (2), it's
[recommended](https://cloud.google.com/apis/design/errors#propagating_errors) the error is translated into a generic "internal server" error without any additional information. This is easily achieved by
using the provided constructors mentioned [here](#error-constructors). See the example below.

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    // Returns an INTERNAL error no matter what the external call returned.
    return xerror.NewInternal(SimpleOptions{Error: err}).
        // marks the xerror as having sensistive details so it will be removed before returning it to the caller
        HideDetails()
}
```

## Using xerrors in HTTP APIs

Isn't it a hassle having to support multiple error models? One for the application and another for the HTTP API?
Well, with xerror you don't have to. The xerror package has a subpackage called `xhttp`, which contains a single
function `RespondFailed(w http.ResponseWriter, err error)`.

That function expects an xerror and if that is passed to it, then it gets translated into a JSON-formatted
representation of a Google Cloud APIs error model. For example, when an xerror of type "ABORTED" is passed to the
"RespondFailed()" function then the body in the HTTP response will look similar to this:

```json
{
    "error": {
        "code": 10,
        "details": [
            {
                "@type": "type.googleapis.com/google.rpc.ErrorInfo",
                "domain": "myservice.example.com",
                "metadata": {
                    "resource": "projects/123",
                    "service": "pubsub.googleapis.com"
                },
                "reason": "VERSION_MISMATCH"
            }
        ],
        "message": "optimistic concurrency control conflict: resource revision mismatch",
        "status": "ABORTED"
    }
}
```

## Using xerrors in gRPC APIs

Just like for HTTP APIs, xerrors support gRPC APIs as well. The difference is that, for gRPC, an interceptor
has to be registered in the server. Once that's done then one can simply return xerrors in the gRPC endpoint
implementations.

```go
// main.go

import (
    "google.golang.org/grpc"
    "github.com/tobbstr/xerror/xgrpc"
)

func main() {
    // Create a gRPC server instance with the interceptor
    server := grpc.NewServer(
        grpc.UnaryInterceptor(xgrpc.UnaryXErrorInterceptor),
    )
}
```
