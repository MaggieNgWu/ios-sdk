/*
 * Copyright 2014-2020. [fisco-dev]
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 *  except in compliance with the License. You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software distributed under the
 *  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 *  express or implied. See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import "C"

// #include <stdio.h>
// #include <stdlib.h>
// #include "data_struct.h"
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"strings"
)

type Params struct {
	ValueType string `json:"type"`
	Value     string `json:"value"`
}

func main() {
}

var config, group, path, endpoint string
var clientSdk *client.Client

// Make config
//export setConfig
func setConfig(basePath *C.char, groupId *C.char, ipPort *C.char, keyfile *C.char) *C.char {
	path = C.GoString(basePath)
	group = C.GoString(groupId)
	endpoint = C.GoString(ipPort)
	config = "[Network]\n" +
		"Type=\"channel\"\n" +
		"CAFile=\"" + path + "/ca.crt\"\n" +
		"Cert=\"" + path + "/sdk.crt\"\n" +
		"Key=\"" + path + "/sdk.key\"\n" +
		"[[Network.Connection]]\n" +
		"NodeURL=\"" + endpoint + "\"\n" +
		"GroupID=" + group + "\n\n" +
		"[Account]\n" +
		"KeyFile= \"" + path + "/" + C.GoString(keyfile) + "\"\n\n" +
		"[Chain]\n" +
		"ChainID=1\n" +
		"SMCrypto=false"
	// connect node
	configs, err := conf.ParseConfig([]byte(config))
	if err != nil {
		return C.CString(err.Error())
	}
	clientSdk, err = client.Dial(&configs[0])
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("Connect Success to node " + endpoint)
}

// Node RPC call
//export getClientVersion
func getClientVersion() *C.RPCResult {
	cv, err := clientSdk.GetClientVersion(context.Background())
	if err != nil {
		return ToRPCResult(cv, err)
	}
	return ToRPCResult(cv, err)
}

//export deployContract
func deployContract(abiString *C.char, binString *C.char) *C.Transaction {
	ops := clientSdk.GetTransactOpts()
	parsed, err := abi.JSON(strings.NewReader(C.GoString(abiString)))
	if err != nil {
		return ToCTransaction(common.Address{}, nil, err)
	}
	addr, transaction, _, err := bind.DeployContract(ops, parsed, common.FromHex(C.GoString(binString)), clientSdk)
	return ToCTransaction(addr, transaction, err)
}

//export sendTransaction
func sendTransaction(abiString *C.char, address *C.char, method *C.char, params *C.char) *C.TxExecResult {
	parsed, err := abi.JSON(strings.NewReader(C.GoString(abiString)))
	if err != nil {
		return ToTxExecutionResult(nil, nil, err)
	}
	goParams, err := toGoParams(C.GoString(params))
	if err != nil {
		return ToTxExecutionResult(nil, nil, err)
	}
	addString := C.GoString(address)
	addr := common.HexToAddress(addString)
	c := bind.NewBoundContract(addr, parsed, clientSdk, clientSdk, clientSdk)
	var tx *types.Transaction
	var receipt *types.Receipt
	if len(goParams) == 0 {
		tx, receipt, err = c.Transact(clientSdk.GetTransactOpts(), C.GoString(method))
		if err != nil {
			return ToTxExecutionResult(nil, nil, err)
		}
	} else {
		tx, receipt, err = c.Transact(clientSdk.GetTransactOpts(), C.GoString(method), goParams...)
		if err != nil {
			return ToTxExecutionResult(nil, nil, err)
		}
	}
	if err != nil {
		return ToTxExecutionResult(nil, nil, err)
	}
	return ToTxExecutionResult(tx, receipt, err)
}

//export call
func call(abiString *C.char, address *C.char, method *C.char, params *C.char) *C.CallResult {
	parsed, err := abi.JSON(strings.NewReader(C.GoString(abiString)))
	if err != nil {
		return ToCallResult(nil, err)
	}
	goParams, err := toGoParams(C.GoString(params))
	if err != nil {
		return ToCallResult(nil, err)
	}
	addr := common.HexToAddress(C.GoString(address))
	c := bind.NewBoundContract(addr, parsed, clientSdk, clientSdk, clientSdk)
	var result interface{}
	if len(goParams) == 0 {
		err = c.Call(clientSdk.GetCallOpts(), &result, C.GoString(method))
	} else {
		err = c.Call(clientSdk.GetCallOpts(), result, C.GoString(method), goParams...)
	}
	if err != nil {
		return ToCallResult(nil, err)
	}
	resultBytes, err := json.MarshalIndent(result, "", "\t")
	return ToCallResult(resultBytes, err)
}

//export getBlockByHash
func getBlockByHash(bHash *C.char, includeTx *C.char) *C.RPCResult {
	bhashString := C.GoString(bHash)
	includetxBool, err := strconv.ParseBool(C.GoString(includeTx))
	if err != nil {
		return ToRPCResult(nil, err)
	}
	raw, err := clientSdk.GetBlockByHash(context.Background(), bhashString, includetxBool)
	return ToRPCResult(raw, err)
}

//export getBlockByNumber
func getBlockByNumber(bNum *C.char, includeTx *C.char) *C.RPCResult {
	bnumString := C.GoString(bNum)
	includetxBool, err := strconv.ParseBool(C.GoString(includeTx))
	if err != nil {
		return ToRPCResult(nil, err)
	}
	raw, err := clientSdk.GetBlockByNumber(context.Background(), bnumString, includetxBool)
	return ToRPCResult(raw, err)
}

//export getBlockHashByNumber
func getBlockHashByNumber(bNum *C.char) *C.RPCResult {
	bnumString := C.GoString(bNum)
	raw, err := clientSdk.GetBlockHashByNumber(context.Background(), bnumString)
	return ToRPCResult(raw, err)
}

//export getTransactionByHash
func getTransactionByHash(bHash *C.char) *C.RPCResult {
	bhashString := C.GoString(bHash)
	raw, err := clientSdk.GetTransactionByHash(context.Background(), bhashString)
	return ToRPCResult(raw, err)
}

//export getTransactionByBlockHashAndIndex
func getTransactionByBlockHashAndIndex(bHash *C.char, txIndex *C.char) *C.RPCResult {
	bhashString := C.GoString(bHash)
	txindexString := C.GoString(txIndex)
	raw, err := clientSdk.GetTransactionByBlockHashAndIndex(context.Background(), bhashString, txindexString)
	return ToRPCResult(raw, err)
}

//export getTransactionByBlockNumberAndIndex
func getTransactionByBlockNumberAndIndex(bNum *C.char, txIndex *C.char) *C.RPCResult {
	bnumString := C.GoString(bNum)
	txindexString := C.GoString(txIndex)
	raw, err := clientSdk.GetTransactionByBlockNumberAndIndex(context.Background(), bnumString, txindexString)
	return ToRPCResult(raw, err)
}

//export getTransactionReceipt
func getTransactionReceipt(txHash *C.char) *C.RPCReceiptResult {
	txHashString := C.GoString(txHash)
	raw, err := clientSdk.GetTransactionReceipt(context.Background(), txHashString)
	return ToRPCReceiptResult(raw, err)
}

//export getContractAddress
func getContractAddress(txHash *C.char) *C.RPCResult {
	txHashString := C.GoString(txHash)
	raw, err := clientSdk.GetContractAddress(context.Background(), txHashString)
	var emptAddr common.Address
	if raw != emptAddr {
		return ToRPCResult([]byte(raw.Hex()), err)
	} else {
		return ToRPCResult([]byte(""), err)
	}
}

//export getPendingTransactions
func getPendingTransactions() *C.RPCResult {
	raw, err := clientSdk.GetPendingTransactions(context.Background())
	return ToRPCResult(raw, err)
}

//export getPendingTxSize
func getPendingTxSize() *C.RPCResult {
	raw, err := clientSdk.GetPendingTxSize(context.Background())
	return ToRPCResult(raw, err)
}

//export getCode
func getCode(addr *C.char) *C.RPCResult {
	addrString := C.GoString(addr)
	raw, err := clientSdk.GetCode(context.Background(), addrString)
	return ToRPCResult(raw, err)
}

//export getTotalTransactionCount
func getTotalTransactionCount() *C.RPCResult {
	raw, err := clientSdk.GetTotalTransactionCount(context.Background())
	return ToRPCResult(raw, err)
}

//export getSystemConfigByKey
func getSystemConfigByKey(configKey *C.char) *C.RPCResult {
	configKeyString := C.GoString(configKey)
	raw, err := clientSdk.GetSystemConfigByKey(context.Background(), configKeyString)
	return ToRPCResult(raw, err)
}

//export getPBFTView
func getPBFTView() *C.RPCResult {
	raw, err := clientSdk.GetPBFTView(context.Background())
	return ToRPCResult(raw, err)
}

//export getSealerList
func getSealerList() *C.RPCResult {
	raw, err := clientSdk.GetSealerList(context.Background())
	return ToRPCResult(raw, err)
}

//export getObserverList
func getObserverList() *C.RPCResult {
	raw, err := clientSdk.GetObserverList(context.Background())
	return ToRPCResult(raw, err)
}

//export getConsensusStatus
func getConsensusStatus() *C.RPCResult {
	raw, err := clientSdk.GetConsensusStatus(context.Background())
	return ToRPCResult(raw, err)
}

//export getSyncStatus
func getSyncStatus() *C.RPCResult {
	raw, err := clientSdk.GetSyncStatus(context.Background())
	return ToRPCResult(raw, err)
}

//export getPeers
func getPeers() *C.RPCResult {
	raw, err := clientSdk.GetPeers(context.Background())
	return ToRPCResult(raw, err)
}

//export getGroupPeers
func getGroupPeers() *C.RPCResult {
	raw, err := clientSdk.GetGroupPeers(context.Background())
	return ToRPCResult(raw, err)
}

//export getNodeIDList
func getNodeIDList() *C.RPCResult {
	raw, err := clientSdk.GetNodeIDList(context.Background())
	return ToRPCResult(raw, err)
}

//export getGroupList
func getGroupList() *C.RPCResult {
	raw, err := clientSdk.GetGroupList(context.Background())
	return ToRPCResult(raw, err)
}

//export getBlockLimit
func getBlockLimit() *C.RPCResult {
	raw, err := clientSdk.GetBlockLimit(context.Background())
	if err != nil {
		return ToRPCResult([]byte(""), err)
	} else {
		return ToRPCResult([]byte(raw.String()), err)
	}
}

//export getBlockNumber
func getBlockNumber() *C.RPCResult {
	raw, err := clientSdk.GetBlockNumber(context.Background())
	return ToRPCResult(raw, err)
}

//export getChainID
func getChainID() *C.RPCResult {
	raw, err := clientSdk.GetChainID(context.Background())
	if err != nil {
		return ToRPCResult([]byte(""), err)
	} else {
		return ToRPCResult([]byte(raw.String()), err)
	}
}

func toGoParams(param string) ([]interface{}, error) {
	var objs []Params
	if err := json.Unmarshal([]byte(param), &objs); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var par []interface{}
	for _, t := range objs {
		switch t.ValueType {
		case "int":
			i, err := strconv.ParseInt(t.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			par = append(par, i)
		case "uint":
			i, err := strconv.ParseUint(t.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			par = append(par, i)
		case "bool":
			b, err := strconv.ParseBool(t.Value)
			if err != nil {
				return nil, err
			}
			par = append(par, b)
		case "string":
			par = append(par, t.Value)
		case "address":
			addr := common.HexToAddress(t.Value)
			par = append(par, addr)
		case "bytes":
			par = append(par, []byte(t.Value))
		}
	}
	return par, nil
}
