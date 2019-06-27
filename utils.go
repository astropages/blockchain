package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
