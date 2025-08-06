package result

import "net/http"

type Result struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Ok() *Result {
	return &Result{
		Code:    http.StatusOK,
		Message: "请求成功",
		Data:    nil,
	}
}

func Err() *Result {
	return &Result{
		Code:    http.StatusBadRequest,
		Message: "请求失败",
		Data:    nil,
	}
}

func (r *Result) SetData(value any) *Result {
	r.Data = value
	return r
}

func (r *Result) SetMsg(msg string) *Result {
	r.Message = msg
	return r
}

func (r *Result) SetCode(code int32) *Result {
	r.Code = code
	return r
}
