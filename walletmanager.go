package main

import "fmt"

//WalletManager 钱包管理：对外管理生成的钱包（公钥,私钥）
//私钥 -> 公钥 -> 地址
type WalletManager struct {
	Wallets map[string]*Wallet //管理所有钱包的集合(key为地址,value为钱包)
}

//NewWalletManager 创建WalletManager
func NewWalletManager() *WalletManager {
	//创建一个钱包管理
	var wm WalletManager

	//从磁盘加载已创建的钱包到集合
	//TODO

	//返回钱包集合
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
	adderss := w.getAddress()

	//将密钥对写入磁盘
	//TDDO

	//返回地址
	return adderss

}
