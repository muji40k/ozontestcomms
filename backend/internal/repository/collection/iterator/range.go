package iterator

type rangeIterator struct {
	i, end int
}

func (self *rangeIterator) Next() (int, bool) {
	if self.i >= self.end {
		return 0, false
	}

	out := self.i
	self.i++
	return out, true
}

func RangeIterator(start, end int) Iterator[int] {
	return &rangeIterator{min(start, end), max(start, end)}
}

