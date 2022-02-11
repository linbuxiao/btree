package btree

import "sort"

type items []int

func (items items) find(key int) (int, bool) {
	i := sort.Search(len(items), func(i int) bool {
		return items[i] >= key
	})
	if i < len(items) && items[i] == key {
		return i, true
	}

	return i, false
}

func (items *items) insertAt(i int, key int) {
	// 神奇的引用！
	*items = append(*items, -1)
	if i < len(*items) {
		copy((*items)[i+1:], (*items)[i:])
	}
	(*items)[i] = key
}

func (items *items) removeAt(i int) {
	*items = append((*items)[:i], (*items)[i+1:]...)
}
