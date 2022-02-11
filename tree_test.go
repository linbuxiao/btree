package btree

import (
	"fmt"
	"testing"
)

func TestTree_Get(t *testing.T) {
	tr := &tree{
		root: &node{
			items: items{3},
			children: children{
				{
					items: items{1, 2},
				},
				{
					items: items{4, 5, 6},
				},
			},
		},
	}

	_, ok := tr.Get(5)
	if !ok {
		t.Fatal("should get 5")
	}
	_, ok = tr.Get(7)
	if ok {
		t.Fatal("should not get 7")
	}
}

func TestOther(t *testing.T) {
	fmt.Println(5 / 2)
}

func insertAt(items []int, i, key int) []int {
	items = append(items, -1)
	copy(items[i+1:], items[i:len(items)-1])
	items[i] = key
	return items
}

func TestInsertAt(t *testing.T) {
	items := []int{1, 3, 4}
	fmt.Println(insertAt(items, 1, 2))
}

func TestTree_Add(t *testing.T) {
	tr := &tree{
		root: &node{
			items: items{3, 6, 9, 12},
			children: []*node{
				{items: items{1, 2}},
				{items: items{4, 5}},
				{items: items{7, 8}},
				{items: items{10, 11}},
				{items: items{13, 14, 15, 16}},
			},
		},
	}

	tr.Add(17)
	fmt.Println(tr)
}

func TestTree_Delete(t *testing.T) {
	tr := &tree{
		root: &node{
			items: items{3, 6},
			children: children{
				{items: items{1, 2}},
				{items: items{4, 5}},
				{items: items{7, 8}},
			},
		},
	}

	tr.Delete(4)
	fmt.Println(tr)
}
