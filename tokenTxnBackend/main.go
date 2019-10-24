package main

import (
	"github.com/hoisie/web"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/tokenTxnBackend/v1"
)

func main() {
	web.Post("/v1", v1.Route)
	web.Run("0.0.0.0:8090")
}
