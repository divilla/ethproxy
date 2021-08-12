package ethcache

type (
	expiries []*item
)

func (i expiries) Len() int {
	return len(i)
}

func (i expiries) Less(x, y int) bool {
	return i[x].expires < i[y].expires
}

func (i expiries) Swap(x, y int) {
	i[x], i[y] = i[y], i[x]
}
