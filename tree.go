package main

import "math"

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
func HoffmanTree(str string) Node {
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
