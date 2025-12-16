package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess      = 200
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeTooMany      = 429
	CodeServerError  = 500
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
	})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, CodeBadRequest, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, CodeUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, CodeForbidden, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, CodeNotFound, msg)
}

func TooManyRequests(c *gin.Context, msg string) {
	Error(c, http.StatusTooManyRequests, CodeTooMany, msg)
}

func ServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, CodeServerError, msg)
}
