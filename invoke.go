// Package invoke provides Lambda sync and async invocation helpers.
//
// The Sync() and Async() varients utilize the default Lambda client,
// while the InvokeSync() and InvokeAsync() variants may be passed
// a client in order to specify the region etc.
//
// All functions invoke DefaultAlias ("current").
package invoke

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/pkg/errors"
)

// DefaultClient is the default Lambda client.
var DefaultClient = lambda.New(session.New(aws.NewConfig()))

// DefaultAlias is the alias for function invocations.
var DefaultAlias = aws.String("current")

// Lambda interface.
type Lambda interface {
	Invoke(*lambda.InvokeInput) (*lambda.InvokeOutput, error)
}

// Error represents a Lambda error.
type Error struct {
	// Message is the error message returned from Lambda.
	Message string `json:"errorMessage"`

	// Handled specifies if the error was controlled or not.
	// For example a timeout is unhandled, while an error returned from
	// the function is handled.
	Handled bool
}

// Error message.
func (e *Error) Error() string {
	if e.Handled {
		return fmt.Sprintf("handled: %s", e.Message)
	} else {
		return fmt.Sprintf("unhandled: %s", e.Message)
	}
}

// Sync invokes function `name` synchronously with the default client.
func Sync(name string, in, out interface{}) error {
	return InvokeSync(DefaultClient, name, in, out)
}

// Async invokes function `name` asynchronously with the default client.
func Async(name string, in interface{}) error {
	return InvokeAsync(DefaultClient, name, in)
}

// InvokeSync invokes function `name` synchronously with the given `client`.
func InvokeSync(client Lambda, name string, in, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshalling input")
	}

	res, err := client.Invoke(&lambda.InvokeInput{
		FunctionName: &name,
		Qualifier:    DefaultAlias,
		Payload:      b,
	})

	if err != nil {
		return errors.Wrap(err, "invoking function")
	}

	if res.FunctionError != nil {
		err := &Error{
			Handled: *res.FunctionError == "Handled",
		}

		if e := json.Unmarshal(res.Payload, &err); e != nil {
			return errors.Wrap(e, "unmarshalling error response")
		}

		return err
	}

	if err := json.Unmarshal(res.Payload, &out); err != nil {
		return errors.Wrap(err, "unmarshalling response")
	}

	return nil
}

// InvokeAsync invokes function `name` asynchronously with the given `client`.
func InvokeAsync(client Lambda, name string, in interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return errors.Wrap(err, "marshalling input")
	}

	_, err = client.Invoke(&lambda.InvokeInput{
		FunctionName:   &name,
		InvocationType: aws.String("Event"),
		Qualifier:      DefaultAlias,
		Payload:        b,
	})

	if err != nil {
		return errors.Wrap(err, "invoking function")
	}

	return nil
}
