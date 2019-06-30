package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

//UintToByteSlice Uint64转换为[]byte
func UintToByteSlice(num uint64) []byte {
	//创建一个字节缓冲区
	var buffer bytes.Buffer
	//使用二进制编码(以小端对齐方式将数值以二进制编码方式写入到字节缓冲区)
	err := binary.Write(&buffer, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//返回字节切片
	return buffer.Bytes()
}

//IsFileExist 判断文件是否存在
func IsFileExist(filename string) bool {
	//获取文件状态
	_, err := os.Stat(filename)
	//判断文件是否存在
	if os.IsNotExist(err) {
		return false
	}
	return true
}
