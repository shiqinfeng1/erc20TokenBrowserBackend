package v1

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hoisie/web"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/types"
	"github.com/shiqinfeng1/erc20TokenBrowserBackend/utiles"
)

const (
	_ int = 1000 + iota
	ERRCODE_HolderBalance1
	ERRCODE_HolderBalance2
	ERRCODE_HolderBalance3
	ERRCODE_HolderBalance4
	ERRCODE_HolderBalance5
	ERRCODE_TokenRegister1
	ERRCODE_TokenRegister2
	ERRCODE_GetTokenInfo1
	ERRCODE_GetTokenInfo2
	ERRCODE_IsContract1
	ERRCODE_GetTokenTxnList1
	ERRCODE_GetTokenTxnList2
	ERRCODE_GetTokenTxnList3
	ERRCODE_GetTokenTxnList4
	ERRCODE_GetHolderTokenList1
	ERRCODE_GetHolderTokenList2
	ERRCODE_GetHolderTokenList3
	ERRCODE_GetHolderTokenList4
	ERRCODE_GetHolderTokenList5
	ERRCODE_GetTokenHolderList1
	ERRCODE_GetTokenHolderList2
	ERRCODE_GetTokenHolderList3
	ERRCODE_GetTokenHolderList4
	ERRCODE_Handle
	ERRCODE_Route
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
func tokenRegister(id int, token string) *types.JSONRPCResponse {
	if token == "" || (token[:2] != "0x" && token[:2] == "0X") {
		return ErrorReturns(id, ERRCODE_TokenRegister1, "Token Address Invalid: "+token)
	}
	err := utiles.InsertTokenAddress(token)
	if err != nil {
		return ErrorReturns(id, ERRCODE_TokenRegister2, "Token Address Insert Fail: "+err.Error())
	}
	return ResultNoPageReturns(id, "Register OK")
}
func getTokenInfo(id int, params string) *types.JSONRPCResponse {
	if params == "" {
		tokeninfos, err := utiles.GetTokenInfoList()
		if err != nil {
			return ErrorReturns(id, ERRCODE_GetTokenInfo1, "Get Token Info Fail: "+err.Error())
		}
		return ResultNoPageReturns(id, tokeninfos)
	}
	tokeninfo, err := utiles.GetTokenInfo(params)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenInfo2, "Get Token Info Fail: "+err.Error())
	}
	return ResultNoPageReturns(id, tokeninfo)
}
func isContract(id int, params string) *types.JSONRPCResponse {
	istrue, err := utiles.IsContract(params)
	if err != nil {
		return ErrorReturns(id, ERRCODE_IsContract1, "isContract Fail: "+err.Error())
	}
	return ResultNoPageReturns(id, istrue)
}

func getTokenTxnList(id int, params string) *types.JSONRPCResponse {

	p := types.TokenListParams{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenTxnList1, "Token Txn List Params Error: "+err.Error())
	}
	if p.Token == "" {
		return ErrorReturns(id, ERRCODE_GetTokenTxnList4, "Token is NULL ")
	}
	var symbol string
	if p.Token[:2] == "0x" || p.Token[:2] == "0X" {
		if symbol, err = utiles.GetTokenNameByAddress(p.Token); err != nil {
			return ErrorReturns(id, ERRCODE_GetTokenTxnList2, "Get Token Symbol Fail: "+err.Error())
		}
	} else {
		symbol = p.Token
	}
	txnList, page, err := utiles.GetTokenTxnList(symbol, p.Page)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenTxnList3, "Get Token Txn List Fail: "+err.Error())
	}
	return ResultWithPageReturns(id, txnList, page)
}
func getHolderTxnList(id int, params string) *types.JSONRPCResponse {

	p := types.HolderTokenListParams{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetHolderTokenList1, "Token Holder Txn List Params Error: "+err.Error())
	}
	if p.Token == "" || p.Holder == "" {
		return ErrorReturns(id, ERRCODE_GetHolderTokenList2, "Token/Holder is NULL ")
	}

	if (p.Token[:2] != "0x" && p.Token[:2] != "0X") || (p.Holder[:2] != "0x" && p.Holder[:2] != "0X") {
		return ErrorReturns(id, ERRCODE_GetHolderTokenList3, "Token/Holder Address Invalid")
	}
	var symbol string
	if symbol, err = utiles.GetTokenNameByAddress(p.Token); err != nil {
		return ErrorReturns(id, ERRCODE_GetHolderTokenList4, "Get Token Symbol Fail: "+err.Error())
	}
	txnList, page, err := utiles.GetTokenHolderTxnList(symbol, p.Holder, p.Page)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenTxnList3, "Get Token Holder Txn List Fail: "+err.Error())
	}
	return ResultWithPageReturns(id, txnList, page)

}
func getTokenHolderList(id int, params string) *types.JSONRPCResponse {

	p := types.HolderTokenListParams{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenHolderList1, "Token Holder List Params Error: "+err.Error())
	}
	if p.Token == "" {
		return ErrorReturns(id, ERRCODE_GetTokenHolderList2, "Token is NULL ")
	}
	var symbol string
	if p.Token[:2] == "0x" || p.Token[:2] == "0X" {
		if symbol, err = utiles.GetTokenNameByAddress(p.Token); err != nil {
			return ErrorReturns(id, ERRCODE_GetTokenHolderList3, "Get Token Symbol Fail: "+err.Error())
		}
	} else {
		symbol = p.Token
	}
	holderList, page, err := utiles.GetTokenHolderList(symbol, p.Page)
	if err != nil {
		return ErrorReturns(id, ERRCODE_GetTokenHolderList4, "Get Token Holder List Fail: "+err.Error())
	}
	return ResultWithPageReturns(id, holderList, page)
}

//fd
func getHolderBalance(id int, params string) *types.JSONRPCResponse {
	p := types.HolderTokenParams{}
	err := json.Unmarshal([]byte(params), &p)
	if err != nil {
		return ErrorReturns(id, ERRCODE_HolderBalance1, "Token Holder Balance Params Error: "+err.Error())
	}
	if p.Token == "" || p.Holder == "" {
		return ErrorReturns(id, ERRCODE_HolderBalance2, "Token/Holder is NULL ")
	}
	var balance string
	if (p.Token[:2] != "0x" && p.Token[:2] != "0X") || (p.Holder[:2] != "0x" && p.Holder[:2] != "0X") {
		return ErrorReturns(id, ERRCODE_HolderBalance3, "Token/Holder Address Invalid")
	}

	if symbol, err := utiles.GetTokenNameByAddress(p.Token); err != nil {
		return ErrorReturns(id, ERRCODE_HolderBalance4, "Get Token Symbol Fail: "+err.Error())
	} else {
		if balance, err = utiles.GetHolderBalance(symbol, p.Holder); err != nil {
			return ErrorReturns(id, ERRCODE_HolderBalance5, "Get Token Balance Fail: "+err.Error())
		}
	}
	return ResultNoPageReturns(id, balance)
}

//Handle 分发
func Handle(req *types.JSONRPCRequest) *types.JSONRPCResponse {

	switch req.Method {
	case "get_tokenInfo": //指定token的信息
		return getTokenInfo(req.ID, req.Params)
	case "get_tokenTxnList": //指定token的所有交易列表
		return getTokenTxnList(req.ID, req.Params)
	case "get_tokenHolderList": //指定token的所有持有者列表，按照余额大小排序
		return getTokenHolderList(req.ID, req.Params)
	case "get_holderTxnList": //holder在指定token中的交易记录
		return getHolderTxnList(req.ID, req.Params)
	case "get_holderBalance": //holder在指定token中的余额
		return getHolderBalance(req.ID, req.Params)
	case "token_register": //注册token
		return tokenRegister(req.ID, req.Params)
	case "is_contract": //注册token
		return isContract(req.ID, req.Params)
	}
	return ErrorReturns(req.ID, ERRCODE_Handle, "Unkown Method: "+req.Method)
}

//Route Route
func Route(ctx *web.Context) string {

	ctx.ContentType("json")
	ctx.SetHeader("Access-Control-Allow-Origin", "*", true)
	req, err := praseRequest(ctx)
	if err != nil {
		log.Println("praseRequest Fial: ", err.Error())
		data, _ := json.Marshal(ErrorReturns(req.ID, ERRCODE_Route, "Unkown Request Body"))
		return string(data)
	}
	bs, _ := json.MarshalIndent(req, "", "    ")
	log.Println("RPC Req Data: ", string(bs))
	data, _ := json.MarshalIndent(Handle(req), "", "    ")
	return string(data)
}
