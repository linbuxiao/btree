package btree

func getMiddle(arr []int) int {
	return len(arr) / 2
}

func itemsIsFull(arr []int) bool {
	return len(arr) == M
}
