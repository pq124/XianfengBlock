package client

const (
	GENERATEGENSIS  = "generategensis"  //ccoinbase -addr
	SENDTRANSACTION = "sendtransaction" //sendTransaction from to amount
	GETBALANCE      = "getbalance"      //获取地址的余额功能
	GETLASTBLOCK    = "getlastblock"
	GETALLBLOCKS    = "getallblocks"
	GETNEWADDRESS   = "getnewaddress" //生成新的比特币地址
	DUMPPRIVKEY     = "dumpprivkey"
	LISTADDRESS     = "listaddress"   //列出所有目前已经生成并管理的地址
	HELP            = "help"
)
