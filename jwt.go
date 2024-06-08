package goutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MainToken = "main-token"
const RefreshToken = "refresh-token"

type RequestCreateJWT struct {
	SignMethod jwt.SigningMethod
	Key        string
	Data       jwt.MapClaims
}

// CreateJWT is a function for generate token
func CreateJWT(req RequestCreateJWT) (token string, err error) {

	if _, ok := req.Data["jti"]; !ok {
		req.Data["jti"] = uuid.New().String()
	}

	// create jwt token
	t := jwt.New(req.SignMethod)
	t.Claims = req.Data
	token, err = t.SignedString([]byte(req.Key))
	if err != nil {
		return
	}

	return

}

type KeyContext string

// ParseJWT is a function for parse of token string
func ParseJWT(key string, signMethod jwt.SigningMethod, attributesJWT []string) func(c *gin.Context) {
	return func(c *gin.Context) {

		errToken := errors.New("token is not valid")

		// get value authorization from header
		var authorization = c.GetHeader("authorization")
		if ok := strings.Contains(authorization, "Bearer "); !ok {
			ResponseError(c, http.StatusUnauthorized, errToken, nil)
			c.Abort()
			return
		}

		// split value without bearer
		authorization = strings.Split(authorization, "Bearer ")[1]

		// parse token
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			} else if method != signMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(key), nil
		})

		// handle error
		if err != nil {
			ResponseError(c, http.StatusUnauthorized, errToken, nil)
			c.Abort()
			return
		}

		// claim token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ResponseError(c, http.StatusUnauthorized, errToken, nil)
			c.Abort()
			return
		}

		// set value of token to context
		ctx := ParseContext(c)
		ctx = context.WithValue(ctx, KeyToken, authorization)
		for _, attribute := range attributesJWT {
			value, ok := claims[attribute]
			if ok {
				ctx = context.WithValue(ctx, KeyContext(attribute), value)
			}
		}

		// set up context.Context to gin.Context
		c.Set("context", ctx)

		// next handler
		c.Next()

	}
}
