package main

import (
	"bytes"
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

	//获得公钥哈希
	pubKeyHash := GetPubKeyHashFromPublicKey(w.PublicKey)

	//拼接version和公钥哈希，得到21字节的数据
	payload := append([]byte{byte(0x00)}, pubKeyHash...)

	//生成4个字节的校验码
	checksum := CheckSum(payload)
	//25字节数据
	payload = append(payload, checksum...)
	//地址
	address := base58.Encode(payload)

	return address
}

//GetPubKeyHashFromPublicKey 通过公钥计算公钥哈希
func GetPubKeyHashFromPublicKey(publickey []byte) []byte {

	hash := sha256.Sum256(publickey)

	//哈希160处理
	hasher := ripemd160.New()
	hasher.Write(hash[:])
	//计算出公钥哈希（锁定output）
	pubKeyHash := hasher.Sum(nil)

	return pubKeyHash
}

//GetPubKeyHashFromAddress 通过地址获取公钥哈希
func GetPubKeyHashFromAddress(address string) []byte {

	//base58解码
	deInfo := base58.Decode(address)

	if len(deInfo) != 25 {
		fmt.Println("地址无效")
		return nil
	}

	//截取
	pubKeyHash := deInfo[1 : len(deInfo)-4]

	return pubKeyHash
}

//CheckSum 获取4字节的校验码
func CheckSum(payload []byte) []byte {
	frist := sha256.Sum256(payload)
	second := sha256.Sum256(frist[:])
	//4字节校验码
	checksum := second[0:4]
	return checksum
}

//IsValidAddress 地址校验：判断地址是否有效
func IsValidAddress(address string) bool {
	//解码，得到25字节数据
	deInfo := base58.Decode(address)
	if len(deInfo) != 25 {
		fmt.Println("地址校验失败")
		return false
	}
	//截取前21字节的payload
	payload := deInfo[:len(deInfo)-4]
	//截取后4字节的checksum1
	checksum1 := deInfo[len(deInfo)-4:]

	//计算payload, 获得checksum2
	checksum2 := CheckSum(payload)
	//对比checksum1和checksum2
	return bytes.Equal(checksum1, checksum2)
}
