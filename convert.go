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

// #include <stdio.h>
// #include <stdlib.h>
// #include "data_struct.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func ToTxExecutionResult(transaction *types.Transaction, receipt *types.Receipt, err error) *C.TxExecResult {
	er := C.calloc(1, C.sizeof_TxExecResult)
	execResult := (*C.TxExecResult)(er)

	if transaction != nil {
		execResult.tx = ToCTransaction(receipt.ContractAddress, transaction, err)
	}

	if receipt != nil {
		execResult.receipt = ToCReceipt(receipt)
	}

	if err != nil {
		execResult.err = C.CString(err.Error())
	}
	return execResult
}

func ToCallResult(result []byte, err error) *C.CallResult {
	cr := C.calloc(1, C.sizeof_CallResult)
	callResult := (*C.CallResult)(cr)
	if err != nil {
		callResult.err = C.CString(err.Error())
	}
	if result != nil {
		callResult.result = C.CString(string(result))
	}
	return callResult
}

func ToCTransaction(address common.Address, transaction *types.Transaction, err error) *C.Transaction {
	t := C.calloc(1, C.sizeof_Transaction)
	tx := (*C.Transaction)(t)
	if err != nil {
		tx.err = C.CString(err.Error())
	} else {
		tx.address = C.CString(address.Hex())
	}
	if transaction != nil {
		//fmt.Println(string(transaction.Data()))
		tx.hash = C.CString(transaction.Hash().Hex())
		tx.data = C.CString(string(transaction.Data()))
		tx.size = C.double(float64(transaction.Size()))
	}
	return tx
}

func ToCReceipt(receipt *types.Receipt) *C.Receipt {
	if receipt == nil {
		return nil
	}
	r := C.calloc(1, C.sizeof_Receipt)
	rec := (*C.Receipt)(r)
	rec.TransactionHash = C.CString(receipt.TransactionHash)
	rec.TransactionIndex = C.CString(receipt.TransactionIndex)
	rec.BlockHash = C.CString(receipt.BlockHash)
	rec.BlockNumber = C.CString(receipt.BlockNumber)
	rec.GasUsed = C.CString(receipt.GasUsed)
	rec.ContractAddress = C.CString(receipt.ContractAddress.Hex())
	rec.Root = C.CString(receipt.Root)
	rec.Status = C.int(receipt.Status)
	rec.From = C.CString(receipt.From)
	rec.To = C.CString(receipt.To)
	rec.Input = C.CString(receipt.Input)
	rec.GasUsed = C.CString(receipt.GasUsed)
	rec.Output = C.CString(receipt.Output)
	logs, err := json.MarshalIndent(receipt.Logs, "", "\t")
	if err != nil {
		fmt.Println("err")
	} else {
		rec.Logs = C.CString(string(logs))
	}
	//rec.Logs = ToCLogs(receipt.Logs)
	rec.LogsBloom = C.CString(receipt.LogsBloom)
	return rec
}

func ToRPCResult(result []byte, err error) *C.RPCResult {
	r := C.calloc(1, C.sizeof_RPCResult)
	rs := (*C.RPCResult)(r)
	if err != nil {
		rs.err = C.CString(err.Error())
	}
	rs.result = C.CString(string(result))
	return rs
}

func ToRPCReceiptResult(receipt *types.Receipt, err error) *C.RPCReceiptResult {
	er := C.calloc(1, C.sizeof_RPCReceiptResult)
	execResult := (*C.RPCReceiptResult)(er)

	if receipt != nil {
		execResult.receipt = ToCReceipt(receipt)
	}

	if err != nil {
		execResult.err = C.CString(err.Error())
	}
	return execResult
}
