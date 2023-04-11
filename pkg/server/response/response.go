package response

import "github.com/gin-gonic/gin"

const (
	MsgOk   = "操作成功"
	MsgFail = "操作失败"
	Msg200  = "操作成功"
	Msg400  = "数据错误"
	Msg401  = "身份验证错误"
	Msg403  = "无权访问"
	Msg404  = "找不到数据"
	Msg405  = "方法不支持"
	Msg500  = "服务器错误"
)

type RespJson struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Response(code int, msg string, data any, c *gin.Context) {
	var ok bool
	if code == 200 {
		ok = true
	}
	if code < 1000 {
		code = code + 1000
	}
	c.JSON(200, RespJson{
		Success: ok,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func Ok(c *gin.Context) {
	Response(200, MsgOk, nil, c)
}

func OkWithData(data any, c *gin.Context) {
	Response(200, Msg200, data, c)
}

func Fail(code int, c *gin.Context) {
	var msg string
	switch code {
	case 400:
		msg = Msg400
	case 401:
		msg = Msg401
	case 403:
		msg = Msg403
	case 404:
		msg = Msg404
	case 405:
		msg = Msg405
	case 500:
		msg = Msg500
	default:
		msg = MsgFail
	}
	Response(code, msg, nil, c)
}

func FailWithMsg(code int, msg string, c *gin.Context) {
	Response(code, msg, nil, c)
}
