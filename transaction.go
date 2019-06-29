/*
	比特币中，没有付款人和收款人，只有输入(input)和输出(output)，
	每个输入都对应着之前别人给你转账时产生的某个输出。
	一笔交易中可以有多个输入和多个输出，给自己找零就是给自己生成一个输出。

		输出产生：
			先从张三给李四转账开始说起，张三给李四转账时，比特币系统会生成一个output，这个output里面包括两个东西：
				1. 转的金额，例如100
				2. 一个锁定脚本，使用李四的**公钥哈希**对转账金额1btc进行锁定，可以理解为用公钥哈希加密了。
			真实的锁定脚本
				锁定脚本：给我收款人的地址，我用这个人公钥进行锁定
				解锁脚本：提供支付人的私钥签名（公钥）

		输入产生：
			与output对应的是input结构，每一个input都源自一个output，在李四对王五进行转账时，系统会创建input，为了定位这笔钱的来源，这个input结构包含以下内容：
				1. 在哪一笔交易中，即需要张三->李四这笔转账的交易ID(hash)
				2. 所引用交易的那个output，所以需要一个output的索引(int)
				3. 定位到了这个output，如何证明能支配呢，所以需要一个张三的签名。（解锁脚本，包括签名和自己的公钥）

		未消费输出（UTXO）：
			1. UTXO：unspent transaction output，是比特币交易中最小的支付单元，不可分割，每一个UTXO必须一次性消耗完，然后生成新的UTXO，存放在比特币网络的UTXO池中。
			2. UTXO是不能再分割、被所有者锁住或记录于区块链中的并被整个网络识别成货币单位的一定量的比特币货币。
			3. 比特币网络监测着以百万为单位的所有可用的（未花费的）UTXO。当一个用户接收比特币时，金额被当作UTXO记录到区块链里。这样，一个用户的比特币会被当作UTXO分散到数百个交易和数百个区块中。
			4. 实际上，并不存在储存比特币地址或账户余额的地点，只有被所有者锁住的、分散的UTXO。
			5. "一个用户的比特币余额"，这个概念是一个通过比特币钱包应用创建的派生之物。比特币钱包通过扫描区块链并聚合所有属于该用户的UTXO来计算该用户的余额。
			6. UTXO被每一个全节点比特币客户端在一个储存于内存中的数据库所追踪，该数据库也被称为“UTXO集”或者"UTXO池"。新的交易从UTXO集中消耗（支付）一个或多个输出。

*/

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

//Transaction 交易
type Transaction struct {
	TXID      []byte     //交易ID
	TXInputs  []TXInput  //交易输入(N个)
	TXOutputs []TXOutput //交易输出（N个）
	TimeStamp uint64     //创建交易的时间
}

//TXInput 交易输入：指明交易发起人可支付资金的来源
type TXInput struct {
	TXID       []byte //引用output所在交易的ID
	Index      int64  //引用output在output中的索引值
	ScriptSign string //锁定脚本：付款人对当前新交易的签名
}

//TXOutput 交易输出：包含资金接收方的相关信息
type TXOutput struct {
	Value        float64 //转账金额
	ScriptPubKey string  //锁定脚本：收款人的公钥哈希（地址）
}

//获取交易ID：计算交易哈希
func (tx *Transaction) setHash() error {
	//对tx进行gob编码获得字节流，然后计算sha256，赋值给TXID
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
	return nil
}

//挖矿奖励
var reward = 12.5

//NewCoinbaseTX 创建挖矿交易(没有input因此不需要签名，只有一个output获得挖矿奖励)
func NewCoinbaseTX(miner /*矿工*/ string, data string) *Transaction {
	input := TXInput{TXID: nil, Index: -1, ScriptSign: data} //挖矿不需要签名，由矿工任意填写
	output := TXOutput{Value: reward, ScriptPubKey: miner}
	timStamp := time.Now().Unix()

	tx := Transaction{
		TXID:      nil,
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
		TimeStamp: uint64(timStamp),
	}
	tx.setHash()
	return &tx
}
