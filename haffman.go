package main

import (
	"fmt"
	"math"
)

const input_string = "A poem is a piece of creative writing written in lines and stanzas that uses the sound, rhythm, and artistic meaning of words to share ideas and feelings.\nIf you would like, tell me what topic or feeling you want to write about, and I can help you compose a short poem!"

// In x86_64 Node struct take 5 registers
// So it can be provided/returned by value
type Node struct {
	c           string
	score       int
	left, right *Node
}

// While range on map in Go work randomly
// the result on each function call
// will be different
func haffmanTree(str string) Node {
	scores := map[string]int{}
	nodes := []Node{}

	for _, c := range str {
		scores[string(c)] += 1
	}

	for k, v := range scores {
		nodes = append(nodes, Node{c: k, score: v})
	}

	smallest := func() Node {
		value := math.MaxInt
		index := 0
		for i, n := range nodes {
			if n.score < value {
				value = n.score
				index = i
			}
		}

		nodes_count := len(nodes)
		result := nodes[index]
		if nodes_count-1 != index {
			nodes[index] = nodes[nodes_count-1]
		}
		nodes = nodes[:nodes_count-1]

		return result
	}

	for len(nodes) > 1 {
		smallest1 := smallest()
		smallest2 := smallest()
		newNode := Node{
			c:     smallest1.c + smallest2.c,
			score: smallest1.score + smallest2.score,
			left:  &smallest1,
			right: &smallest2,
		}
		nodes = append(nodes, newNode)
	}

	return nodes[0]
}

// Looks like maximum Nodes that can be,
// is all printable symbols * 2 + 1.
// So this function safely implemented using recursion.
func PrintTree(tree Node) {
	fmt.Println(printTreeRec(&tree))
}

func printTreeRec(tree *Node) string {
	if tree == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{%s %d 0:%s 1:%s}",
		tree.c, tree.score,
		printTreeRec(tree.left),
		printTreeRec(tree.right))
}

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

type BinCode []byte

func Encrypt(input string, haffmanTree Node) BinCode {
	codes := Codes(haffmanTree)

	return EncryptedBinary(EncryptedString(input, codes))
}

func EncryptedString(input string, codes MapCodes) string {
	var result string

	for _, c := range input {
		v := codes[string(c)]
		result += fmt.Sprintf("%0*b", v.counter, v.code)
	}

	return result
}

func EncryptedBinary(input string) BinCode {
	result := BinCode{}
	counter := 0
	var b byte = 0

	for _, c := range input {
		if c == '1' {
			b |= 1
			counter += 1
		} else if c == '0' {
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

func Decrypt(tree Node, encbin BinCode) string {
	result := ""
	node := tree

	for _, v := range encbin {
		counter := 0

		for counter < 8 {
			if node.left == nil {
				result += string(node.c)
				node = tree
			}

			k := 0b10000000 & v
			v <<= 1
			counter += 1

			if k == 0 {
				node = *node.left
			} else {
				node = *node.right
			}
		}
	}

	return result
}

func main() {
	haffmanTree := haffmanTree(input_string)

	encbin := Encrypt(input_string, haffmanTree)

	result := Decrypt(haffmanTree, encbin)
	fmt.Println(result)
}
