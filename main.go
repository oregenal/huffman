package main

import "fmt"

const input_string = "hello, world!"

func main() {
	hoffmanTree := HoffmanTree(input_string)
	// PrintTree(hoffmanTree)
	codes := Codes(hoffmanTree)
	fmt.Println(codes)
	PrintCodes(codes)
}
