package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type HTTPResponseBody struct {
	c echo.Context
}

// ErrorReturns 发生错误的时候的返回值封装
func ErrorReturns(c echo.Context, errcode string, msg string) error {
	returns := &ReturnBodyNoPage{
		Errcode: errcode,
		ErrMsg:  msg,
	}
	return c.JSON(200, returns)
}


func Hello(ctx *web.Context, val string) {
	for k, v := range ctx.Params {
		println(k, v)
	}
}
