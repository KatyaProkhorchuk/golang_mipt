//go:build !solution

package treeiter

func DoInOrder[Tree interface {
	Right() *Tree
	Left() *Tree
}](root *Tree, f func(node *Tree)) {
	if root == nil {
		return
	}
	DoInOrder((*root).Left(), f)
	f(root)
	DoInOrder((*root).Right(), f)
}
