package main

import (
	"fmt"
	"math"
)

const input_sting = "hello, world!"

type Node struct {
	c string
	score int;
	left, right *Node;
}

func HoffmanTree(str string) Node {
	scores := map[string]int{}
	leefs := []Node{}

	for _, c := range str {
		scores[string(c)] += 1
	}

	// ramge on map is random
	// so the result on each programm run
	// wuld be different
	for k, v := range scores {
		leefs = append(leefs, Node{c: k, score: v})
	}

	smallest := func () Node {
		value := math.MaxInt
		index := 0
		for i, leef := range leefs {
			if leef.score < value {
				value = leef.score
				index = i
			}
		}

		leefs_count := len(leefs)
		result := leefs[index]
		if leefs_count-1 != index {
			leefs[index] = leefs[leefs_count-1]
		}
		leefs = leefs[:leefs_count-1]

		return result
	}

	for len(leefs) > 1 {
		smallest1 := smallest()
		smallest2 := smallest()
		str := smallest1.c + smallest2.c
		val := smallest1.score + smallest2.score
		result := Node{c: str, score: val, left: &smallest1, right: &smallest2}
		leefs = append(leefs, result)
	}

	return leefs[0]
}

func main() {
	hoffmanTree := HoffmanTree(input_sting)
	fmt.Println(hoffmanTree)
}
