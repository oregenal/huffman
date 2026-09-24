package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"unicode/utf8"
)

const dataFile = "data.haz"

// BUG disappearing point at the end
const input_string = "A poem is a piece of creative writing written in lines and stanzas that uses the sound, rhythm, and artistic meaning of words to share ideas and feelings."

// const input_string = "A poem is a piece of creative writing written in lines and stanzas that uses the sound, rhythm, and artistic meaning of words to share ideas and feelings.\nIf you would like, tell me what topic or feeling you want to write about, and I can help you compose a short poem!"

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
	Code    int32
	Counter int32
}

// May be faster to use not String but Rune
type MapCodes map[string]Code

func PrintCodes(codes MapCodes) {
	for k, v := range codes {
		fmt.Printf("|%s: %0*b", k, v.Counter, v.Code)
	}
	fmt.Println("|")
}

func Codes(tree Node) MapCodes {
	codes := MapCodes{}

	codesRec(&tree, codes, Code{0, 0})

	return codes
}

// Map pass thru by reference by default
func codesRec(node *Node, codes MapCodes, code Code) {
	if node.left != nil {
		code.Code <<= 1
		code.Counter += 1
		codesRec(node.left, codes, code)

		code.Code |= 0b1
		codesRec(node.right, codes, code)
	} else {
		codes[node.c] = code
	}
}

func TreeFromCodes(codes MapCodes) Node {
	tree := Node{}
	for k, v := range codes {
		strcode := fmt.Sprintf("%0*b", v.Counter, v.Code)
		node := &tree

		for _, c := range strcode {
			if c == '1' {
				if node.right == nil {
					// This Node goes in heap??
					node.right = &Node{}
				}
				node = node.right
			} else if c == '0' {
				if node.left == nil {
					node.left = &Node{}
				}
				node = node.left
			} else {
				log.Fatalf("UREACHABLE")
			}
		}

		node.c = k
	}
	return tree
}

func DataToBin(encbin BinCode) []byte {
	data := make([]byte, 8)

	encoded, err := binary.Encode(data, binary.LittleEndian, encbin.size)
	if err != nil {
		log.Fatalf("encode fail %v", err)
	}
	if encoded != 8 {
		panic("encoding wrong size")
	}

	data = append(data, encbin.data...)
	return data
}

func CodesToBin(codes MapCodes) []byte {
	result := make([]byte, 4)
	size := len(codes)

	encoded, err := binary.Encode(result, binary.LittleEndian, int32(size))
	if err != nil {
		log.Fatalf("encode fail %v", err)
	}
	if encoded != 4 {
		panic("encoding wrong size")
	}

	for k, v := range codes {
		runeBuf := make([]byte, 4)

		r, utfSize := utf8.DecodeRuneInString(k)
		if utfSize > 4 {
			panic("wrong symbol")
		}

		encoded, err := binary.Encode(runeBuf, binary.LittleEndian, r)
		if err != nil {
			log.Fatalf("encode fail %v", err)
		}
		if encoded != 4 {
			panic("encoding wrong size")
		}
		result = append(result, runeBuf...)

		scoreBuf := make([]byte, 8)
		encoded, err = binary.Encode(scoreBuf, binary.LittleEndian, v)
		if err != nil {
			log.Fatalf("encode fail %v", err)
		}
		if encoded != 8 {
			panic("encoding wrong size")
		}
		result = append(result, scoreBuf...)
	}
	return result
}

func BinToCodes(buf []byte, size int32) MapCodes {
	result := MapCodes{}
	for i := range size {
		var r rune
		sliseNum := i * 12
		decoded, err := binary.Decode(
			buf[sliseNum:],
			binary.LittleEndian, &r)
		if err != nil {
			log.Fatalf("decode fail %v", err)
		}
		if decoded != 4 {
			panic("encoding wrong size")
		}

		payload := Code{}
		decoded, err = binary.Decode(
			buf[sliseNum+4:],
			binary.LittleEndian, &payload)
		if err != nil {
			log.Fatalf("decode fail %v", err)
		}
		if decoded != 8 {
			panic("encoding wrong size")
		}
		result[string(r)] = payload
	}

	return result
}

func BinToData(binDataBuffer []byte) BinCode {
	result := BinCode{
		data: binDataBuffer[8:],
	}
	_, err := binary.Decode(
		binDataBuffer[:8],
		binary.LittleEndian,
		&result.size)
	if err != nil {
		log.Fatalf("decode fail %v", err)
	}

	return result
}

type BinCode struct {
	size int64
	data []byte
}

func Encrypt(input string, haffmanTree Node) BinCode {
	codes := Codes(haffmanTree)

	return EncryptedBinary(EncryptedString(input, codes))
}

func EncryptedString(input string, codes MapCodes) string {
	var result string

	for _, c := range input {
		v := codes[string(c)]
		result += fmt.Sprintf("%0*b", v.Counter, v.Code)
	}

	return result
}

func EncryptedBinary(input string) BinCode {
	result := BinCode{}
	result.size = 0
	var b byte = 0

	for _, c := range input {
		if c == '1' {
			b |= 1
			result.size += 1
		} else if c == '0' {
			result.size += 1
		} else {
			panic("unsupported symbol")
		}

		if result.size%8 != 0 {
			b <<= 1
		} else if result.size%8 == 0 {
			result.data = append(result.data, b)
			b = 0
		} else {
			panic("UNREACHABLE")
		}
	}

	if result.size%8 != 0 {
		b <<= (7 - result.size%8)
		result.data = append(result.data, b)
	}

	return result
}

func Decrypt(tree Node, encbin BinCode) string {
	result := ""
	node := tree
	counter := encbin.size

	for _, v := range encbin.data {
		bit := 0

		for bit < 8 {
			if node.left == nil {
				result += string(node.c)
				node = tree
			}

			if counter == 0 {
				break
			}

			k := 0b10000000 & v
			v <<= 1
			bit += 1
			counter -= 1

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
	codes := Codes(haffmanTree)
	encbin := EncryptedBinary(EncryptedString(input_string, codes))

	// Convert to binary data
	// so we can save it to file
	bincodes := CodesToBin(codes)
	bindata := DataToBin(encbin)
	bin := []byte{}
	bin = append(bin, bincodes...)
	bin = append(bin, bindata...)

	// File manipulations
	os.WriteFile(dataFile, bin, 0644)

	fromFile, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("read file fail %v", err)
	}

	// Decode BinCode from binary format
	var binCodesSize int32
	binary.Decode(fromFile[:4], binary.LittleEndian, &binCodesSize)
	binCodeBufferEnd := 12*binCodesSize + 4
	binCodesBuffer := fromFile[4:binCodeBufferEnd]
	binDataBuffer := fromFile[binCodeBufferEnd:]

	newCodes := BinToCodes(binCodesBuffer, binCodesSize)
	newData := BinToData(binDataBuffer)
	newTree := TreeFromCodes(newCodes)

	result := Decrypt(newTree, newData)
	fmt.Println(result)
}
