package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	Success = 0
	Error   = -1
)

func Result(code int, message string, data any, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context) {
	Result(Success, "操作成功", map[string]any{}, c)
}

func OKWithMessage(message string, c *gin.Context) {
	Result(Success, message, map[string]any{}, c)
}

func OKWithData(data any, c *gin.Context) {
	Result(Success, "操作成功", data, c)
}

func OKWithDetailed(data any, message string, c *gin.Context) {
	Result(Success, message, data, c)
}

func Fail(c *gin.Context) {
	Result(Error, "操作失败", map[string]any{}, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(Error, message, map[string]any{}, c)
}

func FailWithDetailed(data any, message string, c *gin.Context) {
	Result(Error, message, data, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    Error,
		Message: message,
		Data:    nil,
	})
}

type LoginResponse struct {
	User      User   `json:"user"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}

type User struct {
	Username string `json:"username"`
}
