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
		fmt.Printf("MerkleRoot: %x\n", block.MerkleRoot)
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

//转账：每次转账时便添加一个区块
func (cli *CLI) send(from string, to string, amount float64, miner string, data string) {
	//获取一个区块链实例
	bc, err := GetBlockChainInstance()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bc.db.Close()

	//创建挖矿交易
	coinbaseTX := NewCoinbaseTX(miner, data)

	//创建交易集合，添加有效交易
	txs := []*Transaction{coinbaseTX}

	//创建普通交易
	tx := NewTransaction(from, to, amount, bc)
	if tx != nil { //找到有效交易
		txs = append(txs, tx)
	} else {
		fmt.Println("未找到有效交易")
	}

	//添加区块
	err = bc.AddBlock(txs)
	if err != nil {
		fmt.Println("转账失败")
		return
	}
	fmt.Println("转账成功")
}
