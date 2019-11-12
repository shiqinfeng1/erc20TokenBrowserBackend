package utiles

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/shiqinfeng1/erc20TokenBrowserBackend/types"
)

var (
	db              *sql.DB
	defaultPassword = "sleewa2018!"
)

//LogInSQL LogInSQL
type LogInSQL struct {
	Value       string `json:"value"`
	From        string `json:"from"`
	To          string `json:"to"`
	BlockNumber uint64 `json:"blockNumber"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func CloseMysql() {
	db.Close()
}
func InitMysql(dbserver, dbPassword string, reset bool) {
	if dbPassword == "" {
		dbPassword = defaultPassword
	}
	//连接数据库
	mysql, err := sql.Open("mysql", "root:"+dbPassword+"@tcp("+dbserver+")/?charset=utf8mb4&parseTime=true&loc=Local")
	checkErr(err)

	err = mysql.Ping()
	if err != nil {
		panic("PANIC when pinging db: " + err.Error()) // proper error handling instead of panic in your app
	}

	mysql.SetMaxIdleConns(0)
	mysql.SetMaxOpenConns(1024)
	mysql.SetConnMaxLifetime(5 * time.Minute)

	db = mysql
	if reset {
		fmt.Println("drop database tokendata ...")
		_, err = db.Exec(`DROP DATABASE IF EXISTS tokendata;`)
		checkErr(err)
		fmt.Println("create database tokendata ...")
		_, err = db.Exec("CREATE DATABASE if not exists tokendata")
		checkErr(err)
	}

	fmt.Println("use tokendata ...")
	_, err = db.Exec("USE tokendata")
	checkErr(err)
}

func CreateTokenAddressTable() {
	sql := `CREATE TABLE IF NOT EXISTS  tokendata.tokenAddress (
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
		status VARCHAR(64) NULL DEFAULT NULL,
        created DATETIME NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (address)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	smt, err := db.Prepare(sql)
	checkErr(err)
	_, err = smt.Exec()
	checkErr(err)
}
func CreateTransferTable(table string) {
	sql := `CREATE TABLE IF NOT EXISTS  tokendata.tokenTransfer` + table + `(
		id INT(10) NOT NULL AUTO_INCREMENT,
		blockNumber INT(10) NULL DEFAULT NULL,
		blockHash VARCHAR(128) NULL DEFAULT NULL,
		timestamp VARCHAR(32) NULL DEFAULT NULL, 
		transferHash VARCHAR(128) NULL DEFAULT NULL,
		sender VARCHAR(64) NULL DEFAULT NULL,
		receiver VARCHAR(64) NULL DEFAULT NULL,
		value VARCHAR(32) NULL DEFAULT NULL,
        created DATETIME NULL DEFAULT NULL,
		PRIMARY KEY(id),
		KEY (blockNumber),
		KEY (sender),
		KEY (receiver)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	smt, err := db.Prepare(sql)
	checkErr(err)
	_, err = smt.Exec()
	checkErr(err)
}

func CreateTokenBalanceTable(table string) {
	sql := `CREATE TABLE IF NOT EXISTS tokendata.tokenBalance` + table + `(
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
        balance VARCHAR(64) NULL DEFAULT NULL,
        created DATETIME NULL DEFAULT NULL,
		PRIMARY KEY(id),
		UNIQUE KEY (address)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	smt, err := db.Prepare(sql)
	checkErr(err)
	_, err = smt.Exec()
	checkErr(err)
}
func CreateTokenInfoTable() {
	sql := `CREATE TABLE IF NOT EXISTS tokendata.tokenMetaInfo (
        id INT(10) NOT NULL AUTO_INCREMENT,
        address VARCHAR(64) NULL DEFAULT NULL,
		name VARCHAR(64) NULL DEFAULT NULL,
		totalSupply VARCHAR(64) NULL DEFAULT NULL,
		decimals VARCHAR(64) NULL DEFAULT NULL,
        created DATETIME NULL DEFAULT NULL,
		PRIMARY KEY(id),
		UNIQUE KEY (address)
    )ENGINE=InnoDB DEFAULT CHARSET=utf8;`

	smt, err := db.Prepare(sql)
	checkErr(err)
	_, err = smt.Exec()
	checkErr(err)
}

func UpdateTokenBalance(table, address string, balance string) error {
	sql := "INSERT INTO tokendata.tokenBalance" + table +
		" (address,balance,created) VALUES (?,?,?) ON DUPLICATE KEY UPDATE balance=?,created=?"
	smt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(address, balance, time.Now().Local(), balance, time.Now().Local())
	if err != nil {
		return err
	}
	return nil
}
func IsContract(address string) (bool, error) {
	row := db.QueryRow("select COUNT(id) from tokendata.tokenAddress where address='" + address + "';")
	var amount int
	row.Scan(&amount)
	if amount > 0 {
		return true, nil
	}
	return false, nil
}
func GetTokenInfo(address string) (types.TokenMetaInfo, error) {
	var metainfo types.TokenMetaInfo
	var name, totalSupply, decimals sql.NullString
	row := db.QueryRow(
		"select name, totalSupply, decimals from tokendata.tokenMetaInfo " +
			"where address='" + address + "'")

	err := row.Scan(&name, &totalSupply, &decimals)
	if err != nil {
		return metainfo, err
	}

	if name.Valid && totalSupply.Valid && decimals.Valid {
		metainfo = types.TokenMetaInfo{
			Address:     address,
			Name:        name.String,
			TotalSupply: totalSupply.String,
			Decimals:    decimals.String,
		}
	}
	return metainfo, nil
}
func GetTokenInfoList() ([]types.TokenMetaInfo, error) {
	var metainfos []types.TokenMetaInfo
	rows, err := db.Query(`select address, name, totalSupply, decimals from tokendata.tokenMetaInfo`)
	if err != nil {
		return metainfos, err
	}
	defer rows.Close()
	for rows.Next() {
		var address, name, totalSupply, decimals sql.NullString
		err = rows.Scan(&address, &name, &totalSupply, &decimals)
		if err != nil {
			return metainfos, err
		}
		if address.Valid && name.Valid && totalSupply.Valid && decimals.Valid {
			metainfos = append(metainfos, types.TokenMetaInfo{
				Address:     address.String,
				Name:        name.String,
				TotalSupply: totalSupply.String,
				Decimals:    decimals.String,
			})
		}
	}
	return metainfos, nil
}

func GetTokenTxnList(table string, pagein types.PageParams) ([]types.TokenTxnInfo, types.PageBody, error) {
	var txninfos []types.TokenTxnInfo
	var pageout types.PageBody
	var total int

	perpage := 0
	if pagein.PerPage > 100 {
		perpage = 100
	} else {
		perpage = pagein.PerPage
	}
	row := db.QueryRow("SELECT COUNT(id) as amount from tokendata.tokenTransfer" + table + ";")
	err := row.Scan(&total)
	if err != nil {
		return txninfos, pageout, err
	}

	rows, err := db.Query(
		"select blockNumber, blockHash, timestamp,transferHash, sender,receiver,value " +
			"from tokendata.tokenTransfer" + table + " " +
			"order by blockNumber desc limit " +
			strconv.Itoa(pagein.CurrentPage*perpage) + "," +
			strconv.Itoa(perpage) + ";")
	if err != nil {
		return txninfos, pageout, err
	}
	defer rows.Close()
	for rows.Next() {
		var blockHash, timestamp, transferHash, sender, receiver, value sql.NullString
		var blockNumber sql.NullInt64
		err = rows.Scan(&blockNumber, &blockHash, &timestamp, &transferHash, &sender, &receiver, &value)
		if err != nil {
			return txninfos, pageout, err
		}
		if blockNumber.Valid && blockHash.Valid && transferHash.Valid && sender.Valid && receiver.Valid {
			txninfos = append(txninfos, types.TokenTxnInfo{
				BlockNumber:  blockNumber.Int64,
				BlockHash:    blockHash.String,
				Timestamp:    timestamp.String,
				TransferHash: transferHash.String,
				Sender:       sender.String,
				Receiver:     receiver.String,
				Value:        value.String,
			})
		}
	}
	pageout.CurrentPage = pagein.CurrentPage
	pageout.PerPage = perpage
	pageout.Total = total
	return txninfos, pageout, nil
}

func GetTokenHolderTxnList(table, holder string, pagein types.PageParams) ([]types.TokenTxnInfo, types.PageBody, error) {
	var txninfos []types.TokenTxnInfo
	var pageout types.PageBody
	var total int

	perpage := 0
	if pagein.PerPage > 100 {
		perpage = 100
	} else {
		perpage = pagein.PerPage
	}
	row := db.QueryRow(
		"SELECT COUNT(id) as amount from tokendata.tokenTransfer" + table +
			" where sender='" + holder + "' or receiver='" + holder + "';")
	err := row.Scan(&total)
	if err != nil {
		return txninfos, pageout, err
	}

	rows, err := db.Query(
		"select blockNumber, blockHash, timestamp,transferHash, sender,receiver,value " +
			"from tokendata.tokenTransfer" + table +
			" where sender='" + holder + "' or receiver='" + holder +
			"' order by blockNumber desc limit " +
			strconv.Itoa(pagein.CurrentPage*perpage) + "," +
			strconv.Itoa(perpage) + ";")
	if err != nil {
		return txninfos, pageout, err
	}
	defer rows.Close()
	for rows.Next() {
		var blockHash, timestamp, transferHash, sender, receiver, value sql.NullString
		var blockNumber sql.NullInt64
		err = rows.Scan(&blockNumber, &blockHash, &timestamp, &transferHash, &sender, &receiver, &value)
		if err != nil {
			return txninfos, pageout, err
		}
		if blockNumber.Valid && blockHash.Valid && transferHash.Valid && sender.Valid && receiver.Valid {
			txninfos = append(txninfos, types.TokenTxnInfo{
				BlockNumber:  blockNumber.Int64,
				BlockHash:    blockHash.String,
				Timestamp:    timestamp.String,
				TransferHash: transferHash.String,
				Sender:       sender.String,
				Receiver:     receiver.String,
				Value:        value.String,
			})
		}
	}
	pageout.CurrentPage = pagein.CurrentPage
	pageout.PerPage = perpage
	pageout.Total = total
	return txninfos, pageout, nil
}

func InsertTokenAddress(address string) error {
	row := db.QueryRow("select COUNT(id) from tokendata.tokenAddress where address='" + address + "';")
	var amount int
	row.Scan(&amount)
	if amount > 0 {
		return nil
	}
	sqldesc := `insert into tokendata.tokenAddress set address=?, created=?`
	smt, err := db.Prepare(sqldesc)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(address, time.Now().Local())
	if err != nil {
		return err
	}
	return nil
}
func UpdateTokenInfo(address, name, totalSupply, decimals string) error {
	sqldesc := `insert into tokendata.tokenMetaInfo` +
		` (address, name, totalSupply, decimals, created) VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE ` +
		`name=?, totalSupply=?, decimals=?, created=?`
	smt, err := db.Prepare(sqldesc)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(address, name, totalSupply, decimals, time.Now().Local(), name, totalSupply, decimals, time.Now().Local())
	if err != nil {
		return err
	}
	return nil
}
func UpdateTokenStatus(address, status string) error {
	sqldesc := `update tokendata.tokenAddress set status=? WHERE address=?`
	smt, err := db.Prepare(sqldesc)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(status, address)
	if err != nil {
		return err
	}
	return nil
}
func InsertTokenTransfer(table, blockNumber, blockHash, timestamp, transferHash, from, to, value string) error {
	sqldesc := `insert into tokendata.tokenTransfer` + table +
		` set blockNumber=?, blockHash=?, timestamp=?, transferHash=?, sender=?, receiver=?, value=?, created=?`
	smt, err := db.Prepare(sqldesc)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec(HexStoTenSWith0x(blockNumber), blockHash, timestamp, transferHash, from, to, value, time.Now().Local())
	if err != nil {
		return err
	}
	return nil
}
func CheckIfExsistTokenTransferByHashInSQL(table, to, transferHash string) (int, error) {
	var amount int
	row := db.QueryRow("SELECT COUNT(id) as amount from tokendata.tokenTransfer" + table +
		" where transferHash = '" + transferHash + "' and receiver = '" + to + "';")
	err := row.Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
func GetTokengTransferByAddressInSQL(table string, startBlockNumber, endBlockNumber uint64, tokenAddress string) ([]LogInSQL, error) {
	var loginsql []LogInSQL
	rows, err := db.Query("SELECT (blockNumber,sender,receiver,value) from tokendata.tokenTransfer" + table +
		" where tokenAddress = " + tokenAddress +
		" and blockNumber >= " + strconv.FormatUint(startBlockNumber, 10) +
		" and blockNumber <= " + strconv.FormatUint(endBlockNumber, 10) +
		";")
	if err != nil {
		return []LogInSQL{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var from, to, value sql.NullString
		var blockNumber uint64
		err = rows.Scan(&blockNumber, &from, &to, &value)
		if err != nil {
			return loginsql, err
		}
		if from.Valid && to.Valid {
			loginsql = append(loginsql, LogInSQL{BlockNumber: blockNumber, From: from.String, To: to.String, Value: value.String})
		}
	}
	return loginsql, nil
}
func GetBlockHashInSQL(table string, blockNumber uint64) (string, error) {
	var hash string
	row := db.QueryRow("select blockHash from tokendata.tokenTransfer" + table + " where blockNumber = " + strconv.FormatUint(blockNumber, 10) + ";")
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}
func DeleteTokenTransfer(table string, block uint64) error {
	sql := "delete FROM tokendata.tokenTransfer" + table + " where blockNumber >= " + strconv.FormatUint(block, 10) + ";"
	fmt.Println("\n" + sql + "\n")
	smt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer smt.Close()
	_, err = smt.Exec()
	if err != nil {
		return err
	}
	return nil
}
func GetTokenAddressByName(name string) (string, error) {
	var address sql.NullString
	row := db.QueryRow("select address from tokendata.tokenMetaInfo where name='" + name + "';")
	err := row.Scan(&address)
	if err != nil {
		return "", err
	}
	if address.Valid {
		return address.String, nil
	}
	return "", nil
}

func GetTokenNameByAddress(address string) (string, error) {
	var name sql.NullString
	row := db.QueryRow("select name from tokendata.tokenMetaInfo where address='" + address + "';")
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	if name.Valid {
		return name.String, nil
	}
	return "", nil
}
func GetHolderBalance(table, holder string) (string, error) {
	var balance sql.NullString
	row := db.QueryRow("select balance from tokendata.tokenBalance" + table + " where address='" + holder + "';")
	err := row.Scan(&balance)
	if err != nil {
		return "", err
	}
	if balance.Valid {
		return balance.String, nil
	}
	return "", nil
}
func GetTokenAddressesInSQL() ([]types.TokenInfo, error) {
	var tokenInfo []types.TokenInfo
	rows, err := db.Query("select address,status from tokendata.tokenAddress;")
	if err != nil {
		return tokenInfo, err
	}
	defer rows.Close()

	for rows.Next() {
		var a sql.NullString
		var s sql.NullString
		err = rows.Scan(&a, &s)
		if err != nil {
			return tokenInfo, err
		}
		tokenInfo = append(tokenInfo, types.TokenInfo{Address: a.String, Status: s.String})
	}

	return tokenInfo, nil
}

func GetBlockNumberInSQL(table string) (uint64, error) {
	var blocknumber sql.NullInt64
	row := db.QueryRow("SELECT MAX(blockNumber) from tokendata.tokenTransfer" + table + ";")
	err := row.Scan(&blocknumber)
	if err != nil {
		return 0, err
	}
	if blocknumber.Valid {
		return uint64(blocknumber.Int64), nil
	} else {
		return 0, nil
	}
}

func GetTokenHoldersInSQL(table string) ([]string, []string, error) {
	var holders []string
	var balances []string
	rows, err := db.Query("select address,balance from tokendata.tokenBalance" + table + ";")
	if err != nil {
		return holders, balances, err
	}
	defer rows.Close()

	for rows.Next() {
		var s sql.NullString
		var b sql.NullString
		err = rows.Scan(&s, &b)
		if err != nil {
			return holders, balances, err
		}
		if s.Valid {
			holders = append(holders, s.String)
			balances = append(balances, b.String)
		}
	}
	return holders, balances, nil
}

func GetTokenHolderList(table string, pagein types.PageParams) ([]types.TokenHolderInfo, types.PageBody, error) {
	var holders []types.TokenHolderInfo
	var pageout types.PageBody
	var total int

	perpage := 0
	if pagein.PerPage > 100 {
		perpage = 100
	} else {
		perpage = pagein.PerPage
	}
	row := db.QueryRow("SELECT COUNT(id) as amount from tokendata.tokenBalance" + table + ";")
	err := row.Scan(&total)
	if err != nil {
		return holders, pageout, err
	}

	rows, err := db.Query(
		"select address,balance " +
			"from tokendata.tokenBalance" + table + " " +
			"order by balance+0 desc limit " +
			strconv.Itoa(pagein.CurrentPage*perpage) + "," +
			strconv.Itoa(perpage) + ";")
	if err != nil {
		return holders, pageout, err
	}
	defer rows.Close()
	for rows.Next() {
		var a sql.NullString
		var b sql.NullString
		err = rows.Scan(&a, &b)
		if err != nil {
			return holders, pageout, err
		}
		if a.Valid && b.Valid {
			holders = append(holders,
				types.TokenHolderInfo{
					Address: a.String,
					Balance: b.String})
		}
	}
	pageout.CurrentPage = pagein.CurrentPage
	pageout.PerPage = perpage
	pageout.Total = total
	return holders, pageout, nil
}
