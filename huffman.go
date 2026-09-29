// Toy file compressor realized using Huffman algorithm
//
// Compressed file structure:
// 		[header(".HUZ")][tableSize(int32)][Table][Data]

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"unicode/utf8"
)

const (
	dataFile   = "data.huz"
	header     = ".HUZ"
	headerSize = int32(len(header))
	tableSize  = 4  // Int32
	codeSize   = 12 // Rune + 2 * int32
)

// BUG disappearing point at the end
const input_string = "A poem is a piece of creative writing written in lines and stanzas that uses the sound, rhythm, and artistic meaning of words to share ideas and feelings."

// In x86_64 Node struct take 5 registers
// So it can be provided/returned by value
type Node struct {
	C     string `json:"c,omitempty"`
	Score int    `json:"s,omitzero"`
	Lo    *Node  `json:"lo,omitempty"`
	Hi    *Node  `json:"hi,omitempty"`
}

// May be faster to use not String but Rune
type MapCodes map[rune]Code

type Code struct {
	Code    int32
	Counter int32
}

type BinCode struct {
	size int64
	data []byte
}

// While range on map in Go work randomly
// the result on each function call
// will be different
func huffmanTree(str string) Node {
	scores := map[rune]int{}
	nodes := []Node{}

	for _, c := range str {
		scores[c] += 1
	}

	for k, v := range scores {
		nodes = append(nodes, Node{C: string(k), Score: v})
	}

	smallest := func() Node {
		value := math.MaxInt
		index := 0
		for i, n := range nodes {
			if n.Score < value {
				value = n.Score
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
			C:     smallest1.C + smallest2.C,
			Score: smallest1.Score + smallest2.Score,
			Lo:    &smallest1,
			Hi:    &smallest2,
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
		tree.C, tree.Score,
		printTreeRec(tree.Lo),
		printTreeRec(tree.Hi))
}

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
	if node.Lo != nil {
		code.Code <<= 1
		code.Counter += 1
		codesRec(node.Lo, codes, code)

		code.Code |= 0b1
		codesRec(node.Hi, codes, code)
	} else {
		r, utfSize := utf8.DecodeRuneInString(node.C)
		if utfSize > 4 {
			panic("wrong symbol")
		}

		codes[r] = code
	}
}

func TreeFromCodes(codes MapCodes) Node {
	tree := Node{}
	for k, v := range codes {
		strcode := fmt.Sprintf("%0*b", v.Counter, v.Code)
		node := &tree

		for _, c := range strcode {
			if c == '1' {
				if node.Hi == nil {
					// This Node goes in heap??
					node.Hi = &Node{}
				}
				node = node.Hi
			} else if c == '0' {
				if node.Lo == nil {
					node.Lo = &Node{}
				}
				node = node.Lo
			} else {
				log.Fatalf("UREACHABLE")
			}
		}

		node.C = string(k)
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
	result := make([]byte, tableSize)
	size := int32(len(codes))

	encoded, err := binary.Encode(result, binary.LittleEndian, size)
	if err != nil {
		log.Fatalf("encode fail %v", err)
	}
	if encoded != tableSize {
		panic("encoding wrong size")
	}

	for r, v := range codes {
		runeBuf := make([]byte, 4)

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
		slice := i * codeSize
		decoded, err := binary.Decode(
			buf[slice:],
			binary.LittleEndian, &r)
		if err != nil {
			log.Fatalf("decode fail %v", err)
		}
		if decoded != 4 {
			panic("encoding wrong size")
		}

		payload := Code{}
		decoded, err = binary.Decode(
			buf[slice+4:],
			binary.LittleEndian, &payload)
		if err != nil {
			log.Fatalf("decode fail %v", err)
		}
		if decoded != 8 {
			panic("encoding wrong size")
		}
		result[r] = payload
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

func Encrypt(input string, codes MapCodes) BinCode {
	var newByte byte = 0
	var index int32 = 0
	result := BinCode{}

	for _, r := range input {
		data := codes[r]

		// 00000000        111             00111000
		//   ^         counter = 3  numToShift = 8 - 2 - 3 = 3
		// index = 2                                /     \
		//                                     index      counter
		numToShift := 8 - index - data.Counter
		if numToShift > 0 {
			newByte |= byte(data.Code) << numToShift
			result.size += int64(data.Counter)
			index += data.Counter
		} else if numToShift == 0 {
			newByte |= byte(data.Code) << numToShift
			result.data = append(result.data, newByte)
			index = 0
			result.size += int64(data.Counter)
			newByte = 0
		} else if numToShift < 0 {
			// 00000000       11111             00000111|11000000
			//      ^      counter = 5  numToShift = 8 - 5 - 5 = -2
			// index = 5                                /     \
			//                                     index      counter
			numToShift = -numToShift
			newByte |= byte(data.Code) >> numToShift
			result.data = append(result.data, newByte)
			result.size += int64(data.Counter)
			index = numToShift
			newByte = 0 | byte(data.Code)<<(8-numToShift)
		} else {
			panic("UNREACHABLE")
		}
	}

	return result
}

// Very slow
func EncryptedString(input string, codes MapCodes) string {
	var result string

	for _, c := range input {
		v := codes[c]
		result += fmt.Sprintf("%0*b", v.Counter, v.Code)
	}

	return result
}

func EncryptedStringToBinary(input string) BinCode {
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

// Vary slow
func Decrypt(tree Node, encbin BinCode) string {
	result := ""
	node := tree
	counter := encbin.size

	for _, v := range encbin.data {
		bit := 0

		for bit < 8 {
			if node.Lo == nil {
				result += string(node.C)
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
				node = *node.Lo
			} else {
				node = *node.Hi
			}
		}
	}

	return result
}

func main() {
	huffmanTree := huffmanTree(input_string)
	codes := Codes(huffmanTree)
	// encbin := EncryptedStringToBinary(EncryptedString(input_string, codes))
	encbin := Encrypt(input_string, codes)

	// Convert to binary data
	// so we can save it to file
	bincodes := CodesToBin(codes)
	bindata := DataToBin(encbin)
	bin := []byte{}
	bin = append(bin, []byte(header)...)
	bin = append(bin, bincodes...)
	bin = append(bin, bindata...)

	// File manipulations
	os.WriteFile(dataFile, bin, 0644)

	fromFile, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatalf("read file fail %v", err)
	}

	// Decode BinCode from binary format
	if string(fromFile[:headerSize]) != header {
		panic("vrong file type")
	}

	var binCodesSize int32
	binary.Decode(
		fromFile[headerSize:headerSize+tableSize],
		binary.LittleEndian,
		&binCodesSize)
	binCodeBufferEnd := codeSize*binCodesSize + headerSize + tableSize
	binCodesBuffer := fromFile[headerSize+tableSize : binCodeBufferEnd]
	binDataBuffer := fromFile[binCodeBufferEnd:]

	newCodes := BinToCodes(binCodesBuffer, binCodesSize)
	newData := BinToData(binDataBuffer)
	newTree := TreeFromCodes(newCodes)

	result := Decrypt(newTree, newData)
	fmt.Println(result)
}
