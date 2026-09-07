package main

import "fmt"

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
