package v1

import (
	"net/http"

	"github.com/hoisie/web"
)

type HTTPResponseBody struct {
	c web.Context
}

// ErrorReturns 发生错误的时候的返回值封装
func ErrorReturns(c web.Context, errcode string, msg string) error {
	returns := &ReturnBodyNoPage{
		Errcode: errcode,
		ErrMsg:  msg,
	}
	return c.JSON(http.StatusOK, returns)
}

func Hello(ctx *web.Context, val string) {
	for k, v := range ctx.Params {
		println(k, v)
	}
}
