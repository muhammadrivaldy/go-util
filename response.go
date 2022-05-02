package goutil

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseOK is a function for send result to client
func ResponseOK(c *gin.Context, obj interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    obj,
	})
}

// ResponseError is a function for send result error to client
func ResponseError(c *gin.Context, code int, msg error, obj interface{}) {
	message := strings.Split(msg.Error(), " ")
	if len(message) > 0 {
		message[0] = strings.ToTitle(message[0])
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    obj,
	})
}
