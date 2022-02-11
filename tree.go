package btree

const (
	M   = 5
	min = 2 // math.ceil(5 / 3 - 1)
)

type tree struct {
	root *node
}

func (tr *tree) Get(item int) (int, bool) {
	return tr.root.get(item)
}

func (tr *tree) Add(item int) *tree {
	newRootItem, prev, next := tr.root.add(item)
	if newRootItem == -1 {
		tr.root = next
	} else {
		tr.root = &node{
			items:    items{newRootItem},
			children: []*node{prev, next},
		}
	}
	return tr
}

func (tr *tree) Delete(item int) *tree {
	n := tr.root.delete(item)
	tr.root = n
	return tr
}
