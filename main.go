package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os/exec"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	_ "github.com/go-sql-driver/mysql"

	"github.com/ethereum/go-ethereum/common"
	//"time"
)

//BlockNumberResult BlockNumberResult
type BlockNumberResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

//BalanceResult BalanceResult
type BalanceResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
}

//ReceiptResult ReceiptResult
type ReceiptResult struct {
	ID      int     `json:"id"`
	JSONRPC string  `json:"jsonrpc"`
	Result  Receipt `json:"result"`
}

//BlcokResult BlcokResult
type BlcokResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  Block  `json:"result"`
}

//LogsResult LogsResult
type LogsResult struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Result  []Log  `json:"result"`
}

type sync struct {
	CurrentBlock  string `json:"currentBlock"`
	HighestBlock  string `json:"highestBlock"`
	KnownStates   string `json:"knownStates"`
	PulledStates  string `json:"pulledStates"`
	StartingBlock string `json:"startingBlock"`
}

//Block Block
type Block struct {
	Difficulty      string        `json:"difficulty"`
	GasLimit        string        `json:"gasLimit"`
	GasUsed         string        `json:"gasUsed"`
	Hash            string        `json:"hash"`
	Miner           string        `json:"miner"`
	Nonce           string        `json:"nonce"`
	Number          string        `json:"number"`
	ParentHash      string        `json:"parentHash"`
	Size            string        `json:"size"`
	Timestamp       string        `json:"timestamp"`
	TotalDifficulty string        `json:"totalDifficulty"`
	Transactions    []Transaction `json:"transactions"`
}

//Transaction Transaction
type Transaction struct {
	TransactionHash  string `json:"hash"`
	Input            string `json:"input"`
	From             string `json:"from"`
	To               string `json:"to"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasprice"`
	Nonce            string `json:"nonce"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
}

//Receipt Receipt
type Receipt struct {
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	LogsBloom        string `json:"logsBloom"`
	Logs             []Log  `json:"logs"`
	FromAddress      string `json:"from"`
	ToAddress        string `json:"to"`
	GasUsed          string `json:"gasUsed"`
}

//Log Log
type Log struct {
	Address          string   `json:"address"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Data             string   `json:"data"`
	Topics           []string `json:"topics"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
}

//LogInSQL LogInSQL
type LogInSQL struct {
	Data        string `json:"data"`
	Topics      string `json:"topics"`
	BlockNumber uint64 `json:"blockNumber"`
}

const (
	debugMode    = false
	dbServer     = "wuchou@tcp(138.68.3.61:3306)/tokendata"
	ethereumNode = "127.0.0.1:8545"
)

var (
	tokenInfos   TokenInfoArray
	tradersChan  chan []string
	newTokenChan chan TokenInfo
	db           *sql.DB
)

//TokenInfoArray TokenInfoArray
type TokenInfoArray []TokenInfo

//Contains Contains
func (t TokenInfoArray) Contains(ti TokenInfo) bool {
	for _, t := range t {
		if t.Address == ti.Address {
			return true
		}
	}
	return false
}

//TokenCall TokenCall
type TokenCall struct {
	client          *ethclient.Client
	contractAddress string
}

//Symbol Symbol
func (s *TokenCall) Symbol() (string, error) {
	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return "", fmt.Errorf("Symbol.NewToken Fail:%v", err)
	}
	return token.Symbol(&bind.CallOpts{})
}

//BalanceOf BalanceOf
func (s *TokenCall) BalanceOf(holder string) (*big.Int, error) {
	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return nil, fmt.Errorf("BalanceOf.NewToken Fail:%v", err)
	}
	return token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(holder))
}

//TotalSupply TotalSupply
func (s *TokenCall) TotalSupply() (*big.Int, error) {
	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return nil, fmt.Errorf("TotalSupply.NewToken Fail:%v", err)
	}
	return token.TotalSupply(&bind.CallOpts{})
}

//InitTokenCall InitTokenCall
func InitTokenCall(url, contractAddress string) *TokenCall {
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(rpcDial)
	return &TokenCall{client: client, contractAddress: contractAddress}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// TokenInfo TokenInfo
type TokenInfo struct {
	Address      string
	LatestNumber uint64
	Symbol       string
}

func dropTable(table string) {
	sql := `DROP TABLE IF EXISTS ` + table + `;`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func createTokenAddressTable() {
	sql := `CREATE TABLE IF NOT EXISTS  TokenAddress (
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
        latestNumber INT(10) NULL DEFAULT NULL,
        created DATE NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (address),
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func createTransferTable() {
	sql := `CREATE TABLE IF NOT EXISTS  TokenTransfer (
		id INT(10) NOT NULL AUTO_INCREMENT,
		tokenAddress VARCHAR(64) NULL DEFAULT NULL,
		blockNumber INT(10) NULL DEFAULT NULL,
        blockHash VARCHAR(64) NULL DEFAULT NULL,
		transferHash VARCHAR(64) NULL DEFAULT NULL,
		logTopic VARCHAR(256) NULL DEFAULT NULL,
		logData VARCHAR(256) NULL DEFAULT NULL,
        created DATE NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (tokenAddress),
		KEY (blockNumber)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}

func createTokenBalanceTable(table string) {
	sql := `CREATE TABLE IF NOT EXISTS ` + table + `(
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
        balance VARCHAR(64) NULL DEFAULT NULL,
        created DATE NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (address)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func createTokenInfoTable(table string) {
	sql := `CREATE TABLE IF NOT EXISTS ` + table + `(
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
		name VARCHAR(64) NULL DEFAULT NULL,
		totalSupply VARCHAR(64) NULL DEFAULT NULL,
        created DATE NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (address)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}

func updateTokenAddresses(address string, latestBlock uint64) {
	sql := "UPDATE tokendata.TokenAddress SET latestNumber=" + strconv.FormatUint(latestBlock, 10) + " WHERE address=" + address + ";"
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func updateTokenBalance(table, address string, balance uint64) {
	sql := "UPDATE tokendata.TokenBalance SET balance=" + strconv.FormatUint(balance, 10) + " WHERE address=" + address + ";"
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func insertTokenInfo(table, address, name, totalSupply string) {
	sql := `insert into ` + table +
		`(address, name, totalSupply) value (` + address + ", " + name + ", " + totalSupply + `);`

	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}

func insertTokenTransfer(tokenAddress, blockNumber, blockHash, transferHash, logTopic, logData string) {
	sql := "insert into tokenTransfer (tokenAddress,blockNumber, blockHash, transferHash, logTopic, logData) values (\"" + tokenAddress + "\", " + blockNumber + "\", " + blockHash + ", \"" + transferHash + "\", \"" + logTopic + "\", " + logData + ");"
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}
func getTokengTransferByHashInSQL(transferHash string) (uint64, error) {
	var blocknumber uint64
	row := db.QueryRow("SELECT (blockNumber) from TokenTransfer where transferHash = " + transferHash + ";")
	err := row.Scan(&blocknumber)
	if err != nil {
		return 0, err
	}
	return blocknumber, nil
}
func getTokengTransferByAddressInSQL(startBlockNumber, endBlockNumber uint64, tokenAddress string) ([]LogInSQL, error) {
	var loginsql []LogInSQL
	rows, err := db.Query("SELECT (blockNumber,logTopic, logData) from TokenTransfer where tokenAddress = " +
		tokenAddress +
		" and blockNumber >= " + strconv.FormatUint(startBlockNumber, 10) +
		" and blockNumber <= " + strconv.FormatUint(endBlockNumber, 10) +
		";")
	if err != nil {
		return []LogInSQL{}, err
	}
	for rows.Next() {
		var topics, data string
		var blockNumber uint64
		err = rows.Scan(&blockNumber, &topics, &data)
		if err != nil {
			log.Fatalln(err)
		}
		loginsql = append(loginsql, LogInSQL{BlockNumber: blockNumber, Topics: topics, Data: data})
	}
	rows.Close()
	return loginsql, nil
}
func getBlockHashInSQL(blockNumber uint64) (string, error) {
	var hash string
	row := db.QueryRow("select blockHash from tokenTransfer where blockNumber = " + strconv.FormatUint(blockNumber, 10) + ";")
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}
func deleteTokenTransfer(block uint64) {
	sql := "delete FROM tokenTransfer where blockNumber >= " + strconv.FormatUint(block, 10) + ";"
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	checkErr(err)
	smt.Exec()
}

func getTokenAddressesInSQL() []TokenInfo {
	var tokenInfo []TokenInfo
	rows, err := db.Query("select address from TokenAddress;")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var s string
		var n uint64
		err = rows.Scan(&s, &n)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("TokenAddress containing %q", s)
		tokenInfo = append(tokenInfo, TokenInfo{Address: s, LatestNumber: n})
	}
	rows.Close()
	return tokenInfo
}
func getBlockNumberInSQL() (uint64, error) {
	var blocknumber uint64
	row := db.QueryRow("SELECT MAX(blockNumber) from TokenTransfer;")
	err := row.Scan(&blocknumber)
	if err != nil {
		return 0, err
	}
	return blocknumber, nil
}
func getSyncedBlockNumberInSQL(table string) (uint64, error) {
	var blocknumber uint64
	row := db.QueryRow("SELECT MAX(blockNumber) from " + table + ";")
	err := row.Scan(&blocknumber)
	if err != nil {
		return 0, err
	}
	return blocknumber, nil
}

func getTokenHoldersInSQL(table string) ([]string, []uint64) {
	var holders []string
	var balances []uint64
	rows, err := db.Query("select address,balance from " + table + ";")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var s string
		var b uint64
		err = rows.Scan(&s, &b)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("TokenBalance containing %q", s)
		holders = append(holders, s)
		balances = append(balances, b)
	}
	rows.Close()
	return holders, balances
}

func newTokenBalanceTable(table string) {
	// dropTable(db, table)
	createTokenBalanceTable(table)
}
func newTokenInfoTable(table string) {
	// dropTable(db, table)
	createTokenInfoTable(table)
}

// func newTokenTransactionTable(table string) {
// 	// dropTable(db, table)
// 	createTokenTransactionTable(table)
// }

func queryTokenInfo(token string) (string, *big.Int, error) {
	tockenCall := InitTokenCall("http://localhost:8545", token)
	symbol, err := tockenCall.Symbol()
	if err != nil {
		return "", big.NewInt(0), fmt.Errorf("Get Token %s Symbol Fail", token)
	}
	supply, err := tockenCall.TotalSupply()
	if err != nil {
		return "", big.NewInt(0), fmt.Errorf("Get Token %s supply Fail", token)
	}
	return symbol, supply, nil
}

func queryTokenBalance(token, holder string) (*big.Int, error) {
	tockenCall := InitTokenCall("http://localhost:8545", token)
	balance, err := tockenCall.BalanceOf(holder)
	if err != nil {
		return big.NewInt(0), fmt.Errorf("Get Token %s :%s Balance Fail", token, holder)
	}

	return balance, nil
}

//5秒检查一次数据库表tokendata.TokenAddress
func refreshTokenAddress() {
	c := time.Tick(time.Duration(5) * time.Second)
	for {
		tokenInfo := getTokenAddressesInSQL()
		for _, info := range tokenInfo {
			if tokenInfos.Contains(info) == false {
				//查询token信息
				symbol, supply, err := queryTokenInfo(info.Address)
				if err != nil {
					fmt.Printf("queryTokenInfo fail:%v\n", err)
					break
				}
				//创建表项，保存token信息
				newTokenInfoTable("tokenMetaInfo" + symbol)
				newTokenBalanceTable("tokenBalance" + symbol)
				// newTokenTransactionTable("tokenTransaction" + symbol)
				insertTokenInfo("tokenMetaInfo"+symbol, info.Address, symbol, supply.String())
				info.Symbol = symbol

				tokenInfos = append(tokenInfos, info)
				newTokenChan <- info
			}
		}
		<-c
	}
}
func checkBlockFork(lastCheckedBlock, latestBlockNumber uint64) uint64 {

	// 区块高度太低，不做检查
	if lastCheckedBlock <= 12 || latestBlockNumber <= 12 {
		return lastCheckedBlock
	}
	// 最后检查的区块高度比最新的区块高度小，丢弃超过最新高度的数据
	if latestBlockNumber < lastCheckedBlock {
		deleteTokenTransfer(latestBlockNumber)
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
		dbhash, err := getBlockHashInSQL(i)
		checkErr(err)
		if dbhash == "" {
			continue
		}
		b := block(i)
		if dbhash != b.Hash {
			fmt.Println("change the ", i, "block")
			deleteTokenTransfer(i)
			return i - 1
		}
	}
	return lastCheckedBlock
}
func grabTransferLog() {
	var syncedBlockNumber uint64
	//数据库中解析过的最新高度
	lastBlockNumber, _ := getBlockNumberInSQL()
	//链最新高度
	latestBlockNumber := blockNumber()
	//根据最新高度检查链是否存在分叉，如果存在，需要删除分叉点后的数据，并重新同步
	comfiredBlockNumber := checkBlockFork(lastBlockNumber, latestBlockNumber)

	for i := comfiredBlockNumber + 1; i <= syncedBlockNumber; i++ {
		//获取区块数据
		b := block(i)
		fmt.Println("扫描块: ", hextoten(b.Number))
		//过滤没有交易的区块
		if len(b.Transactions) == 0 {
			continue
		}

		for _, v := range b.Transactions {
			//获取交易信息
			r := transactionReceipt(v.TransactionHash)

			//过滤没有事件的交易
			logs := r.Logs
			if len(logs) == 0 {
				continue
			}
			for _, event := range logs {

				//过滤不是合约Transfer的交易
				if event.Topics[0] != "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
					continue
				}
				topic, _ := json.Marshal(event.Topics[1:])
				//检查交易是否重复
				blockNumber, _ := getTokengTransferByHashInSQL(event.TransactionHash)
				number, _ := strconv.ParseUint(event.BlockNumber, 0, 64)
				if blockNumber != number { //不存在该交易
					insertTokenTransfer(v.To, event.BlockNumber, b.Hash, event.TransactionHash, string(topic), event.Data)
				}
				from := "0x" + string([]byte(event.Topics[1])[26:])
				to := "0x" + string([]byte(event.Topics[2])[26:])
				tradersChan <- []string{from, to}
			}
		}
	}
}

//3秒检查一次数据库表tokendata.TokenAddress
func refreshBlock() {
	c1 := time.Tick(time.Duration(3) * time.Second)
	for {
		select {
		case <-c1:
			grabTransferLog()
		}
	}
}
func refreshTokenBalance(token TokenInfo) {
	var holders []string
	var balance []uint64
	c := time.Tick(time.Duration(300) * time.Second)
	for {
		select {
		case <-c:
			holders, balance = getTokenHoldersInSQL("tokenBalance" + token.Symbol)
		case holders = <-tradersChan:
			balance = []uint64{}
		}

		for i, holder := range holders {
			if b, err := queryTokenBalance(token.Address, holder); err == nil {
				if b.Uint64() != balance[i] {
					updateTokenBalance("tokenBalance"+token.Symbol, holder, b.Uint64())
				}
			} else {
				fmt.Printf("refreshTokenBalance.queryTokenBalance Fail:%v", err)
			}
		}
	}
}
func createDatabase() {
	sql := `create database tokendata; `
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	fmt.Printf("createDatabase: %v\n", err)
	smt.Exec()
}

func init() {
	//连接数据库
	mysql, err := sql.Open("mysql", "root:"+dbServer+"?charset=utf8")
	checkErr(err)
	defer mysql.Close()

	db = mysql

	createDatabase()
}

func main() {

	//创建保存注册token地址的表
	createTokenAddressTable()
	//定时查询当前注册的token地址名称列表，如果有更新，从链上查询代币信息保存到本地，创建相关的表项
	//并通过newTokenChan通知外界，进行对该token的监听操作
	go refreshTokenAddress()

	//创建交易记录表，定时刷新区块交易
	createTransferTable()
	go refreshBlock()

	//如果有新的token注册
	select {
	case token := <-newTokenChan:

		// 监听区块数据，对需要监听的token进行交易解析，并保存到数据库
		go refreshTokenBalance(token)
	}
}

func block(number uint64) Block {
	method := "eth_getBlockByNumber"
	params := `["` + "0x" + strconv.FormatUint(number, 16) + `",true]`
	jsrp, err := callGeth(method, params)
	checkErr(err)
	r := BlcokResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	checkErr(err)
	return r.Result
}

func logs(fromBlock, toBlock uint64) []Log {
	method := "eth_getLogs"
	params := `[{"fromBlock":"0x` + strconv.FormatUint(fromBlock, 16) + `","toBlock":"0x` + strconv.FormatUint(toBlock, 16) + `","address": "0xb7cb1c96db6b22b0d3d9536e0108d062bd488f74","topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]`
	jsrp, err := callGeth(method, params)
	checkErr(err)
	r := LogsResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	checkErr(err)
	return r.Result
}

func blockNumber() uint64 {
	method := "eth_blockNumber"
	params := `[]`
	jsrp, err := callGeth(method, params)
	checkErr(err)
	r := BlockNumberResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	checkErr(err)
	number, err := strconv.ParseUint(r.Result, 0, 64)
	checkErr(err)
	return number
}

func transactionReceipt(address string) Receipt {
	method := "eth_getTransactionReceipt"
	params := `["` + address + `"]`
	jsrp, err := callGeth(method, params)
	checkErr(err)
	r := ReceiptResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	checkErr(err)
	return r.Result
}

func getbalance(address string, block uint64) *big.Int {
	method := "eth_getBalance"
	params := `["` + address + `","0x` + strconv.FormatUint(block, 16) + `"]`
	jsrp, err := callGeth(method, params)
	checkErr(err)
	r := BalanceResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	checkErr(err)
	balance := big.NewInt(0)
	balance.SetString(r.Result, 0)
	return balance
}

func callGeth(method, params string) (string, error) {
	cmd := exec.Command("bash", "-c", `curl -H "Content-Type: application/json" -X POST '`+ethereumNode+`' --data '{"jsonrpc":"2.0","method":"`+method+`","params":`+params+`,"id":1}'`)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		return "", err
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}

//RunSQLQuery RunSQLQuery
func RunSQLQuery(command string) error {
	if debugMode {
		fmt.Println(command)
	} else {
		_, err := db.Exec(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertBlockCom(b Block, reward string) string {
	diff := HexStoTenSWith0x(b.Difficulty)
	return `insert into Block (BlockNumber, BlockHash, BlockMiner, BlockDifficulty, GasLimit, GasUsed, Nonce, BlockSize, ParentHash, TransactionCount, UncleCount, BlockTime, BlockReward) values (` + strconv.FormatInt(hextoten(b.Number), 10) + `, "` + b.Hash + `", "` + b.Miner + `", ` + diff + `, ` + strconv.FormatInt(hextoten(b.GasLimit), 10) + `, ` + strconv.FormatInt(hextoten(b.GasUsed), 10) + `, "` + b.Nonce + `", ` + strconv.FormatInt(hextoten(b.Size), 10) + `, "` + b.ParentHash + `", ` + strconv.Itoa(len(b.Transactions)) + `, ` + strconv.Itoa(0) + `, ` + strconv.FormatInt(hextoten(b.Timestamp), 10) + `, ` + reward + `);`
}

func updateRewardCom(block int64, reward *big.Int) string {
	sql := "update Block set BlockReward = " + reward.String() + " where BlockNumber = " + strconv.FormatInt(block, 10) + ";"
	return sql
}

func deleteBlockCom(block uint64) string {
	sql := "delete FROM Block where BlockNumber >= " + strconv.FormatUint(block, 10) + ";"
	return sql
}

//HexStoTenSWith0x HexStoTenSWith0x
func HexStoTenSWith0x(value string) string {
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	return bignumber.String()
}

func hextoten(num string) int64 {
	v := num[2:]
	if s, err := strconv.ParseInt(v, 16, 32); err == nil {
		return s
	}
	return 0
}

//HexStoTenBigInt HexStoTenBigInt
func HexStoTenBigInt(value string) *big.Int {
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	return bignumber
}

//HexStoTenSAndDiv10e9 HexStoTenSAndDiv10e9
func HexStoTenSAndDiv10e9(value string) int64 {
	var B10e9 = big.NewInt(1000000000)
	bignumber := big.NewInt(0)
	bignumber.SetString(value, 0)
	bignumber = bignumber.Div(bignumber, B10e9)
	return bignumber.Int64()
}
func getReward(block, balance *big.Int) *big.Int {
	oneyear := big.NewInt(60 * 24 * 365)
	var reward *big.Int
	big5000 := new(big.Int).Mul(big.NewInt(5000), big.NewInt(1e+18))
	years := new(big.Int).Div(block, oneyear).Int64()

	if years < 2 {
		reward = big.NewInt(5e+18)
		if balance.Cmp(big5000) > 0 {
			reward = reward.Add(reward, big.NewInt(1e+18))
		}
	} else {
		var r int64 = 5e+18
		if balance.Cmp(big5000) >= 0 {
			r = 5e+18
		}
		y := years / 2
		for i := 0; int64(i) < y; i++ {
			r = int64(float64(r) * 0.75)
		}
		reward = big.NewInt(r)
	}
	return reward
}
