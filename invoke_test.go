package invoke_test

import (
	"testing"

	"github.com/apex/invoke"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/stretchr/testify/assert"
)

type client struct {
	FunctionError *string
}

func (c *client) Invoke(in *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return &lambda.InvokeOutput{
		FunctionError: c.FunctionError,
		Payload:       in.Payload,
	}, nil
}

type input struct {
	Value string `json:"name"`
}

type output struct {
	Value string `json:"name"`
}

func TestInvokeSync(t *testing.T) {
	c := new(client)
	var out output
	err := invoke.InvokeSync(c, "test", input{"hello"}, &out)
	assert.NoError(t, err)
	assert.Equal(t, "hello", out.Value)
}

func TestInvokeSync_noInput(t *testing.T) {
	c := new(client)
	var out output
	err := invoke.InvokeSync(c, "test", nil, &out)
	assert.NoError(t, err)
	assert.Equal(t, "", out.Value)
}

func TestInvokeSync_noOutput(t *testing.T) {
	c := new(client)
	err := invoke.InvokeSync(c, "test", input{"hello"}, nil)
	assert.NoError(t, err)
}

func TestInvokeSync_noInput_noOutput(t *testing.T) {
	c := new(client)
	err := invoke.InvokeSync(c, "test", nil, nil)
	assert.NoError(t, err)
}

func TestInvokeSync_error(t *testing.T) {
	c := new(client)
	var out output
	c.FunctionError = aws.String("Unhandled")

	err := invoke.InvokeSync(c, "test", &invoke.Error{Message: "Task timed out after 5.00 seconds"}, &out)
	assert.Equal(t, "unhandled: Task timed out after 5.00 seconds", err.Error())

	e := err.(*invoke.Error)
	assert.False(t, e.Handled)
	assert.Equal(t, "Task timed out after 5.00 seconds", e.Message)
}

func TestInvokeAsync(t *testing.T) {
	c := new(client)
	err := invoke.InvokeAsync(c, "test", input{"hello"})
	assert.NoError(t, err)
}
