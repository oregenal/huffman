package main

import (
	"fmt"
	"os"
)

const input_string = "hello, world!"

func main() {
	hoffmanTree := HoffmanTree(input_string)
	// PrintTree(hoffmanTree)

	codes := Codes(hoffmanTree)
	// fmt.Println(codes)
	PrintCodes(codes)

	encstr := EncryptedString(input_string, codes)
	// encstr := "01101001"
	fmt.Println(encstr)

	encbin := EncryptedBinary(encstr)
	fmt.Println(encbin)
	for _, b := range(encbin) {
		fmt.Printf("|%08b", b)
	}
	fmt.Println("|")
	for _, b := range(encbin) {
		fmt.Printf("|%x", b)
	}
	fmt.Println("|")

	err := os.WriteFile("test.bin", encbin, 446)
	if err != nil {
		panic("write file error")
	}
}
