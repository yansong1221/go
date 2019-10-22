package timer

import(
	"time"
	"container/heap"
)

type CallBackFunc func(*TimerItem,interface{})

type TimerManager struct{

	times	TimerContainer
}

func New()  *TimerManager{

	return &TimerManager{

	}	
}

func(this *TimerManager) AddTimer(d time.Duration, cb CallBackFunc,parm interface{})  {
	
	item := NewTimerItem(d,cb,parm)
	heap.Push(&this.times,item)
}

func(this *TimerManager) Update(){

	if len(this.times) == 0{
		return
	}

	for{
		item := this.times.Front()
		if item.leftTime() != 0{
			return
		}
		heap.Pop(&this.times)

		if item.canceled() == false{

			item.trigger()
			item.reset()

			heap.Push(&this.times,item)
		}	
	}
}