package main

import "fmt"

type Code struct {
	code    int
	counter int
}

// May be faster to use not String but Rune
type MapCodes map[string]Code

func PrintCodes(codes MapCodes) {
	for k, v := range codes {
		fmt.Printf("%s: %0*b\n", k, v.counter, v.code)
	}
}

func Codes(tree Node) MapCodes {
	codes := MapCodes{}

	codesRec(&tree, codes, Code{0, 0})

	return codes
}

// Map pass thru by reference by default
func codesRec(node *Node, codes MapCodes, code Code) {
	if node.left != nil {
		code.code <<= 1
		code.counter += 1
		codesRec(node.left, codes, code)

		code.code |= 0b1
		codesRec(node.right, codes, code)
	} else {
		codes[node.c] = code
	}
}
