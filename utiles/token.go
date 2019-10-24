package utiles

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

//NewTokenCall --
func NewTokenCall(url, contractAddress string) *TokenCall {
	if url[:4] != "http" {
		url = "http://" + url
	}
	rpcDial, err := rpc.Dial(url)
	if err != nil {
		panic(err)
	}

	client := ethclient.NewClient(rpcDial)
	return &TokenCall{client: client, contractAddress: contractAddress}
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

//Decimals Decimals
func (s *TokenCall) Decimals() (*big.Int, error) {
	token, err := NewToken(common.HexToAddress(s.contractAddress), s.client)
	if err != nil {
		return nil, fmt.Errorf("Decimals.NewToken Fail:%v", err)
	}
	return token.Decimals(&bind.CallOpts{})
}
