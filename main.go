package main

import "fmt"

const input_sting = "hello, world!"

func main() {
	hoffmanTree := HoffmanTree(input_sting)
	// PrintTree(hoffmanTree)
	codes := Codes(hoffmanTree)
	fmt.Println(codes)
	PrintCodes(codes)
}
