package main

import (
	"fmt"
)

const input_sting = "hello, world!"

type Node struct {
	c string
	score int;
	left, right *Node;
}

func Get_Tree(str string) Node {
	scores := map[string]int{}
	leefs := []Node{}
	// fmt.Println(input_sting)

	for _, c := range str {
		scores[string(c)] += 1
	}
	// fmt.Println(scores)

	for k, v := range scores {
		leefs = append(leefs, Node{c: k, score: v})
	}
	// fmt.Println(leefs)

	// I can just make func in func with closure
	Smallest := func () Node {
		value := 99999999999
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
		smallest1 := Smallest()
		smallest2 := Smallest()
		str := smallest1.c + smallest2.c
		val := smallest1.score + smallest2.score
		result := Node{c: str, score: val, left: &smallest1, right: &smallest2}
		leefs = append(leefs, result)
		// fmt.Println(result)
	}

	return leefs[0]
}

func main() {
	result := Get_Tree(input_sting)
	fmt.Println(result)
}
