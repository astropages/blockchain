package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

//Wallet 钱包
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey //私钥
	//X,Y类型一致，将X和Y拼接成字节流赋值给publicKey字段用于传输
	//验证时将X和Y截取出来再创建一条曲线，还原公钥以进行校验
	PublicKey []byte //公钥
}

//NewWalletKeyPair 创建钱包：密钥对
func NewWalletKeyPair() *Wallet {
	//创建私钥
	curve := elliptic.P256()                                 //创建曲线
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader) //生成私钥
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//通过私钥获得公钥
	publicKey := privateKey.PublicKey

	//将公钥的X,Y进行拼接
	pubKey := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	//返回
	wallet := Wallet{privateKey, pubKey}
	return &wallet
}

//根据私钥生成地址
func (w *Wallet) getAddress() string {
	//公钥
	pubKey := w.PublicKey
	hash1 := sha256.Sum256(pubKey)

	//哈希160处理
	hasher := ripemd160.New()
	hasher.Write(hash1[:])
	//计算出公钥哈希（锁定output）
	pubKeyHash := hasher.Sum(nil)

	//拼接version和公钥哈希，得到21字节的数据
	payload := append([]byte{byte(0x00)}, pubKeyHash...)

	//生成4个字节的校验码
	frist := sha256.Sum256(payload)
	second := sha256.Sum256(frist[:])
	//4字节校验码
	checksum := second[0:4]
	//25字节数据
	payload = append(payload, checksum...)
	//地址
	address := base58.Encode(payload)

	return address
}
