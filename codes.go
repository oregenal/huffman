package main

func Codes(tree Node) MapCodes {
	codes := MapCodes{}

	codesRec(&tree, codes, Code{0,0})

	return codes
}

// Map pass thru by reference by default
func codesRec(node *Node, codes MapCodes, code Code) {
	if len(node.c) != 1 {
		code.code <<= 1
		code.counter += 1
		codesRec(node.left, codes, code)

		code.code |= 0b1
		codesRec(node.right, codes, code)
	} else {
		codes[node.c] = code
	}
}
