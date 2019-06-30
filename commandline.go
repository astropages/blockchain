package main

import "fmt"

/*
	命令行方法
*/

//创建区块链
func (cli *CLI) createBlockChain() {
	//创建区块链
	err := CreateBlockChain()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// //添加区块
// func (cli *CLI) addBlock(data string) {
// 	//获取一个区块链实例
// 	bc, err := GetBlockChainInstance()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer bc.db.Close()
// 	//添加区块
// 	if err := bc.AddBlock(data); err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// }

//获取地址对应的金额
func (cli *CLI) getBalance(address string) {
	//获取一个区块链实例
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取地址的utxo详情
	utxoInfos := bc.FindMyUTXO(address)
	//遍历累加金额
	total := 0.0
	for _, utxo := range utxoInfos {
		total += utxo.TXOutput.Value
	}

	fmt.Printf("%s的金额为: %f\n", address, total)
}

//打印区块链
func (cli *CLI) printBlockChain() {
	//获取一个区块链实例
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bc.db.Close()
	//使用迭代器打印区块信息
	it := bc.NewIterator()
	for {
		//使用迭代器Next方法获取区块并移动游标
		block := it.Next()
		//打印区块链
		fmt.Println("===============================")
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("MerKleRoot: %x\n", block.MerKleRoot)
		fmt.Printf("TimeStamp: %d\n", block.TimeStamp)
		fmt.Printf("Bits: %d\n", block.Bits)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Transactions[0].TXInputs[0].ScriptSign)

		//校验区块（工作量验证）
		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v\n", pow.IsValid())

		//如果区块前哈希为空则退出循环
		if block.PrevHash == nil {
			break
		}
	}
}
