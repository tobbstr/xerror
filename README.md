# xerror

xerror is a powerful error handling library designed to enhance error management in applications that utilize both HTTP and gRPC APIs. With xerror, you can easily define and wrap errors with additional context, add runtime states for logging purposes, and seamlessly handle errors in both gRPC and HTTP environments. This library also enables you to hide sensitive data in error responses and provides flexible log level configuration.

By using xerror, you can streamline your error handling process and ensure consistent and compliant error responses. Whether you're building a microservice or a complex distributed system, xerror offers a comprehensive set of features to simplify error management and improve the reliability of your applications.

Don't let errors hinder the reliability and stability of your applications. Try xerror today and experience a new level of error handling sophistication.

The error model that serves as the foundation for xerrors is the Google Cloud APIs error model, which is documented in detail in the [Google Cloud APIs error model documentation](https://google.aip.dev/193). This error model provides a robust and standardized approach to handling errors in applications. This is referred to as the `google.rpc.status` error model in this text.

## What's Included

- Simplifies error management in applications that utilize both HTTP and gRPC APIs.
- Provides additional context and runtime states for error handling and logging purposes.
- Enables hiding sensitive data in error responses.
- Offers flexible log level configuration.
- Ensures consistent and compliant error responses.
- Streamlines error handling process.
- Improves the reliability of applications.
- Provides a comprehensive set of features for error management.
- Aligns with the `google.rpc.status` error model.
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

## XError Properties

When working with xerror, you can take advantage of the following properties:

1. **Error Model**: xerror utilizes the [google.rpc.status](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto) model to effectively return consistent errors from your application to API callers.
2. **Debug Info and Error Info**: You have the option to remove debug info and error info details from the error response sent to API callers.
3. **Runtime State**: Capture and log relevant information by storing key-value pairs in the runtime state. This allows for better error analysis and troubleshooting.
4. **Log Level**: Set the log level for your application, ensuring that errors are logged at the appropriate severity level.
5. **Retries**: Utilize the xerror API to easily determine if an error is retryable.

By leveraging these properties, you can enhance your error handling process and improve the overall reliability of your application.
Remember to always consider the specific needs and requirements of your application when utilizing xerror.


## Errors Originating Within Your System

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
    return xerror.NewInvalidArgument("age", "This operation is reserved for people 65+")
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

## Errors Originating From External gRPC APIs

When working with errors returned by gRPC APIs, the `xgrpc` package provides a convenient function called `ErrorFrom()`. This function allows you to keep the error status from the external system without mapping it into a specific `xerror` such as when
constructing an xerror using the [constructors](#error-constructors) in the `xerror` package. Here's an example:

```go
// Example: Keeping the returned error status from the external call
resp, err := grpcClient.DoThat()
if err != nil {
    return xgrpc.ErrorFrom(err) // This xerror retains the error status
}
```

By using `ErrorFrom()`, you can handle errors from gRPC APIs effectively and maintain the original error status. Give it a try in your application!

## Advanced Error Handling

In some cases, simply wrapping a returned error in an `xerror` may not be sufficient. You may need to inspect the error and handle different error types differently. For example, let's say your service needs to order more pencils when it runs out of stock. If the order fails due to being out of stock, your service should make another call to restock the pencils.

The `google.rpc.status` error model fully supports this use case. It leverages the `ErrorInfo` detail, which is included in the error status. The `ErrorInfo` detail serves the purpose of uniquely identifying errors, allowing for seamless propagation across multiple service hops. This enables edge services to effectively inspect and take appropriate action based on downstream errors. Rest assured, with the `google.rpc.status` error model, your error handling process will be robust, reliable, and efficient.

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

The `reason` field (`error.details[0].reason` in this example) is designed to be examined for domain-specific errors. It is scoped to the domain and should be combined with the domain to distinguish between different services that share the same reason value.

Thankfully, the `xerror` library offers a convenient method for checking domain-specific errors, making it easier to handle and manage errors in your application.

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
        xerr := xgrpc.ErrorFrom(err)
        if xerr.IsDomainError(order.Domain, order.ReasonOutOfStock) { 
                // handle the case when pencils are out of stock
                restockPencils()
        }
}
```

In the example above, the order service exports its domain (e.g., "order.greatpencils.com") and domain-specific reasons (as enums) for the check. You can refer to [Google's enum definitions](https://github.com/googleapis/googleapis/blob/master/google/api/error_reason.proto) for inspiration.

By leveraging this advanced error handling technique, you can effectively handle different errors and ensure the smooth operation of your service.

# Error Propagation Outside of Your Domain or Bounded Context

This section discusses the handling of errors when they need to be returned to callers of your service. It is important to consider the trustworthiness of the caller in such scenarios. Internal services within the same organization are often considered trusted, but if the caller is on a public network, such as the Internet, it may be necessary to exercise caution.

## Trusted Callers

For trusted callers, your service has the flexibility to choose what error information to return.

## Untrusted Callers

Untrusted callers can be categorized into two types:

1. Applications, such as mobile apps or web apps, that belong to your organization but are running on untrusted networks. In this case, only the network is untrusted.
2. Completely external code that calls your service. In this case, both the caller and the network it is running on are untrusted.

When dealing with untrusted callers, regardless of type, avoid including sensitive information like stack traces in error responses. The `xerror` library offers a convenient method to strip sensitive details from the `google.rpc.status` error model. Specifically, it removes "debug info" and "error info" details. "Debug info" can include stack traces and diagnostic information intended for developers to diagnose issues, while "error info" provides structured data about the error. Removing these details ensures your error responses are secure and do not expose unnecessary information.


```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    return xgrpc.ErrorFrom(err).
        HideDetails() // Hides sensitive details before returning the error to the caller
}
```

For untrusted callers of type (1), the error may be propagated to the caller as-is, but it is still recommended to strip it of sensitive information. For untrusted callers of type (2), it is recommended to translate the error into a generic "internal server" error without any additional information. This can be achieved using the provided constructors mentioned in the [error constructors](#error-constructors) section. See the example below.

```go
resp, err := orderClientpb.OrderPencils()
if err != nil {
    return xerror.NewInternal(err).
        HideDetails() // Hides sensitive details before returning the error to the caller
}
```

# Using xerrors in HTTP APIs

When integrating xerrors into your HTTP APIs, it's important to handle sensitive information appropriately before returning error responses to callers. The "debug info" and "error info" details may contain sensitive data that should be sanitized.

With xerror, you don't have to deal with multiple error models for your application and HTTP API. The xerror package provides a subpackage called `xhttp`, which includes a convenient function called `RespondFailed(w http.ResponseWriter, err error)`.

**TIP:** Import the `xhttp` package using dot notation (`import . "github.com/tobbstr/xerror/xhttp"`) so you don't have to write out the package name when invoking the `RespondFailed` function. This can make your code more concise and readable. 

By passing an xerror to this function, it will automatically translate it into a JSON-formatted representation of a `google.rpc.status` error model. This simplifies error handling in your HTTP APIs and ensures consistent and standardized error responses.
See the example below:

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

In addition to HTTP APIs, xerrors can also be utilized in gRPC APIs. The process involves registering an interceptor in the server, which allows for the seamless integration of xerrors in the endpoint implementations. After registering the interceptor, xerrors should be returned in endpoint implementations. The interceptor takes care of responding with a `google.rpc.status` error. This allows for seamless integration and enhances the error handling capabilities of your gRPC APIs, ensuring consistent and standardized error responses.

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

## Logging Errors in Your Application

When it comes to logging errors in your application, there are two key considerations. First, you want to ensure that all relevant details of the error are captured. Second, you need to determine the appropriate log level for the error.

Let's consider an example scenario. Suppose you are making a gRPC call to an external system, and the call fails. In this case, you want to create an error value that accurately captures what went wrong, along with any relevant details. This error value will then be propagated up the call stack until it reaches a point where it should be logged.

By effectively logging errors in your application, you can gain valuable insights into the root causes of failures and troubleshoot issues more efficiently. Additionally, logging errors at the appropriate log level ensures that you can prioritize and address them effectively.

Remember, logging errors is an essential practice for maintaining the reliability and stability of your application. So make sure to incorporate robust error logging mechanisms into your development process.

### Capturing Runtime State

Capturing the runtime state is a useful technique when you want to preserve the values of relevant variables at the time an error occurs. This can provide valuable insights into the context and help with troubleshooting and debugging.

By including the runtime state in your error handling, you can easily identify the specific conditions that led to the error. This can be especially helpful when dealing with complex scenarios or hard-to-reproduce issues.

To capture the runtime state, you can include relevant variables or data structures in the error value itself. This ensures that the information is readily available when the error is logged or inspected.

For example, let's say you have a function that performs a calculation and encounters an error. Instead of just returning the error, you can create an xerror that includes the input parameters, intermediate results, or any other relevant data. This way, when the error is logged or reported, you have all the necessary information to understand what went wrong.

```go
func PerformCalculation(input int) error {
    // Perform the calculation
    result, err := calculate(input)
    if err != nil {
        // Create an xerror value with the runtime state
        return xerror.NewInternal(err).
            AddVar("input": input).  // Adds the input to the runtime state
            AddVar("result", result) // Adds the intermediate result to the runtime state
    }
    // Continue with the rest of the code
    return nil
}
```

By capturing the runtime state in your error handling, you can enhance the effectiveness of your debugging process and improve the overall reliability of your application.

### Logging xerrors

To effectively log xerrors in your application, you can follow these steps:

1. Identify the location in your code where the error is furthest in the call stack.
2. Initialize a new xerror using `xerror.NewInternal(...)` or any other appropriate constructor or helper function.
3. Let the xerror bubble up the call stack, adding more context to it along the way using the `xerror.Wrap()` function.
4. At the top of the call stack, use a logging library like [zap](https://github.com/uber-go/zap) to log the error.
5. Extract the runtime state from the xerror using `xerr.RuntimeState()` and log it.
6. Determine the log level based on the severity of the error using `xerr.LogLevel()`.

Here's an example of how you can log an xerror using the zap library:

```go
err := function_1()
if err != nil {
    xerr := xerror.From(err) // converts the err value into an xerror
    runtimeState := xerr.RuntimeState()
    zapFields := make([]zap.ZapField, len(runtimeState))
    for i, v := range runtimeState {
        zapFields[i] = zap.Any(v.Name, v.Value)
    }

    switch xerr.LogLevel() {
    case xerror.LogLevelInfo:
        logger.Info("invoking function_1()", zapFields...)

    // Handle other log levels...

    }
}
```

By following these steps, you can ensure that all relevant details of the xerror are captured and logged appropriately, helping you troubleshoot and debug issues more efficiently.

Remember, logging errors is an essential practice for maintaining the reliability and stability of your application. Incorporate robust error logging mechanisms into your development process to gain valuable insights into the root causes of failures.

### Setting the Log Level

Setting the log level is straightforward. Here's an example:

```go
return xerror.From(err).SetLogLevel(xerror.LogLevelWarn) // This sets a warning log level
```


## Retries

When performing functions, methods, or external calls, failures can occur for various reasons. While some failures may not be worth retrying, others may be worth attempting again.

Retries can be categorized into two types: immediate retries and higher-level retries. Immediate retries involve retrying the failed call itself, while higher-level retries occur further up the call stack, potentially involving the retry of an entire transaction.

To simplify the retry process, the `google.rpc.status` error model supports these two categories. The `xerror` library
offers two convenience functions that relieve developers from the burden of remembering which error codes allow for retries,
and of which category.

The two provided methods are:

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


