
package timer

import(
	"time"
)

type TimerItem struct{

	
	start 	time.Duration
	elapse 	time.Duration
	cb		CallBackFunc
	cancel	bool
	parm 	interface{}
}

func NewTimerItem(elapse time.Duration,cb CallBackFunc,parm interface{})  *TimerItem{
	
	return &TimerItem{
		start : time.Duration(time.Now().Unix()),
		elapse : elapse,
		cb : cb,
		cancel : false,
		parm : parm,
	}
}

func (this *TimerItem) leftTime() time.Duration  {
	
	spanTime := time.Now().Unix() - int64(this.start)
	leftTime := time.Duration(int64(this.elapse) - spanTime)

	if leftTime < 0{
		leftTime = 0
	}
	return leftTime
}

func (this *TimerItem) trigger()  {
	this.cb(this,this.parm)
}

func (this *TimerItem) reset(){
	this.start = time.Duration(time.Now().Unix())
}
func(this *TimerItem) Cancel(){
	this.cancel = true
}
func(this *TimerItem) canceled() bool{
	return this.cancel
}