package btree

// b 树是向上生长的，而非向下生长

// 由于 b 树是向上生长，所以底部分裂时可能会设计到上层的分裂，也就是要逐层从下到上递归

// 对于这种操作，avl 树的处理方式时，返回新的节点，重新赋值，直至返回到最上层

// google/btree 使用了另一种方法，它在插入时不理会上层的分裂，在第二次插入时，如果遇到了上次遗留需要分裂的，则向上抛

// 这决定了 google/btree 并不能实时保证它是标准的 btree

// 在它插入数据时，它被破坏了。第二次的插入可能会修复，但也有可能造成遗留更上层的分裂。

// 但越上层是越难破坏的

// Do a single pass down the tree, but before entering (visiting) a node, restructure the tree so that once the key to be deleted is encountered, it can be deleted without triggering the need for any further restructuring

//在树上做一次传递，但在进入（访问）一个节点之前，重组树，以便一旦遇到要删除的键，就可以删除它，而不需要触发任何进一步的重组

// 这是一种高效的方式，减少了重组的次数

// 但我们先不使用，我们仍然使用之前的方法，在递归到 children 时， 接收 递归的返回值作为新的 child

type node struct {
	items    items
	children children
}

func (n *node) get(key int) (int, bool) {
	i, found := n.items.find(key)
	if found {
		return n.items[i], true
	} else if len(n.children) > 0 {
		// 向下查找
		return n.children[i].get(key)
	}
	return 0, false
}

func (n *node) add(key int) (int, *node, *node) {
	i, found := n.items.find(key)
	if found {
		return -1, nil, n
	}

	// 这里存在一个误区：只要有children，必然是要去children里找的。而不是当前node格子空，就可以放到这个格子里去。b树的插入永远是在叶子节点进行。但插入的内容不会永远保持在叶子节点上。
	if len(n.children) == 0 {
		// 叶子节点
		// 没有 children 直接插入
		n.items.insertAt(i, key)
	} else {
		// 找下游
		item, prev, next := n.children[i].add(key)
		// 接到下游的处理结果
		// 检查下游是否分裂
		if item == -1 {
			// 未分裂
			// 直接 return
			return -1, nil, next
		} else {
			// 下游分裂了
			// 删掉原本的 children
			n.children = append(n.children[:i], n.children[i+1:]...)
			// 与当前节点汇合
			n.items.insertAt(i, item)
			// 子节点也要插入
			n.children.insertAt(i, prev, next)
		}
	}
	if itemsIsFull(n.items) {
		return n.split()
	}
	return -1, nil, n
}

func (n *node) delete(key int) *node {
	// 这一步我们仍然使用层序递归的方式，自下而上依次改变结构
	i, found := n.items.find(key)
	if found {
		n.items.removeAt(i)
		return n
	} else {
		if len(n.children) > 0 {
			child := n.children[i].delete(key)
			// 我们只需对子节点进行判断，判断它是否合法，不合法时再进行处理
			if len(child.items) < min {
				// 过少
				// 1. child 为叶子节点，可以找左或右借
				if len(child.children) == 0 {
					// 从兄弟节点借
					switch n.children.canBorrow(i) {
					case directionNotFound:
						// 两边都借不到
						// 合并当前child 与 左侧的child
						n.mergeChild(i)
					case directionLeft:
						leftItems := n.children[i-1].items
						newRootItem := leftItems[len(leftItems)-1] // 拿到最右侧的元素
						oldRootItem := n.items[i]                  // 保存一下
						n.children[i].items = append(items{oldRootItem}, n.children[i].items...)
						n.items[i] = newRootItem // 替换 root
						// 删除左子最右侧元素
						n.children[i-1].items.removeAt(len(n.children[i-1].items) - 1)
					case directionRight:
						rightItems := n.children[i+1].items
						newRootItem := rightItems[0] // 拿到最左侧的元素
						oldRootItem := n.items[i]    // 保存一下
						n.children[i].items = append(n.children[i].items, oldRootItem)
						n.items[i] = newRootItem // 替换 root
						// 删除右子最左侧元素
						n.children[i+1].items.removeAt(0)
					}
				} else {
					// 非叶子节点， 先尝试左旋右旋进行
					// 先尝试右旋
					if i > 0 {
						if len(child.children[i-1].items)-1 >= min {
							child.items = append(child.items, -1)
							child.rightRotation(len(child.items))
							return n
						} else if len(child.children[i].items)-1 >= min {
							child.items = append(child.items, -1)
							child.leftRotation(len(child.items))
							return n
						}
					} else {
						// 无法右旋
						// 尝试左旋
						if len(child.children[i].items)-1 >= min {
							child.items = append(child.items, -1)
							child.leftRotation(len(child.items))
							return n
						} else {
							// 无法左旋 也 无法右旋， 此时需要合并
							n.mergeChild(i)
							return n
						}
					}
				}
			}
		}
	}
	return n
}

func (n *node) split() (item int, prev, next *node) {
	m := getMiddle(n.items)
	item = n.items[m]
	prev = &node{
		items: n.items[:m],
	}
	next = &node{
		items: n.items[m+1:],
	}

	if len(n.children) > 0 {
		prev.children = n.children[:m+1]
		next.children = n.children[m+1:]
	}

	return item, prev, next
}

func (n *node) mergeChild(i int) {
	if i > 0 {
		// 统一向左 merge
		root := n.items[i-1]
		prev := n.children[i-1]
		next := n.children[i]
		newNode := &node{
			items:    append(prev.items, append(items{root}, next.items...)...),
			children: append(prev.children, prev.children...),
		}
		n.children[i-1] = newNode
		n.children = append(n.children[:i], n.children[i+1:]...)
		n.items.removeAt(i - 1)
	} else {
		// 向 index = + 1 merge
		root := n.items[i+1]
		prev := n.children[i+1]
		next := n.children[i]
		newNode := &node{
			items:    append(prev.items, append(items{root}, next.items...)...),
			children: append(prev.children, prev.children...),
		}
		n.children[i] = newNode
		n.children = append(n.children[:i+1], n.children[i+2:]...)
		n.items.removeAt(i)
	}
}

func (n *node) rightRotation(i int) {
	// i 为 pivot item
	leftChild := n.children[i]
	rightChild := n.children[i+1]
	// 把 left 的最后一个 把 item 覆盖掉
	// 然后 left 最后一个删除
	n.items[i] = leftChild.items[len(leftChild.items)-1]
	leftChild.items.removeAt(len(leftChild.items) - 1)
	if len(leftChild.children) > 0 {
		// 移动 child
		waitMoveChild := leftChild.children[len(leftChild.children)-1]
		rightChild.children = append(children{waitMoveChild}, rightChild.children...)
		leftChild.children = leftChild.children[:len(leftChild.children)-1]
	}
}

func (n *node) leftRotation(i int) {
	// i 为 pivot item
	leftChild := n.children[i]
	rightChild := n.children[i+1]
	// 把 right 的第一个 把 item 覆盖掉
	// 然后 right 第一个删除
	n.items[i] = rightChild.items[0]
	rightChild.items.removeAt(0)
	if len(rightChild.children) > 0 {
		// 移动 child
		waitMoveChild := rightChild.children[0]
		leftChild.children = append(leftChild.children, waitMoveChild)
		rightChild.children = rightChild.children[1:]
	}
}
