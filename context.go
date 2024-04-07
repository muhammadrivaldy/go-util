package goutil

import (
	"context"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Key string

const (
	KeyRequestID Key = "request_id"
	KeyMethod    Key = "method"
	KeyEndpoint  Key = "endpoint"
	KeyUserID    Key = "user_id"
	KeyName      Key = "name"
	KeyEmail     Key = "email"
	KeyUserType  Key = "user_type"
	KeyExp       Key = "exp"
	KeyToken     Key = "token"
	KeyJti       Key = "jti"
	KeyType      Key = "type"
)

func (k Key) String() string {
	return string(k)
}

type Context struct {
	RequestID string
	Method    string
	Endpoint  string
	Payload   string
	UserID    int64
	Name      string
	Phone     string
	Email     string
	UserType  int
	Exp       time.Time
	Token     string
	Jti       string
	Type      string
}

// SetContext is a function for set value on context
func SetContext(c *gin.Context) {

	// declare variable context.Context
	ctx := context.Background()

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

// GetContext is a function for get value from context
func GetContext(ctx context.Context) (c Context) {

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

		if ctx.Value(KeyUserID) != nil {
			c.UserID = int64(ctx.Value(KeyUserID).(float64))
		}

		if ctx.Value(KeyName) != nil {
			c.Name = ctx.Value(KeyName).(string)
		}

		if ctx.Value(KeyEmail) != nil {
			c.Email = ctx.Value(KeyEmail).(string)
		}

		if ctx.Value(KeyUserType) != nil {
			c.UserType = int(ctx.Value(KeyUserType).(float64))
		}

		if ctx.Value(KeyExp) != nil {
			timestamp := int64(ctx.Value(KeyExp).(float64))
			c.Exp = time.Unix(timestamp, 0)
		}

		if ctx.Value(KeyToken) != nil {
			c.Token = ctx.Value(KeyToken).(string)
		}

		if ctx.Value(KeyJti) != nil {
			c.Jti = ctx.Value(KeyJti).(string)
		}

		if ctx.Value(KeyType) != nil {
			c.Type = ctx.Value(KeyType).(string)
		}

	}

	return
}
