package main

type MaxFn (func() int)

type Index struct {
	max MaxFn
	idx int
}

func NewIndex(maxRows MaxFn) *Index {
	return &Index{
		max: maxRows,
		idx: 0,
	}
}

func (i *Index) IsSelected(idx int) bool {
	return i.idx == idx
}

func (i *Index) Cur() int {
	return i.idx
}

func (i *Index) Inc() {
	i.idx = (i.idx + 1) % i.max()
}

func (i *Index) Dec() {
	if i.idx == 0 {
		i.idx = i.max() - 1
	} else {
		i.idx = (i.idx - 1) % i.max()
	}
}
