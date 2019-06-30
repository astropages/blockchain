package main

import (
	"fmt"
	"os"
	"strconv"
)

//CLI 命令行(Command Line)
type CLI struct {
}

//Usage 使用说明
const Usage = `
Usage:
	create "创建区块链"
	getbalance <address> "获取地址对应的金额"
	print "打印区块链" 
	send <from> <to> <amount> <miner> <data> "转账：付款人 收款人 转账金额 矿工 数据"
	createwallet "创建钱包"
`

//Run 解析用户输入命令的方法
func (cli *CLI) Run() {

	//获取输入参数
	cmds := os.Args
	if len(cmds) < 2 {
		fmt.Println("输入参数错误")
		fmt.Println(Usage)
		return
	}

	//根据输入参数调用函数
	switch cmds[1] {
	case "create":
		fmt.Println("创建区块链")
		cli.createBlockChain()
	case "print":
		fmt.Println("打印区块链")
		cli.printBlockChain()
	case "getbalance":
		fmt.Println("获取地址金额")
		if len(cmds) < 3 {
			fmt.Println("请输入地址")
			return
		}
		address := cmds[2]
		cli.getBalance(address)
	case "send":
		fmt.Println("转账")
		if len(cmds) < 7 {
			fmt.Println("转账参数错误")
			return
		}
		from := cmds[2]
		to := cmds[3]
		amount, _ := strconv.ParseFloat(cmds[4], 64)
		miner := cmds[5]
		data := cmds[6]
		cli.send(from, to, amount, miner, data)
	case "createwallet":
		fmt.Println("创建钱包")
		cli.createWallet()
	default:
		fmt.Println("输入参数错误")
	}
}
