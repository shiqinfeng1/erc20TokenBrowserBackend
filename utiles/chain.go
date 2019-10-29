package utiles

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os/exec"
	"strconv"

	"github.com/shiqinfeng1/erc20TokenBrowserBackend/types"
)

func NewChainNode(url string) *Chain {
	return &Chain{
		EthereumNode: url,
	}
}

type Chain struct {
	EthereumNode string
}

func (c *Chain) Block(number uint64) (types.Block, error) {
	method := "eth_getBlockByNumber"
	params := `["` + "0x" + strconv.FormatUint(number, 16) + `",true]`
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return types.Block{}, err
	}
	r := types.BlcokResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return types.Block{}, err
	}
	return r.Result, nil
}

func (c *Chain) GetLogs(token string, fromBlock, toBlock uint64) ([]types.Log, error) {
	method := "eth_getLogs"
	params := `[{"fromBlock":"0x` + strconv.FormatUint(fromBlock, 16) + `","toBlock":"0x` + strconv.FormatUint(toBlock, 16) + `","address": "` + token + `","topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]` //
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return []types.Log{}, err
	}
	r := types.LogsResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return []types.Log{}, err
	}
	return r.Result, nil
}

func (c *Chain) GetBlockNumber() (uint64, error) {
	method := "eth_blockNumber"
	params := `[]`
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return 0, err
	}
	r := types.BlockNumberResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseUint(r.Result, 0, 64)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func (c *Chain) Getbalance(address string, block uint64) (*big.Int, error) {
	method := "eth_getBalance"
	params := `["` + address + `","0x` + strconv.FormatUint(block, 16) + `"]`
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return nil, err
	}
	r := types.BalanceResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return nil, err
	}
	balance := big.NewInt(0)
	balance.SetString(r.Result, 0)
	return balance, nil
}
func (c *Chain) IsContract(address string) (bool, error) {
	method := "eth_getCode"
	params := `["` + address + `","latest"]`
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return false, err
	}
	r := types.BalanceResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return false, err
	}
	if r.Result == "" {
		return false, nil
	}
	return true, nil
}
func (c *Chain) TransactionReceipt(address string) (types.Receipt, error) {
	method := "eth_getTransactionReceipt"
	params := `["` + address + `"]`
	jsrp, err := c.callGeth(method, params)
	if err != nil {
		return types.Receipt{}, err
	}
	r := types.ReceiptResult{}
	err = json.Unmarshal([]byte(jsrp), &r)
	if err != nil {
		return types.Receipt{}, err
	}
	return r.Result, nil
}
func (c *Chain) callGeth(method, params string) (string, error) {
	cmd := exec.Command("bash", "-c", `curl -H "Content-Type: application/json" -X POST '`+c.EthereumNode+`' --data '{"jsonrpc":"2.0","method":"`+method+`","params":`+params+`,"id":1}'`)
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
