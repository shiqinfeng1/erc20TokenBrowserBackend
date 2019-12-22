# erc20TokenBrowserBackend

## 启动交易同步
```
./erc20bb  -wtcnode 192.168.50.184:8545 -dbserver 49.51.138.248:3306 -reset
./erc20bb  -reset
```

+ -reset： 清空数据库中所有数据. 默认不清空数据库。
+ -wtcnode：指定wtc链RPC访问节点地址。
+ -dbserver：mysql的访问地址。

### 测试token地址
```
0x668df218d073f413ed2fcea0d48cfbfd59c030ae
```

## web服务
```
./erc20tb -dbserver 49.51.138.248:3306 
```

+ 判断地址是否为注册的token
"data"字段 true=是 或 false=否
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"is_contract","params":"0x668df218d073f413ed2fcea0d48cfbfd59c030ae","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": true
    }
}

```

+ 注册token

```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"token_register","params":"0x668df218d073f413ed2fcea0d48cfbfd59c030ae","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": "Register OK"
    }
}

```
+ 查询所有token信息
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_tokenInfo","params":"","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": [
            {
                "address": "0x668df218d073f413ed2fcea0d48cfbfd59c030ae",
                "name": "WTA",
                "totalSupply": "500000000000000000000000000",
                "decimals": "18"
            }
        ]
    }
}

```

+ 查询指定token地址和指定账户的交易列表
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_holderTxnList","params":"{\"token\":\"0x668df218d073f413ed2fcea0d48cfbfd59c030ae\",\"holder\":\"0x7f125ec7988a5fb35692ff911be45a4f9c48fc48\",\"page\":{\"current_page\":0,\"per_page\":50}}","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": [
            {
                "blockNumber": 487220,
                "blockHash": "0xe2a9b8158f38c3fbb3e395d01eea38239171be697aa03c9036e7c829914d260b",
                "timestamp": "0x5da84801",
                "transferHash": "0x18ea89f965ae66ad814a4491500c3235c47c7841d07c411fb8821f6ac083f273",
                "sender": "0x7f125ec7988a5fb35692ff911be45a4f9c48fc48",
                "receiver": "0xdaf04e95357baef61cd4932241d681fcf0e9270c",
                "value": "1000000000000000000"
            },
            {
                "blockNumber": 484198,
                "blockHash": "0x6fbc440d28faf42c8e859e2e1ae29ebf5e144f44a4d01a1dd4dc5f4d690ce51d",
                "timestamp": "0x5da6e4b5",
                "transferHash": "0x0c3d247c27f5f4694a483e69d4b942909547ce9b35b59343829246b295bf2e3f",
                "sender": "0x5a37e535a430a9a5b3da17c7e68c3647035bd7bd",
                "receiver": "0x7f125ec7988a5fb35692ff911be45a4f9c48fc48",
                "value": "10000000000000000000"
            }
        ],
        "page": {
            "current_page": 0,
            "per_page": 50,
            "total": 2
        }
    }
}

```

+ 查询指定token地址和账户的余额
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_holderBalance","params":"{\"token\":\"0x668df218d073f413ed2fcea0d48cfbfd59c030ae\",\"holder\":\"0x7f125ec7988a5fb35692ff911be45a4f9c48fc48\"}","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": "9000000000000000000"
    }
}
```

+ 查询指定token地址的信息
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_tokenInfo","params":"0x668df218d073f413ed2fcea0d48cfbfd59c030ae","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": {
            "address": "0x668df218d073f413ed2fcea0d48cfbfd59c030ae",
            "name": "WTA",
            "totalSupply": "500000000000000000000000000",
            "decimals": "18"
        }
    }
}
```

+ 查询指定token的所有交易列表

参数可以为合约地址， 也可以为token符号，例如："WTA"
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_tokenTxnList","params":"{\"token\":\"0x668df218d073f413ed2fcea0d48cfbfd59c030ae\",\"page\":{\"current_page\":0,\"per_page\":50}}","id":1}'

{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": [
            {
                "blockNumber": 487220,
                "blockHash": "0xe2a9b8158f38c3fbb3e395d01eea38239171be697aa03c9036e7c829914d260b",
                "timestamp": "0x5da84801",
                "transferHash": "0x18ea89f965ae66ad814a4491500c3235c47c7841d07c411fb8821f6ac083f273",
                "sender": "0x7f125ec7988a5fb35692ff911be45a4f9c48fc48",
                "receiver": "0xdaf04e95357baef61cd4932241d681fcf0e9270c",
                "value": "1000000000000000000"
            },
            {
                "blockNumber": 484198,
                "blockHash": "0x6fbc440d28faf42c8e859e2e1ae29ebf5e144f44a4d01a1dd4dc5f4d690ce51d",
                "timestamp": "0x5da6e4b5",
                "transferHash": "0x0c3d247c27f5f4694a483e69d4b942909547ce9b35b59343829246b295bf2e3f",
                "sender": "0x5a37e535a430a9a5b3da17c7e68c3647035bd7bd",
                "receiver": "0x7f125ec7988a5fb35692ff911be45a4f9c48fc48",
                "value": "10000000000000000000"
            },
            {
                "blockNumber": 472560,
                "blockHash": "0x6a1b162d0e58672d1a0cf5d5e4b5cf2618619612c1b7b99d24901e716e1bce9f",
                "timestamp": "0x5da19a82",
                "transferHash": "0x13d0c1534b776431053e5e564c91bb92eb88b7ba19cc0c3f5381f8fec257f716",
                "sender": "0x5a37e535a430a9a5b3da17c7e68c3647035bd7bd",
                "receiver": "0x957ecfc82e72d5e49be43e882dbe8db6e3f1b49e",
                "value": "10000000000000000000"
            }
        ],
        "page": {
            "current_page": 0,
            "per_page": 50,
            "total": 3
        }
    }
}

```

+ 查询指定token的所有持有者列表， 按照余额大小排序

参数可以为合约地址， 也可以为token符号，例如："WTA"
```
curl -H "Content-Type: application/json" -X POST 'http://localhost:8090' --data '{"jsonrpc":"2.0","method":"get_tokenHolderList","params":"{\"token\":\"WTA\",\"page\":{\"current_page\":0,\"per_page\":50}}","id":1}'
{
    "id": 1,
    "jsonrpc": "2.0",
    "result": {
        "errcode": 0,
        "errmsg": "",
        "data": [
            {
                "balance": 10000000000000000000,
                "address": "0x957ecfc82e72d5e49be43e882dbe8db6e3f1b49e"
            },
            {
                "balance": 9000000000000000000,
                "address": "0x7f125ec7988a5fb35692ff911be45a4f9c48fc48"
            },
            {
                "balance": 4204666696842084352,
                "address": "0x5a37e535a430a9a5b3da17c7e68c3647035bd7bd"
            },
            {
                "balance": 1000000000000000000,
                "address": "0xdaf04e95357baef61cd4932241d681fcf0e9270c"
            }
        ],
        "page": {
            "current_page": 0,
            "per_page": 50,
            "total": 4
        }
    }
}
```
