package main

import (
	"flag"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hoisie/web"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/tokenTxnBackend/v1"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/utiles"
)

var (
	dbserver = flag.String("dbserver", "localhost:3306", "Database Address.") //"49.51.138.248:3306"
)

func init() {
	flag.Parse()
	fmt.Print("Enter DB Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("\nPassword typed fail: " + err.Error())
	}
	fmt.Println("")
	password := string(bytePassword)
	password = strings.TrimSpace(password)
	utiles.InitMysql(*dbserver, password, false)
}

func main() {
	web.Post("/", v1.Route)
	web.Run("0.0.0.0:8090")
}
