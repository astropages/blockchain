package main

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

//BlockChain 区块链
type BlockChain struct {
	// Blocks []*Block
	db   *bolt.DB //用于存储数据的数据库
	tail []byte   //最后一个区块的哈希值
}

//创世语
const genesisInfo = "I am alpha."

//数据库名
const blockChainDBFile = "blockchain.db"

//数据桶
const blockBucket = "blockBucket"

//数据桶中保存最后一个区块哈希值的字段key
const lastBlockHashKey = "lastBlockHashKey"

//CreateBlockChain 创建区块链（同时添加创世块）
func CreateBlockChain() error {

	//打开数据库，没有则创建
	db, err := bolt.Open(blockChainDBFile, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	//开始创建
	err = db.Update(func(tx *bolt.Tx) error {
		//打开数据桶
		bucket := tx.Bucket([]byte(blockBucket))
		//如果数据桶不存在则创建
		if bucket == nil {
			//创建数据桶
			bucket, err := tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				return err
			}
			//创建挖矿交易
			coinbase := NewCoinbaseTX("中本聪", genesisInfo)
			//拼装交易集合txs
			txs := []*Transaction{coinbase}
			//新建创世快
			genesisBlock := NewBlock(txs, nil)
			//将区块数据流写入数据库（key为区块的哈希，value为区块的数据流）
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			//将最后一个区块的哈希写入数据库（key为lastBlockHash,value为创世块的哈希）
			bucket.Put([]byte(lastBlockHashKey), genesisBlock.Hash)
			fmt.Println("创建区块链成功")
		} else {
			fmt.Println("区块链已存在")
		}

		return nil
	})
	return err
}

//GetBlockChainInstance 获取区块链实例
func GetBlockChainInstance() (*BlockChain, error) {
	//内存中的最后一个区块的哈希值
	var lastHash []byte

	//打开数据库
	db, err := bolt.Open(blockChainDBFile, 0400, nil) //只有读权限
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//不关闭数据库

	//查询数据库事务
	db.View(func(tx *bolt.Tx) error {
		//打开数据桶
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return errors.New("No bucket")
		}
		//从数据桶获取最后一个区块的哈希值
		lastHash = bucket.Get([]byte(lastBlockHashKey))
		return nil
	})

	//返回区块链实例
	bc := BlockChain{db, lastHash}
	return &bc, nil
}

//AddBlock 向区块链中添加区块的方法（传入数据：交易集合）
func (bc *BlockChain) AddBlock(txs []*Transaction) error {

	//获取最后一个区块的哈希
	lastBlockHash := bc.tail

	//创建一个新区块
	newBlock := NewBlock(txs, lastBlockHash)

	//写入数据库
	err := bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return errors.New("No bucket")
		}
		//写入新区块到数据库（key为区块的哈希，value为区块的数据字节流）
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			return err
		}
		//更新lastBlockHashKey（数据库中记录最后一个区块哈希的值）
		err = bucket.Put([]byte(lastBlockHashKey), newBlock.Hash)
		if err != nil {
			return err
		}
		//更新区块链的tali值（最后一个区块的哈希值）
		bc.tail = newBlock.Hash
		fmt.Println("添加区块成功")
		return nil
	})
	return err
}

//Iterator 迭代器（用于实现区块遍历）
type Iterator struct {
	db          *bolt.DB
	currentHash []byte //游标：不断移动的哈希值
}

//NewIterator 初始化迭代器的方法
func (bc *BlockChain) NewIterator() *Iterator {
	it := Iterator{
		db:          bc.db,
		currentHash: bc.tail, //最后一个区块的哈希值
	}
	return &it
}

//Next 迭代器Next方法，返回当前指向的区块并向左移动游标指向前一个区块
func (it *Iterator) Next() (block *Block) {
	//从数据库读取当前哈希
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			return errors.New("No bucket")
		}
		//获取到最后一个区块的字节流
		tmpBlockInfo := bucket.Get([]byte(it.currentHash))
		//获取最后一个区块结构
		block = DeSerialize(tmpBlockInfo)
		//游标前移：从区块结构获取前一个区块的哈希值并赋值给游标
		it.currentHash = block.PrevHash
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return
}

//UTXOInfo UTXO详情
type UTXOInfo struct {
	TXID     []byte //交易ID
	Index    int64  //索引值
	TXOutput        //继承自output
}

//FindMyUTXO 获取指定地址的金额：遍历账本
func (bc *BlockChain) FindMyUTXO(address string) []UTXOInfo {
	var utxoInfos []UTXOInfo                //UTXO集合
	var spentUtxos = make(map[string][]int) //定义一个存放已消耗交易输出集合的集合

	it := bc.NewIterator() //定义迭代器

	for {
		//遍历区块
		block := it.Next()
		//遍历交易
		for _, tx := range block.Transactions {
		LABEL:
			//遍历outputs，判断其锁定脚本是否为目标地址
			for outputIndex, output := range tx.TXOutputs {
				if output.ScriptPubKey == address {
					//过滤
					currentTXID := string(tx.TXID)
					//在集合中查找集合
					indexArray := spentUtxos[currentTXID]
					//判断该交易ID是否有数据，有则代表已被某个output使用
					if len(indexArray) != 0 {
						for _, spendIndex := range indexArray {
							//判断下标
							if outputIndex == spendIndex {
								continue LABEL
							}
						}

					}
					//找到属于目标地址的utxo详情
					utxoInfo := UTXOInfo{tx.TXID, int64(outputIndex), output}
					utxoInfos = append(utxoInfos, utxoInfo)
				}
			}

			//遍历inputs
			for _, input := range tx.TXInputs {
				if input.ScriptSign == address {
					//key交易ID，value为交易输出索引的集合
					spentKey := string(input.TXID)
					//向集合中添加已消耗交易输出的集合
					spentUtxos[spentKey] = append(spentUtxos[spentKey], int(input.Index))
				}
			}
		}
		//退出条件
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return utxoInfos

}

//遍历账本（转账人地址，转账金额）找到from能使用的utxo集合及包含的所有金额
func (bc *BlockChain) findNeedUTXO(from string, amount float64) (map[string][]int64, float64) {
	var retMap = make(map[string][]int64)
	var retValue float64

	//遍历账本，找到所有utxo集合
	utxoInfos := bc.FindMyUTXO(from)
	//遍历utxo,统计总金额
	for _, utxoInfo := range utxoInfos {
		retValue += utxoInfo.Value                        //utxo总额
		key := string(utxoInfo.TXID)                      //
		retMap[key] = append(retMap[key], utxoInfo.Index) //将要使用的utxo集合
		//如果总金额大于转账金额，直接返回
		if retValue >= amount {
			break
		}
		//否则继续遍历
	}

	return retMap, retValue
}

/*





	//Bolt数据库
		//linux下可通过strings命令查看二进制数据库文件的可读内容
		//例：$ strings test.db





	////定义数据库名
	//const testDB = "test.db"

	////打开数据库，没有则创建
	//db, err := bolt.Open(testDB, 0600, nil)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer db.Close()

	////创建bucket
	//err = db.Update(func(tx *bolt.Tx) error {
	//	//打开一个bucket
	// 	b1 := tx.Bucket([]byte("bucket1"))
	// 	//判断bucket是否存在
	// 	if b1 == nil {
	// 		//没有则创建
	// 		b1, err = tx.CreateBucket([]byte("bucket1"))
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return err
	// 		}
	// 	}
	// 	//写入数据
	// 	b1.Put([]byte("key1"), []byte("value1"))
	// 	//读取数据
	// 	v1 := b1.Get([]byte("key1"))
	// 	//打印数据
	// 	fmt.Printf("key1: %s\n", v1)
	// 	return nil
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

*/
