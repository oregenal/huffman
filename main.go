package main

import (
	"fmt"
)

const input_string = "A poem is a piece of creative writing written in lines and stanzas that uses the sound, rhythm, and artistic meaning of words to share ideas and feelings.\nIf you would like, tell me what topic or feeling you want to write about, and I can help you compose a short poem!"

func Decrypt(tree Node, encbin []byte) string {
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
	hoffmanTree := HoffmanTree(input_string)

	encbin := Encrypt(input_string, hoffmanTree)

	result := Decrypt(hoffmanTree, encbin)
	fmt.Println(result)
}
