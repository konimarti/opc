package opc

import (
	"testing"
)

func testingCreateNewTree() *Tree {

	root := Tree{
		Name:     "root",
		Parent:   nil,
		Branches: []*Tree{},
		Leaves: []Leaf{
			Leaf{
				Name: "bandwith",
				Tag:  "bandwith",
			},
		},
	}
	options := Tree{
		Name:     "options",
		Parent:   &root,
		Branches: []*Tree{},
		Leaves: []Leaf{
			Leaf{
				Name: "frequency",
				Tag:  "options.frequency",
			},
			Leaf{
				Name: "amplitute",
				Tag:  "options.amplitude",
			},
		},
	}
	numeric := Tree{
		Name:     "numeric",
		Parent:   &root,
		Branches: []*Tree{},
		Leaves: []Leaf{
			Leaf{
				Name: "sin",
				Tag:  "numeric.sin",
			},
			Leaf{
				Name: "cos",
				Tag:  "numeric.cos",
			},
			Leaf{
				Name: "tan",
				Tag:  "numeric.tan",
			},
		},
	}

	root.Branches = append(root.Branches, &options, &numeric)
	return &root
}

func TestTreeExtractBranchByName(t *testing.T) {
	tree := testingCreateNewTree()
	subtree := ExtractBranchByName(tree, "numeric")
	if subtree == nil {
		t.Fatal("subtree not correctly extracted")
	}
	if len(CollectTags(subtree)) != 3 {
		t.Fatal("subtree not correctly extracted")
	}
}

func TestTreeCollectTags(t *testing.T) {
	tree := testingCreateNewTree()
	collection := CollectTags(tree)
	if len(collection) != 6 {
		t.Fatal("not enough tags collected")
	}
}
