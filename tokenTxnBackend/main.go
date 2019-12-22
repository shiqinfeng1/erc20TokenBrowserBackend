package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hoisie/web"
	v1 "github.com/shiqinfeng1/erc20TokenBrowserBackend/tokenTxnBackend/v1"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/utiles"
)

var (
	dbserver = flag.String("dbserver", "localhost:3306", "Database Address.") //"49.51.138.248:3306"
	dbpwd    = flag.String("dbpwd", "abc123", "Database password.")
)

func init() {
	flag.Parse()
	var fd int
	fmt.Print("Enter DB Password: ")
	switch runtime.GOOS {
	case "darwin", "linux":
		fd = 0
	case "windows":
		fd = int(os.Stdin.Fd())
	}
	bytePassword, err := terminal.ReadPassword(fd)
	if err != nil {
		fmt.Println("\nPassword typed fail: "+err.Error(), "fd=", fd)
	}
	fmt.Println("-----")
	if len(bytePassword) != 0 {
		*dbpwd = string(bytePassword)
	}
	password := *dbpwd
	password = strings.TrimSpace(password)
	utiles.InitMysql(*dbserver, password, false)
}

func main() {
	//
	web.Post("/", v1.Route)
	// 跨域设置
	web.Match("OPTIONS", "/", func(ctx *web.Context) string {
		ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
		ctx.SetHeader("Access-Control-Allow-Methods", "POST", true)
		ctx.SetHeader("Access-Control-Allow-Headers", "accept,content-type,cookieorigin", true)
		ctx.SetHeader("Access-Control-Allow-Credentials", "true", true)
		return ""
	})
	web.Run("0.0.0.0:8090")
}
