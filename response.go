package goutil

import (
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseOK is a function for send result to client
func ResponseOK(c *gin.Context, code int, obj interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: "Success",
		Data:    obj,
	})
}

// ResponseError is a function for send result error to client
func ResponseError(c *gin.Context, code int, msg error, obj interface{}) {

	message := msg.Error()
	messages := strings.Split(msg.Error(), " ")
	if len(messages) > 0 {
		messages[0] = cases.Title(language.AmericanEnglish).String(messages[0])
		message = strings.Join(messages, " ")
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    obj,
	})
}
