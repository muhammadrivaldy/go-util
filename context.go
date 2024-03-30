package goutil

import (
	"context"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// Key ..
type Key string

const (
	// KeyRequestID ..
	KeyRequestID Key = "request_id"
	// KeyMethod ..
	KeyMethod Key = "method"
	// KeyEndpoint ..
	KeyEndpoint Key = "endpoint"
	// KeyUserID ..
	KeyUserID Key = "user_id"
	// KeyFullname ..
	KeyFullname Key = "name"
	// KeyEmail ..
	KeyEmail Key = "email"
	// KeyUserType ..
	KeyUserType Key = "user_type"
	// KeyExp ..
	KeyExp Key = "exp"
	// KeyToken ..
	KeyToken Key = "token"
	// KeyJti ..
	KeyJti Key = "jti"
)

func (k Key) String() string {
	return string(k)
}

// Context ..
type Context struct {
	RequestID string
	Method    string
	Endpoint  string
	Payload   string
	UserID    int64
	Name      string
	Phone     string
	Email     string
	GroupID   int
	Exp       time.Time
	Token     string
	Jti       string
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

		// set request id
		if ctx.Value(KeyRequestID) != nil {
			c.RequestID = ctx.Value(KeyRequestID).(string)
		}

		// set method
		if ctx.Value(KeyMethod) != nil {
			c.Method = ctx.Value(KeyMethod).(string)
		}

		// set endpoint
		if ctx.Value(KeyEndpoint) != nil {
			c.Endpoint = ctx.Value(KeyEndpoint).(string)
		}

		// set user id
		if ctx.Value(KeyUserID) != nil {
			c.UserID = int64(ctx.Value(KeyUserID).(float64))
		}

		// set fullname
		if ctx.Value(KeyFullname) != nil {
			c.Name = ctx.Value(KeyFullname).(string)
		}

		// set email
		if ctx.Value(KeyEmail) != nil {
			c.Email = ctx.Value(KeyEmail).(string)
		}

		// set group id
		if ctx.Value(KeyUserType) != nil {
			c.GroupID = int(ctx.Value(KeyUserType).(float64))
		}

		// set expired
		if ctx.Value(KeyExp) != nil {
			timestamp := int64(ctx.Value(KeyExp).(float64))
			c.Exp = time.Unix(timestamp, 0)
		}

		// set expired
		if ctx.Value(KeyToken) != nil {
			c.Token = ctx.Value(KeyToken).(string)
		}

		// set jti
		if ctx.Value(KeyJti) != nil {
			c.Jti = ctx.Value(KeyJti).(string)
		}

	}

	return
}
