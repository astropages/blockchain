package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"sort"
)

//WalletManager 钱包管理：对外管理生成的钱包（公钥,私钥）
//私钥 -> 公钥 -> 地址
type WalletManager struct {
	Wallets map[string]*Wallet //管理所有钱包的map(key为地址,value为钱包)
}

//NewWalletManager 创建WalletManager
func NewWalletManager() *WalletManager {
	//创建一个钱包管理
	var wm WalletManager

	//创建钱包map
	wm.Wallets = make(map[string]*Wallet)

	//从磁盘加载已创建的钱包到map
	if !wm.loadFile() {
		return nil
	}

	//返回钱包map
	return &wm
}

func (wm *WalletManager) createWallet() string {
	//创密钥对
	w := NewWalletKeyPair()
	if w == nil {
		fmt.Println("钱包密钥对创建失败")
		return ""
	}

	//获取地址
	address := w.getAddress()

	//将钱包写入到map，key为钱包地址
	wm.Wallets[address] = w

	//将密钥对写入磁盘
	if !wm.saveFile() {
		return ""
	}

	//返回地址
	return address

}

//钱包文件
const walletFile = "wallet.dat"

//保存WalletManager到磁盘
func (wm *WalletManager) saveFile() bool {
	//使用gob对wm进行编码
	var buffer bytes.Buffer

	//注册椭圆曲线的接口函数后才能进行编码
	gob.Register(elliptic.P256())

	//编码
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(wm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//将WalletManager写入文件
	err = ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//读取钱包文件并加载到WalletManager
func (wm *WalletManager) loadFile() bool {

	//判断文件是否存在
	if !IsFileExist(walletFile) {
		fmt.Println("钱包文件不存在")
		return true
	}
	//读取文件
	content, err := ioutil.ReadFile(walletFile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//创建解码器
	decoder := gob.NewDecoder(bytes.NewReader(content))

	//注册椭圆曲线的接口函数后才能进行解码
	gob.Register(elliptic.P256())

	//解码并赋值到wm
	err = decoder.Decode(wm)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//获取所有钱包地址
func (wm *WalletManager) listAddresses() []string {
	var addresses []string
	for address := range wm.Wallets {
		addresses = append(addresses, address)
	}

	//排序
	sort.Strings(addresses)

	return addresses
}
