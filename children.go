package btree

type children []*node

const (
	directionNotFound = iota
	directionLeft
	directionRight
)

func (c *children) insertAt(i int, prev, next *node) {
	*c = append(*c, nil, nil)
	if i < len(*c) {
		copy((*c)[i+2:], (*c)[i:])
	}
	(*c)[i] = prev
	(*c)[i+1] = next
}

func (c *children) canBorrow(i int) int {
	if i-1 > 0 && len((*c)[i-1].children)-1 > min {
		return directionLeft
	} else if i+1 < len(*c) && len((*c)[i+1].children)-1 > min {
		return directionRight
	}
	return directionNotFound
}
