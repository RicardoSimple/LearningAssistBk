package basic

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ar-app-api/handler/aerrors"
)

type Resp struct {
	Code    int    `json:"code"` // 0成功,1请求参数错误,2权限错误,-1服务异常
	Message string `json:"msg"`  // 错误信息
	Data    any    `json:"data"` // 数据
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Resp{Code: 0, Data: data})
}

func SuccessWithFailureMessage(c *gin.Context, msg string, code int) {
	c.JSON(http.StatusOK, Resp{Code: code, Data: nil, Message: msg})
}

type FailureType int

func RequestFailure(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Resp{Code: aerrors.NormalError, Message: msg})
	c.Abort()
}

func RequestFailureWithCode(c *gin.Context, msg string, code int) {
	c.JSON(http.StatusBadRequest, Resp{Code: code, Message: msg})
	c.Abort()
}

func RequestParamsFailure(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Resp{Code: aerrors.ParamsError, Message: "参数错误"})
	c.Abort()
}

func RequestFailureWithError(c *gin.Context, err error, code int) {
	c.JSON(http.StatusBadRequest, Resp{Code: code, Message: err.Error()})
	c.Abort()
}

func AuthFailure(c *gin.Context) {
	c.AbortWithStatus(aerrors.NoAuth)
}

func PanicFailure(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Resp{Code: aerrors.PanicError, Message: msg})
	c.Abort()
}

func PanicFailureWithError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Resp{Code: aerrors.PanicError, Message: err.Error()})
	c.Abort()
}
