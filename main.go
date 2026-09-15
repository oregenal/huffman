package main

import "fmt"

const input_string = "hello, world!"

func EncryptedString(input string, codes MapCodes) string {
	var result string
	for _, c := range(input) {
		v := codes[string(c)]
		result += fmt.Sprintf("%0*b", v.counter, v.code)
	}
	return result
}

func EncryptedBinary(input string) []byte {
	result := []byte{}
	counter := 0
	var b byte = 0

	for _, c := range(input) {
		if c == rune('1') {
			b |= 1
			counter += 1
		} else if c == rune('0') {
			counter += 1
		} else {
			panic("unsupported symbol")
		}

		if counter < 8 {
			b <<= 1
		} else if counter == 8 {
			result = append(result, b)
			b = 0
			counter = 0
		} else {
			panic("UNREACHABLE")
		}
	}

	if counter != 0 {
		b <<= (7 - counter)
		result = append(result, b)
	}
	return result
}

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
}
