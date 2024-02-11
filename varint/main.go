package main

import (
	"fmt"
	"encoding/hex"
)

func main() {
	fmt.Println("hello world!!!")
	var num uint64 = 824
	byteArray := []byte{}
	for num > 0 {
		currentNum := num & 0x7F
		num = num >> 7
		if num > 0 {
			currentNum |= 0x80
		}
		byteArray = append(byteArray, byte(currentNum))
	}
	hexString := hex.EncodeToString(byteArray)
    fmt.Println("Byte array as hexadecimal string:", hexString)
	fmt.Println(byteArray)
}
