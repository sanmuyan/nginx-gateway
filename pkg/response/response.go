package response

import "github.com/gin-gonic/gin"

type Code int

const (
	HttpOk                  Code = 200
	HttpBadRequest          Code = 400
	HttpUnauthorized        Code = 401
	HttpForbidden           Code = 413
	HttpInternalServerError Code = 500
)

func (m Code) GetMessage() string {
	switch m {
	case HttpOk:
		return "操作成功"
	case HttpBadRequest:
		return "数据错误"
	case HttpUnauthorized:
		return "身份验证错误"
	case HttpForbidden:
		return "无权限访问"
	case HttpInternalServerError:
		return "服务器内部错误"
	}
	return ""
}

type Response struct {
	Success bool   `json:"success"`
	Code    Code   `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (r *Response) defaultSet() {
	if r.Code < 600 {
		r.Code += 1000
	}
}

func (r *Response) Ok() *Response {
	code := HttpOk
	r.Code = code
	r.Success = true
	r.Message = code.GetMessage()
	return r
}

func (r *Response) Fail(code Code) *Response {
	r.Code = code
	r.Success = false
	r.Message = code.GetMessage()
	return r
}

func (r *Response) WithData(data any) *Response {
	r.Data = data
	return r
}

func (r *Response) WithMsg(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) SetGin(c *gin.Context) {
	r.defaultSet()
	c.JSON(200, r)
}

func NewResponse() *Response {
	return &Response{
		Success: false,
		Code:    HttpOk,
		Message: "",
		Data:    nil,
	}
}
