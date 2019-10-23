package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/types"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/utiles"
	//"time"
)

var (
	tokenInfos   TokenInfoArray
	newTokenChan = make(chan types.TokenInfo)
	reset        = flag.Bool("reset", false, "Clear All DB Data.")
	wtcnode      = flag.String("wtcnode", "192.168.50.184:8545", "WTC chain Node Address.")
	dbserver     = flag.String("dbserver", "49.51.138.248:3306", "Database Address.")
	chainNode    *utiles.Chain
)

//TokenInfoArray TokenInfoArray
type TokenInfoArray []types.TokenInfo

//Contains Contains
func (t TokenInfoArray) Contains(ti types.TokenInfo) bool {
	for _, t := range t {
		if t.Address == ti.Address {
			return true
		}
	}
	return false
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func queryTokenInfo(token string) (string, *big.Int, *big.Int, error) {
	tockenCall := utiles.NewTokenCall(*wtcnode, token)
	symbol, err := tockenCall.Symbol()
	if err != nil {
		return "", big.NewInt(0), big.NewInt(0), fmt.Errorf("Get Symbol Fail")
	}
	supply, err := tockenCall.TotalSupply()
	if err != nil {
		return "", big.NewInt(0), big.NewInt(0), fmt.Errorf("Get TotalSupply Fail")
	}
	decimals, err := tockenCall.Decimals()
	if err != nil {
		return "", big.NewInt(0), big.NewInt(0), fmt.Errorf("Get Decimals Fail")
	}
	return symbol, supply, decimals, nil
}

func queryTokenBalance(token, holder string) (*big.Int, error) {
	tockenCall := utiles.NewTokenCall(*wtcnode, token)
	balance, err := tockenCall.BalanceOf(holder)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("Get Token %s :%s Balance Fail", token, holder)
	}

	return balance, nil
}

//5秒检查一次数据库表tokendata.TokenAddress
func refreshTokenAddress() {
	log.Println("[refreshTokenAddress]Start Refresh ...")
	c := time.Tick(time.Duration(5) * time.Second)
	for {
		tokenInfo := utiles.GetTokenAddressesInSQL()
		for _, info := range tokenInfo {
			if info.Status != "" {
				break
			}
			if tokenInfos.Contains(info) == false {
				//查询token信息
				symbol, supply, decimals, err := queryTokenInfo(info.Address)
				if err != nil {
					utiles.UpdateTokenStatus(info.Address, err.Error())
					break
				}

				//创建表项，保存token信息
				utiles.CreateTokenInfoTable(symbol)
				utiles.CreateTokenBalanceTable(symbol)
				utiles.CreateTransferTable(symbol)
				// newTokenTransactionTable("tokenTransaction" + symbol)
				utiles.UpdateTokenInfo(symbol, info.Address, symbol, supply.String(), decimals.String())
				info.Symbol = symbol

				tokenInfos = append(tokenInfos, info)
				log.Printf("[refreshTokenAddress]Got New Token:%+v\n", info)
				newTokenChan <- info
			}
		}
		<-c
	}
}
func checkBlockFork(symbol string, lastCheckedBlock, latestBlockNumber uint64) uint64 {

	// 区块高度太低，不做检查
	if lastCheckedBlock <= 12 || latestBlockNumber <= 12 {
		return lastCheckedBlock
	}
	// 最后检查的区块高度比最新的区块高度小，丢弃超过最新高度的数据
	if latestBlockNumber < lastCheckedBlock {
		utiles.DeleteTokenTransfer(symbol, latestBlockNumber)
		return latestBlockNumber - 1
	}
	// 最后检查的区块和最新区块高度相差12以上
	if latestBlockNumber-lastCheckedBlock >= 12 {
		return lastCheckedBlock
	}
	// 剩下的情况是 latestBlockNumber - lastCheckedBlock < 12
	for i := latestBlockNumber - uint64(12); i <= lastCheckedBlock; i++ {
		fmt.Println("检查块", i)
		// 如果tokenTransfer表中存在该高度的交易记录，和链上的区块hash进行比较，检查是否分叉
		dbhash, err := utiles.GetBlockHashInSQL(symbol, i)
		checkErr(err)
		if dbhash == "" {
			continue
		}

		b := chainNode.Block(i)
		if dbhash != b.Hash {
			fmt.Println("change the ", i, "block")
			utiles.DeleteTokenTransfer(symbol, i)
			return i - 1
		}
	}
	return lastCheckedBlock
}
func grabTransferLogByIterBlock(from, to uint64) {

	for i := from; i <= to; i++ {
		//获取区块数据
		b := chainNode.Block(i)
		//过滤没有交易的区块
		if len(b.Transactions) == 0 {
			continue
		}
		log.Println("扫描到有交易的块: ", utiles.Hextoten(b.Number))
		for _, v := range b.Transactions {
			//获取交易信息
			r := chainNode.TransactionReceipt(v.TransactionHash)

			//过滤没有事件的交易
			logs := r.Logs
			if len(logs) == 0 {
				continue
			}
			log.Println("扫描到有合约事件的块: ", utiles.Hextoten(b.Number))
			for _, event := range logs {
				//过滤不是合约Transfer的交易
				if event.Topics[0] != "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
					continue
				}
				log.Println("扫描到有Tranfer事件的块: ", utiles.Hextoten(b.Number))
				from := "0x" + string([]byte(event.Topics[1])[26:])
				to := "0x" + string([]byte(event.Topics[2])[26:])
				value := big.NewInt(0)
				value.SetString(event.Data, 0)
				//value = value.Div(value, big.NewInt(1000000000))

				log.Printf("Token:%v From：%v To：%v Value:%v\n", event.Address, from, to, value)

			}
		}
	}
}

func grabTransferLogByLogFilter(token types.TokenInfo, ch chan []string) {
	//数据库中解析过的最新高度
	lastBlockNumber, _ := utiles.GetBlockNumberInSQL(token.Symbol)
	//链最新高度
	latestBlockNumber, err := chainNode.GetBlockNumber()
	if err != nil {
		log.Printf("[grabTransferLog]Get blockNumber Exit:%+v", err)
		return
	}
	//根据最新高度检查链是否存在分叉，如果存在，需要删除分叉点后的数据，并重新同步
	comfiredBlockNumber := checkBlockFork(token.Symbol, lastBlockNumber, latestBlockNumber)

	transferlog := chainNode.GetLogs(token.Address, comfiredBlockNumber+1, latestBlockNumber)

	if len(transferlog) == 0 {
		return
	}
	log.Printf("Catch Transfer Logs From:%v To:%v\n", comfiredBlockNumber+1, latestBlockNumber)
	for _, event := range transferlog {
		from := "0x" + string([]byte(event.Topics[1])[26:])
		to := "0x" + string([]byte(event.Topics[2])[26:])
		value := big.NewInt(0)
		value.SetString(event.Data, 0)
		//value = value.Div(value, big.NewInt(1000000000))

		log.Printf("## %v ##Got Transfer: \n\tAddress:%v\n\t  Block:%v\n\t   From：%v\n\t     To：%v\n\t  Value:%v\n\t TxHash:%v\n\n",
			token.Symbol, token.Address, event.BlockNumber, from, to, value, event.TransactionHash)

		number, _ := strconv.ParseUint(event.BlockNumber, 0, 64)
		b := chainNode.Block(number)
		//检查交易是否重复
		amount, _ := utiles.CheckIfExsistTokenTransferByHashInSQL(token.Symbol, to, event.TransactionHash)
		if amount == 0 { //不存在该交易
			utiles.InsertTokenTransfer(
				token.Symbol,
				event.BlockNumber,
				b.Hash,
				event.TransactionHash,
				from,
				to,
				value.String())
		} else {
			log.Printf("TxHash:%v Is Already Exsist!!!\n", event.TransactionHash)
		}
		ch <- []string{from, to}
	}
}

//3秒检查一次数据库表tokendata.TokenAddress
func refreshTokenLog(token types.TokenInfo, ch chan []string) {
	log.Printf("[refreshLog]Start Refresh ... %+v", token)
	c1 := time.Tick(time.Duration(3) * time.Second)
	for {
		select {
		case <-c1:
			grabTransferLogByLogFilter(token, ch)
		}
	}
}
func refreshTokenBalance(token types.TokenInfo, ch chan []string) {
	var holders []string
	var balance []uint64
	log.Printf("[refreshTokenBalance]Start Refresh ... %+v\n", token)
	c := time.Tick(time.Duration(60) * time.Second)
	for {
		select {
		case <-c:
			holders, balance = utiles.GetTokenHoldersInSQL(token.Symbol)
		case holders = <-ch:
			balance = make([]uint64, len(holders))
		}

		for i, holder := range holders {
			if b, err := queryTokenBalance(token.Address, holder); err == nil {
				if b.Uint64() == 0 || b.Uint64() != balance[i] {
					log.Printf("Got %v At ## %v ## Balance. Old=%v New=%v", holder, token.Symbol, balance[i], b.Uint64())
					utiles.UpdateTokenBalance(token.Symbol, holder, b.Uint64())
				}
			} else {
				fmt.Printf("refreshTokenBalance.queryTokenBalance Fail:%v", err)
			}
		}
	}
}

func init() {
	flag.Parse()
	utiles.InitMysql(*dbserver, *reset)
	chainNode = utiles.NewChainNode(*wtcnode)
}

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGTERM)
	//创建保存注册token地址的表
	utiles.CreateTokenAddressTable()
	//定时查询当前注册的token地址名称列表，如果有更新，从链上查询代币信息保存到本地，创建相关的表项
	//并通过newTokenChan通知外界，进行对该token的监听操作
	go refreshTokenAddress()

	//如果有新的token注册
	log.Printf("Wait New Token Register ...")
	for {
		select {
		case token := <-newTokenChan:
			ch := make(chan []string)
			// 监听区块数据，对需要监听的token进行交易解析，并保存到数据库
			go refreshTokenLog(token, ch)
			go refreshTokenBalance(token, ch)
		case s := <-c:
			fmt.Println("Got signal:", s)
			utiles.CloseMysql()
			return
		}
	}
}
