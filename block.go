package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//Block 区块
type Block struct {
	Version      uint64         //版本号
	PrevHash     []byte         //前区块哈希值
	MerkleRoot   []byte         //梅克尔根（交易的根哈希值）
	TimeStamp    uint64         //时间戳
	Bits         uint64         //调整比特币挖矿难度的数值（用于计算哈希）
	Nonce        uint64         //随机数（挖矿时寻找的数值）
	Hash         []byte         //当前区块哈希值
	Transactions []*Transaction //区块数据：区块的交易集合
}

//NewBlock 创建一个区块(传入交易和前区块的哈希)
func NewBlock(txs []*Transaction, prevHash []byte) *Block {
	b := Block{
		Version:      0,
		PrevHash:     prevHash,
		MerkleRoot:   nil,
		TimeStamp:    uint64(time.Now().UnixNano()),
		Bits:         0,
		Nonce:        0,
		Hash:         nil,
		Transactions: txs,
	}
	//填充梅克尔根值
	b.HashTransactionMerkleRoot()

	//工作量证明(挖矿寻找随机数并计算符合难度目标的哈希值)
	pow := NewProofOfWork(&b)
	hash, nonce := pow.Run()
	b.Hash = hash
	b.Nonce = nonce

	//返回区块
	return &b
}

//Serialize 将区块数据序列化为字节流的方法
func (b *Block) Serialize() []byte {
	//定义buffer容器
	var buffer bytes.Buffer
	//创建编码器
	encoder := gob.NewEncoder(&buffer)
	//编码
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//得字节流
	data := buffer.Bytes()
	return data
}

//DeSerialize 将字节流反序列化为区块数据
func DeSerialize(data []byte) *Block {

	//定义一个区块
	var block Block

	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(data))

	//解码
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &block
}

//HashTransactionMerkleRoot 简易的梅克尔根，把所有交易拼接后计算哈希值
func (b *Block) HashTransactionMerkleRoot() {

	var info [][]byte
	//遍历所有交易并计算哈希值
	for _, tx := range b.Transactions {
		//将所有的哈希值拼接后再计算哈希值
		txHash := tx.TXID
		info = append(info, txHash)
	}
	//拼接字符切片
	value := bytes.Join(info, []byte{})
	hash := sha256.Sum256(value)
	//将最终的哈希值赋值给MerKleRoot
	b.MerkleRoot = hash[:]
}
