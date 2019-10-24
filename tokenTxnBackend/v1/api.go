package v1

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hoisie/web"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/types"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/utiles"
)

// ErrorReturns 发生错误的时候的返回值封装
func ErrorReturns(id, errcode int, errmsg string) *types.JSONRPCResponse {
	returns := &types.JSONRPCResponse{
		ID:      id,
		Jsonrpc: "2.0",
		Result:  types.ReturnBodyNoPage{ErrCode: errcode, ErrMsg: errmsg},
	}
	return returns
}

// ResultNoPageReturns 返回值封装
func ResultNoPageReturns(id int, data interface{}) *types.JSONRPCResponse {
	returns := &types.JSONRPCResponse{
		ID:      id,
		Jsonrpc: "2.0",
		Result:  types.ReturnBodyNoPage{ErrCode: 0, ErrMsg: "", Data: data},
	}
	return returns
}

// ResultWithPageReturns 返回值封装
func ResultWithPageReturns(id int, data interface{}, page types.PageBody) *types.JSONRPCResponse {
	returns := &types.JSONRPCResponse{
		ID:      id,
		Jsonrpc: "2.0",
		Result: types.ReturnBodyWithPage{
			ErrCode: 0, ErrMsg: "",
			Data: data,
			Page: page,
		},
	}
	return returns
}
func praseRequest(ctx *web.Context) (*types.JSONRPCRequest, error) {
	var req = types.JSONRPCRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		return &types.JSONRPCRequest{}, err
	}
	if req.Jsonrpc != "2.0" {
		return &types.JSONRPCRequest{}, fmt.Errorf("JSONRPC Request Version Mismatch: %v", req.Jsonrpc)
	}
	return &req, nil
}
func getTokenInfo(id int) *types.JSONRPCResponse {
	tokeninfos, err := utiles.GetTokenInfo()
	if err != nil {
		return ErrorReturns(id, 1003, "Get Token Info Fail: "+err.Error())
	}
	return ResultNoPageReturns(id, tokeninfos)
}
func getTokenTxnList(id int, params string) *types.JSONRPCResponse {

	p := types.TokenTxnListParams{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		return ErrorReturns(id, 1004, "Token Txn List Params Error: "+err.Error())
	}
	var symbol string
	if p.Token[:2] == "0x" || p.Token[:2] == "0X" {
		if symbol, err = utiles.GetTokenNameByAddress(p.Token); err != nil {
			return ErrorReturns(id, 1005, "Get Token Symbol Fail: "+err.Error())
		}
	} else {
		symbol = p.Token
	}
	txnList, page, err := utiles.GetTokenTxnList(symbol, p.Page)
	if err != nil {
		return ErrorReturns(id, 1005, "Get Token Txn List Fail: "+err.Error())
	}
	return ResultWithPageReturns(id, txnList, page)
}
func Handle(req *types.JSONRPCRequest) *types.JSONRPCResponse {

	switch req.Method {
	case "get_tokenInfo": //指定token的信息
		return getTokenInfo(req.ID)
	case "get_tokenTxnList": //指定token的所有交易列表
		return getTokenTxnList(req.ID, req.Params)
	case "get_tokenHolderList": //指定token的所有持有者列表，按照余额大小排序
	case "get_holderTxnList": //holder在指定token中的交易记录
	case "get_holderBalance": //holder在指定token中的余额
	case "token_register": //注册token

	}
	return ErrorReturns(req.ID, 1002, "Unkown Method: "+req.Method)
}

//Route Route
func Route(ctx *web.Context) string {

	ctx.ContentType("json")

	req, err := praseRequest(ctx)
	if err != nil {
		log.Println("praseRequest Fial: ", err.Error())
		data, _ := json.Marshal(ErrorReturns(req.ID, 1001, "Unkown Request Body"))
		return string(data)
	}
	data, _ := json.Marshal(Handle(req))
	return string(data)
}
