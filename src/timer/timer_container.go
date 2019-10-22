package timer

type TimerContainer []*TimerItem

func (h TimerContainer) Len() int           { return len(h) }
func (h TimerContainer) Less(i, j int) bool { return h[i].leftTime() < h[j].leftTime() }
func (h TimerContainer) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TimerContainer) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*TimerItem))
}

func (h *TimerContainer) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func(h *TimerContainer) Front() *TimerItem{
	if len(*h) == 0{
		return nil
	}
	return (*h)[0]
}