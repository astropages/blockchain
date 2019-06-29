package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//ProofOfWork 工作量证明
type ProofOfWork struct {
	block  *Block   //区块
	target *big.Int //目标值(大数值类型)：与生成的哈希值比较
}

//NewProofOfWork 创建一个工作证明(用户提供区块）系统提供目标值
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	//难度值：
	targetStr := "0001000000000000000000000000000000000000000000000000000000000000"
	//目标值(64位的16进制数)：
	tmpBigInt := new(big.Int)          //创建一个BigInt
	tmpBigInt.SetString(targetStr, 16) //将难度值字符串以16进制赋值给BigInt
	pow.target = tmpBigInt

	return &pow
}

//Run 挖矿（工作量证明）方法：挖矿寻找Nonce,直到随机数+区块数据的sha256值小于难度目标值
func (pow *ProofOfWork) Run() ([]byte, uint64) {
	//定义随机数
	var nonce uint64
	//定以哈希值
	var hash [32]byte

	fmt.Println("开始挖矿...")
	//挖矿
	for {

		//拼接字符串(随机数+区块数据)
		data := pow.PrepareData(nonce)
		//计算哈希值
		hash = sha256.Sum256(data)

		//将哈希值转换为bigInt以进行比较
		tmpInt := new(big.Int)
		tmpInt.SetBytes(hash[:]) //将字符切片转换为BigInt

		//哈希值与难度值比较(返回-1表示x<y，挖矿成功)
		if tmpInt.Cmp(pow.target) == -1 {
			break
		}

		nonce++

	}

	//返回挖矿成功的哈希值和随机数
	return hash[:], nonce
}

//PrepareData 拼接Nonce和区块数据
func (pow *ProofOfWork) PrepareData(nonce uint64) []byte {
	b := pow.block
	//将区块各个字段的字节流进行拼接
	tmp := [][]byte{
		UintToByteSlice(b.Version),
		b.PrevHash,
		b.MerKleRoot, //由所有交易数据计算的哈希值
		UintToByteSlice(b.TimeStamp),
		UintToByteSlice(b.Bits),
		UintToByteSlice(nonce), //随机数
		// b.Hash, //计算后才进行赋值，因此不能参与哈希计算

	}
	//将二维切片使用空切片进行拼接，得到一个切片
	data := bytes.Join(tmp, []byte{})

	return data
}

//IsValid 工作量验证：校验挖矿结果(对求出来的哈希和随机数进行验证)
func (pow *ProofOfWork) IsValid() bool {

	//获取拼接后的数据
	data := pow.PrepareData(pow.block.Nonce)
	//计算哈希值
	hash := sha256.Sum256(data)
	//与难度值比较
	tmpInt := new(big.Int)
	tmpInt.SetBytes(hash[:])
	//返回比较结果
	return tmpInt.Cmp(pow.target) == -1
}
