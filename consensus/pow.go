package consensus

import (
	"XianfengChain04/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

//目的：拿到区块的属性数据(属性值)
//1、通过结构体引用，引用block结构体，然后访问其属性，比如block.Height
//2、接口

const DIFFICULTY = 16 //难度值系数

type PoW struct {
	Block  BlockInterface
	Target *big.Int
}

func (pow PoW) FindNonce() ([32]byte, int64) {
	//1、给定一个nonce值，计算区块hash
	var nonce int64
	nonce = 0
	//无限循环
	hashBig := new(big.Int)
	for {
		hash := CalculateHash(pow.Block, nonce)
		//2、拿到系统的目标值
		target := pow.Target
		//3、比较大小
		//target big.Int
		//hash  [32]byte

		hashBig = hashBig.SetBytes(hash[:])
		//result := bytes.Compare(hash[:], target.Bytes())
		result := hashBig.Cmp(target)
		//4、判断结果
		if result == -1 {
			return hash, nonce
		}
		nonce++ //否则nonce自增
	}
}

/**
 * 根据区块已有的信息和当前nonce的赋值，计算区块的hash
 */
func CalculateHash(block BlockInterface, nonce int64) [32]byte {
	heightByte, _ := utils.Int2Byte(block.GetHeight())
	versionByte, _ := utils.Int2Byte(block.GetVersion())
	timeByte, _ := utils.Int2Byte(block.GetTimeStamp())
	nonceByte, _ := utils.Int2Byte(nonce)

	prev := block.GetPrevHash()

	txs := block.GetTransactions()
	txsBytes := make([]byte, 0)
	for _, tx := range txs {
		// struct -> []byte
		txData, err := utils.Encode(tx)
		if err != nil {
			break
		}
		txsBytes = append(txsBytes, txData...)
	}
	blockByte := bytes.Join([][]byte{heightByte,
		versionByte,
		prev[:],
		timeByte,
		nonceByte,
		txsBytes,
	}, []byte{})
	//计算区块的hash
	hash := sha256.Sum256(blockByte)
	return hash
}
