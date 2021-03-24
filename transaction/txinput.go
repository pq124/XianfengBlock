package transaction

/**
 * 定义交易输入的结构体
 */
type TxInput struct {
	TxId      [32]byte //该字段确定引用自哪笔交易
	Vout      int      //该字段确定引用自该交易的哪个输出
	ScriptSig []byte   //该字段表示使用交易输出的证明，解锁脚本
}
