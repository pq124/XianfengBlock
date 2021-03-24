package transaction

import (
	"XianfengChain04/utils"
	"crypto/sha256"
)

const REWARDSIZE = 50

/**
 * 定义交易的结构体
 */
type Transaction struct {
	//交易哈希
	TxHash [32]byte
	//交易输入
	Inputs []TxInput
	//交易输出
	Outputs []TxOutput
}

/**
 * 该函数用于定义一个coinbase交易，并返回该交易结构体
 */
func CreateCoinBase(addr string) (*Transaction, error) {
	output0 := TxOutput{
		Value:     REWARDSIZE,
		ScriptPub: []byte(addr),
	}

	coinbase := Transaction{
		Outputs: []TxOutput{output0},
	}
	coinbaseBytes, err := utils.Encode(coinbase)
	if err != nil {
		return nil, err
	}
	coinbase.TxHash = sha256.Sum256(coinbaseBytes)

	return &coinbase, nil
}

/**
 * 该函数用于构建一笔普通的交易，返回构建好的交易实例
 */
func CreateNewTransaction(utxos []UTXO, from string, to string, amount float64) (*Transaction, error) {
	//1、构建inputs
	inputs := make([]TxInput, 0) //用于存放交易输入的容器
	var inputAmount float64      //该变量用于记录转账发起者一共付了多少钱
	//input -> 交易输入:对某个交易的交易输出UTXO的引用
	for _, utxo := range utxos {
		input := TxInput{
			TxId:      utxo.TxId,
			Vout:      utxo.Vout,
			ScriptSig: []byte(from),
		}
		inputAmount += utxo.Value
		//把构建好的input存入到交易输入容器中
		inputs = append(inputs, input)
	}

	//2、构建outputs
	outputs := make([]TxOutput, 0) //用于存放交易输出的容器
	//构建转账接收者的交易输出
	output0 := TxOutput{
		Value:     amount,
		ScriptPub: []byte(to),
	}
	outputs = append(outputs, output0) //把第一个交易输出放入到专门存交易输出的容器中

	//判断是否需要找零,如果需要找零，则需要构建一个新的找零输出
	if inputAmount-amount > 0 {
		output1 := TxOutput{
			Value:     inputAmount - amount,
			ScriptPub: []byte(from),
		}
		outputs = append(outputs, output1)
	}

	//3、构建transaction
	newTransaction := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	//4、计算transaction的哈希,并赋值
	transactionBytes, err := utils.Encode(newTransaction)
	newTransaction.TxHash = sha256.Sum256(transactionBytes)
	if err != nil {
		return nil, err
	}

	//5、将构建的transaction实例进行返回
	return &newTransaction, nil
}
