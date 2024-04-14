package goutil

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Key string

const (
	KeyRequestID Key = "request_id"
	KeyMethod    Key = "method"
	KeyEndpoint  Key = "endpoint"
	KeyToken     Key = "token"
)

func (k Key) String() string {
	return string(k)
}

type Context struct {
	RequestID string
	Method    string
	Endpoint  string
}

// SetContext is a function for set value on context
func SetContext(c *gin.Context) {

	// declare variable context.Context
	ctx := c.Request.Context()

	// set up value (request id, method & endpoint) to context
	ctx = context.WithValue(ctx, KeyRequestID, requestid.Get(c))
	ctx = context.WithValue(ctx, KeyMethod, c.Request.Method)
	ctx = context.WithValue(ctx, KeyEndpoint, c.Request.URL.RequestURI())

	// set up context.Context to gin.Context
	c.Set("context", ctx)

}

// ParseContext is a function for parsing gin.Context to context.Context
func ParseContext(c *gin.Context) context.Context {

	// get context
	val, _ := c.Get("context")

	// send result context.Context value
	return val.(context.Context)

}

// getContext is a function for get value from context
func getContext(ctx context.Context) (c Context) {

	// check context value
	if ctx != nil {

		if ctx.Value(KeyRequestID) != nil {
			c.RequestID = ctx.Value(KeyRequestID).(string)
		}

		if ctx.Value(KeyMethod) != nil {
			c.Method = ctx.Value(KeyMethod).(string)
		}

		if ctx.Value(KeyEndpoint) != nil {
			c.Endpoint = ctx.Value(KeyEndpoint).(string)
		}

	}

	return
}
