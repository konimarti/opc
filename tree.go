package opc

import "fmt"

//Tree creates an OPC browser representation
type Tree struct {
	Name     string
	Parent   *Tree
	Branches []*Tree
	Leaves   []Leaf
}

//Leaf contains the OPC tag and forms part of the Tree struct for the  OPC browser
type Leaf struct {
	Name string
	Tag  string
}

//ExtractBranchByName return substree with name
func ExtractBranchByName(tree *Tree, name string) *Tree {
	if tree.Name == name {
		return tree
	}
	for _, b := range tree.Branches {
		subtree := ExtractBranchByName(b, name)
		if subtree != nil {
			return subtree
		}
	}
	return nil
}

//CollectTags traverses tree and collects all tags in string slice
func CollectTags(tree *Tree) []string {
	collection := []string{}
	for _, l := range tree.Leaves {
		collection = append(collection, l.Tag)
	}
	for _, b := range tree.Branches {
		lowerCollection := CollectTags(b)
		collection = append(collection, lowerCollection...)
	}
	return collection
}

//PrettyPrint prints tree in a nice format
func PrettyPrint(tree *Tree) {
	fmt.Println(tree.Name)
	printSubtree(tree, 1)
}

// printSubtree is a recursive helper function to traverse the tree
func printSubtree(tree *Tree, level int) {
	space := ""
	for i := 0; i < level; i++ {
		space += "  "
	}
	for _, l := range tree.Leaves {
		fmt.Println(space, "-", l.Tag)
	}
	for _, b := range tree.Branches {
		fmt.Println(space, "+", b.Name)
		printSubtree(b, level+1)
	}
}
