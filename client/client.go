package client

import (
	"fmt"
	"XianfengChain04/chain"
	"os"
	"flag"
	"math/big"
	"XianfengChain04/utils"
)

/**
 * 该结构体定义了用于实现命令行参数解析的结构体
 */
type CmdClient struct {
	Chain chain.BlockChain
}

/**
 * client运行方法
 */
func (cmd *CmdClient) Run() {
	args := os.Args
	//1、处理用户没有输入任何命令和参数的情况，打印输出说明书
	if len(args) == 1 {
		cmd.Help()
		return
	}
	//2、解析用户输入的第一个参数，作为功能命令进行解析
	switch os.Args[1] {
	case GENERATEGENSIS:
		cmd.GenerateGensis()
	case SENDTRANSACTION: //发送交易..（前提：创世区块已存在）
		cmd.SendTransaction()
	case GETBALANCE: //获取某个地址的余额
		cmd.GetBalance()
	case GETLASTBLOCK:
		cmd.GetLastBlock()
	case GETALLBLOCKS:
		cmd.GetAllBlocks()
	case GETNEWADDRESS: //生成新地址的功能
		cmd.GetNewAddress()
	case LISTADDRESS: //获取所有的地址列表
		cmd.ListAddress()
	case DUMPPRIVKEY:
		cmd.DumpPrivKey()
	case HELP:
		cmd.Help()
	default:
		cmd.Default()
	}
}

func (cmd *CmdClient) Default() {
	fmt.Println("go run main.go：Unknown subcommand.")
	fmt.Println("Run 'go run main.go help' for usage.")
}

func (cmd *CmdClient)DumpPrivKey()  {
	dumpPrivkey:=flag.NewFlagSet(DUMPPRIVKEY,flag.ExitOnError)
	address :=dumpPrivkey.String("address","","要导出的私钥地址")
	dumpPrivkey.Parse(os.Args[2:])
	if len(os.Args[2:]) > 2 {
		fmt.Println("无法解析输入参数,请检查后重试!")
		return
	}
	pri,err := cmd.Chain.DumpPrivkey(*address)
	if err!=nil {
		fmt.Println(err.Error())
		return
	}
    fmt.Printf("私钥是%x",pri.D.Bytes())
}

func (cmd *CmdClient) ListAddress() {
	listAddress := flag.NewFlagSet(LISTADDRESS, flag.ExitOnError)
	listAddress.Parse(os.Args[2:])
	if len(os.Args[2:]) > 0 {
		fmt.Println("无法解析参数，请检查后重试！")
		return
	}
	addList, err := cmd.Chain.GetAddressList()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("获取地址列表成功，地址信息如下：")
	for index, add := range addList {
		fmt.Printf("[%d]:%s\n", index+1, add)
	}
}

/**
 * 定义新的方法：用于生成新的地址
 */
func (cmd *CmdClient) GetNewAddress() {
	getNewAddress := flag.NewFlagSet(GETNEWADDRESS, flag.ExitOnError)
	getNewAddress.Parse(os.Args[2:])

	if len(os.Args[2:]) > 0 {
		fmt.Println("抱歉，生成新地址功能无法解析参数，请重试！")
		return
	}

	address, err := cmd.Chain.GetNewAddress()
	if err != nil {
		fmt.Println("生成地址遇到错误：", err.Error())
		return
	}
	fmt.Println("生成新的地址：", address)
}

func (cmd *CmdClient) GetAllBlocks() {
	blocks, err := cmd.Chain.GetAllBlocks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("恭喜，查询到所有区块数据")
	for _, block := range blocks {
		fmt.Printf("区块高度:%d,区块哈希:%x\n", block.Height, block.Hash)
		fmt.Print("区块中的交易信息：\n")
		for index, tx := range block.Transactions {
			fmt.Printf("   第%d笔交易,交易hash:%x\n", index, tx.TxHash)
			for inputIndex, input := range tx.Inputs {
				fmt.Printf("       第%d笔交易输入,%s花了%x的%d的钱\n", inputIndex, input.ScriptSig, input.TxId, input.Vout)
			}
			for outputIndex, output := range tx.Outputs {
				fmt.Printf("       第%d笔交易输出,%s实现收入%f\n", outputIndex, output.ScriptPub, output.Value)
			}
		}
		fmt.Println()
	}
}

func (cmd *CmdClient) GetLastBlock() {
	lastBlock := cmd.Chain.GetLastBlock()
	//1、判断是否为空
	hashBig := new(big.Int)
	hashBig.SetBytes(lastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有最新区块
		fmt.Println("抱歉，当前暂无最新区块.")
		return
	}
	fmt.Println("恭喜，获取到最新区块数据")
	fmt.Printf("最新区块高度:%d\n", lastBlock.Height)
	fmt.Printf("最新区块哈希:%x\n", lastBlock.Hash)

	for index, tx := range lastBlock.Transactions {
		fmt.Printf("区块交易%d,交易:%v\n", index, tx)
	}
}

/**
 * 用户发起交易
 */
func (cmd *CmdClient) SendTransaction() {
	//-data
	createBlock := flag.NewFlagSet(SENDTRANSACTION, flag.ExitOnError)
	from := createBlock.String("from", "", "交易发起人地址")
	to := createBlock.String("to", "", "交易接收者地址")
	amount := createBlock.String("amount", "", "转账的数量")

	if len(os.Args[2:]) > 6 {
		fmt.Println("sendTransaction命令只支持三个参数和参数值，请重试")
		return
	}
	createBlock.Parse(os.Args[2:])

	//from，to，amount三个参数是字符串类型，同时需要满足符合JSON格式
	fromSlice, err := utils.JSONArray2String(*from)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	toSlice, err := utils.JSONArray2String(*to)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}
	amountSlice, err := utils.JSONArray2Float(*amount)
	if err != nil {
		fmt.Println("抱歉，参数格式不正确，请检查后重试！")
		return
	}

	//先看看参数个数是否一致
	fromLen := len(fromSlice)
	toLen := len(toSlice)
	amountLen := len(amountSlice)
	if fromLen != toLen || fromLen != amountLen || toLen != amountLen {
		fmt.Println("参数个数不一致，请检查参数后重试")
		return
	}

	//1、先判断是否已生成创世区块，如果没有创世区块，提示用户先生成
	//[0000000]
	hashBig := new(big.Int)
	hashBig.SetBytes(cmd.Chain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有创世区块
		fmt.Println("That not a gensis block in blockchain，please use go run main.go generategensis command to create a gensis block first.")
		return
	}

	err = cmd.Chain.SendTransaction(fromSlice, toSlice, amountSlice)
	if err != nil {
		fmt.Println("抱歉，发送交易出现错误：", err.Error())
		return
	}
	fmt.Println("交易发送成功")
}

/**
 * 获取地址的余额方法
 */
func (cmd *CmdClient) GetBalance() {
	getbalance := flag.NewFlagSet(GETBALANCE, flag.ExitOnError)
	var addr string
	getbalance.StringVar(&addr, "address", "", "用户的地址")
	getbalance.Parse(os.Args[2:])

	blockChain := cmd.Chain
	//1、先判断是否有创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 0 { //没有创世区块
		fmt.Println("抱歉，该网络链暂未存在，无法查询")
		return
	}
	//2、调用余额查询功能
	balance, err := blockChain.GetBalance(addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("地址%s的余额是：%f\n", addr, balance)
}

func (cmd *CmdClient) GenerateGensis() {
	//命令参数集合
	generategensis := flag.NewFlagSet(GENERATEGENSIS, flag.ExitOnError)
	//解析参数
	var addr string
	generategensis.StringVar(&addr, "address", "", "用户指定的矿工的地址")
	generategensis.Parse(os.Args[2:])

	blockChain := cmd.Chain
	//1、先判断该blockchain中是否已存在创世区块
	hashBig := new(big.Int)
	hashBig.SetBytes(blockChain.LastBlock.Hash[:])
	if hashBig.Cmp(big.NewInt(0)) == 1 {
		fmt.Println("创世区块已存在，不能重复生成创世区块")
		return
	}

	err := blockChain.CreateCoinBase(addr)
	if err != nil {
		fmt.Println("抱歉，创建coinbase交易失败，遇到错误：", err.Error())
		return
	}
	fmt.Println("恭喜！生成了一笔COINBASE交易，奖励已到账。")
}

/**
 * 该方法用于打印输出项目的使用和说明信息，相当于项目的帮助文档和说明书
 */
func (cmd *CmdClient) Help() {
	fmt.Println("------------Welcome to XianfengChain04 Project-----------")
	fmt.Println("XianfengChain04 is a custom blockchain project, the project plan to build a very simple public chain.")
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println()
	fmt.Println("go run main.go command [arguments]")
	fmt.Println()
	fmt.Println("AVAILABLE COMMANDS")
	fmt.Println()
	fmt.Println("    generategensis    use the command can create a genesis block and save to the boltdb file. use the genesis argument to set the custom data.")
	fmt.Println("    sendtransaction   this command used to send a new transaction, that can specified three argument named from, to and amount.")
	fmt.Println("    getbalance        this is a command that can get the balance of specified address")
	fmt.Println("    getlastblock      get the lastest block data.")
	fmt.Println("    getallblocks      return all blocks data to user.")
	fmt.Println("    getnewaddress     this commadn used to create a new address by bitcoin algorithm")
	fmt.Println("    help              use the command can print usage infomation.")
	fmt.Println()
	fmt.Println("Use go run main.go help [command] for more information about a command.")
}
