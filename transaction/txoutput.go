package transaction

/**
 * 定义交易输出的结构体
 */
type TxOutput struct {
	Value     float64 //转账的数量
	ScriptPub []byte  //锁定脚本
}
