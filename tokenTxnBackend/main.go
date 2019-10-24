package main

import (
	"github.com/hoisie/web"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/tokenTxnBackend/v1"
)

func main() {
	web.Get("/(.*)", v1.Hello)
	web.Run("0.0.0.0:9999")

	web.Post("v1/token/info", v1.Hello)              //指定token的信息
	web.Post("v1/token/register", v1.Hello)          //注册token
	web.Post("v1/token/list/transaction", v1.Hello)  //指定token的所有交易列表
	web.Post("v1/token/list/holders", v1.Hello)      //指定token的所有持有者列表，按照余额大小排序
	web.Post("v1/holder/list/transaction", v1.Hello) //holder在指定token中的交易记录
	web.Post("v1/holder/balance", v1.Hello)          //holder在指定token中的余额
}
