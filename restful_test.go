package goutil_test

import (
	"net/http"
	"testing"

	goutil "github.com/muhammadrivaldy/go-util"
	"github.com/stretchr/testify/assert"
)

// go test -v -run=^TestRequestBasicAuth$
func TestRequestBasicAuth(t *testing.T) {
	restful := goutil.NewRESTful("http://localhost:8080", 2)
	statusCode, _ := restful.RequestBasicAuth(goutil.BasicAuthPayload{
		Path:   "/api/v1/security/users/customer/login",
		Method: http.MethodGet,
		Payload: goutil.BasicAuth{
			Username: "haha",
			Password: "haha",
		}})

	assert.Equal(t, 500, statusCode)
}
