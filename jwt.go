package goutil

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JWT is a object
type JWT struct {
	UserID     int64
	Name       string
	Phone      string
	Email      string
	GroupID    int
	ExpToken   time.Time
	ExpRefresh time.Time
	Jti        string
}

// CreateJWT is a function for generate token & refresh token
func CreateJWT(req JWT, signMethod jwt.SigningMethod, key string) (token, refresh string, err error) {

	if req.Jti == "" {
		req.Jti = uuid.New().String()
	}

	// create jwt token
	t := jwt.New(signMethod)
	tClaims := t.Claims.(jwt.MapClaims)
	tClaims["user_id"] = req.UserID
	tClaims["name"] = req.Name
	tClaims["phone"] = req.Phone
	tClaims["email"] = req.Email
	tClaims["group_id"] = req.GroupID
	tClaims["exp"] = req.ExpToken.Unix()
	tClaims["jti"] = req.Jti
	token, err = t.SignedString([]byte(key))
	if err != nil {
		return
	}

	// create refresh jwt token
	r := jwt.New(signMethod)
	rClaims := r.Claims.(jwt.MapClaims)
	rClaims["user_id"] = req.UserID
	rClaims["exp"] = req.ExpRefresh.Unix()
	refresh, err = r.SignedString([]byte(key))
	if err != nil {
		return
	}

	return

}

// ParseJWT is a function for parse of token string
func ParseJWT(key string, signMethod jwt.SigningMethod) func(c *gin.Context) {
	return func(c *gin.Context) {

		// get value authorization from header
		var authorization = c.GetHeader("authorization")
		if ok := strings.Contains(authorization, "Bearer "); !ok {
			ResponseError(c, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)), nil)
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
			ResponseError(c, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)), nil)
			c.Abort()
			return
		}

		// claim token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ResponseError(c, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)), nil)
			c.Abort()
			return
		}

		// set value of token to gin.context
		ctx := ParseContext(c)
		ctx = context.WithValue(ctx, KeyUserID, claims["user_id"])
		ctx = context.WithValue(ctx, KeyFullname, claims["name"])
		ctx = context.WithValue(ctx, KeyPhone, claims["phone"])
		ctx = context.WithValue(ctx, KeyEmail, claims["email"])
		ctx = context.WithValue(ctx, KeyGroupID, claims["group_id"])
		ctx = context.WithValue(ctx, KeyExp, claims["exp"])
		ctx = context.WithValue(ctx, KeyToken, authorization)

		// set up context.Context to gin.Context
		c.Set("context", ctx)

		// next handler
		c.Next()

	}
}
