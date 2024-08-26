package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kawasoki/gowork/validator_trans"
	"net/http"
)

var (
	ErrParams = errors.New("invalid params")
)

func HttpSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{Code: http.StatusOK, Msg: "success", Data: data})
}
func HttpFail(c *gin.Context, msgS ...string) {
	resp := &Response{Code: http.StatusBadRequest, Msg: "error"}
	if len(msgS) == 1 {
		resp.Msg = msgS[0]
	}
	c.JSON(http.StatusOK, resp)
}

// 兼容以前的
func HttpError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: -1,
		Msg:  msg,
	})
}
func HttpErrorCode(c *gin.Context, msg string, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// 定义一个空接口，用于标记空结构体
type EmptyMarker interface {
	IsEmpty()
}

// 定义一个空结构体，并实现 EmptyMarker 接口
type Empty struct{}

func (e *Empty) IsEmpty() {}

func HttpEmpty(c *gin.Context, msgS ...string) {
	resp := &Response{Code: http.StatusNoContent, Msg: "no content"}
	if len(msgS) == 1 {
		resp.Msg = msgS[0]
	}
	c.JSON(http.StatusOK, resp)
}
func HttpForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusOK, &Response{Code: http.StatusForbidden, Msg: message})
}
func HttpUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusOK, &Response{Code: http.StatusUnauthorized, Msg: message})
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Must[T any](result T, err error) T {
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (r *Response) handleValidateErr(err error) *Response {
	r.Code = http.StatusBadRequest
	r.Msg = validator_trans.Error(err) // 翻译错误信息
	return r
}
