package main

// In x86_64 Node struct take 5 registers
// So it can be provided/returned by value
type Node struct {
	c           string
	score       int
	left, right *Node
}

type Code struct {
	code int;
	counter int;
}

// May be faster to use not String but Rune
type MapCodes map[string]Code

