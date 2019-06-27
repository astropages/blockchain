package main

import (
	"fmt"
	"os"
)

//CLI 命令行(Command Line)
type CLI struct {
}

//Usage 使用说明
const Usage = `
Usage:
	create "创建区块链"
	add <data> "添加区块"
	print "打印区块链" 
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
	case "add":
		fmt.Println("添加区块")
		if len(cmds) < 3 {
			fmt.Println("请输入区块数据")
			return
		}
		data := cmds[2]
		cli.addBlock(data)
	case "print":
		fmt.Println("打印区块链")
		cli.printBlockChain()
	default:
		fmt.Println("输入参数错误")
	}
}
