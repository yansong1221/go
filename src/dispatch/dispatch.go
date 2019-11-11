package dispatch

import(
	"sync"
	"threadpool"
	"timer"
)

type IEvent interface{

	//网络事件
	OnSocketConn(socketID uint32)
	OnSocketMessage(socketID uint32,data []byte)
	OnSocketClose(socketID uint32)
}

type NetDataQue struct{
	dataType int
	socketID uint32
	data []byte
}

type EventDispatch struct{

	netDataQue []*NetDataQue
	netMtx		sync.Mutex
	event  		IEvent


	tp 			*threadpool.ThreadPool
	tm 			*timer.TimerManager
}

const(
	NET_DATATYPE_NEWCONN int = iota
	NET_DATATYPE_NEWMESSAGE
	NET_DATATYPE_CLOSE
)

func NewEventDispatch(e IEvent) *EventDispatch{
	return &EventDispatch{
		netDataQue : make([]*NetDataQue,0),
		event : e,
		tp : threadpool.New(),
		tm : timer.New(),
	}
}

func (this *EventDispatch) OnNewConn(socketID uint32)  {
	this.netMtx.Lock()

	this.netDataQue = append(this.netDataQue,&NetDataQue{dataType : NET_DATATYPE_NEWCONN,
		socketID : socketID,
		data : nil})

	this.netMtx.Unlock()
}
func  (this *EventDispatch) OnNewMessage(socketID uint32,data []byte)  {
	
	this.netMtx.Lock()

	this.netDataQue = append(this.netDataQue,&NetDataQue{dataType : NET_DATATYPE_NEWMESSAGE,
		socketID : socketID,
		data : data})

	this.netMtx.Unlock()
}
func (this *EventDispatch) OnConnClose(socketID uint32)  {

	this.netMtx.Lock()

	this.netDataQue = append(this.netDataQue,&NetDataQue{dataType : NET_DATATYPE_CLOSE,
		socketID : socketID,
		data : nil})

	this.netMtx.Unlock()
}


func (this *EventDispatch) Update()  {
	
	this.netMtx.Lock()

	tempNetData := make([]*NetDataQue,len(this.netDataQue))
	copy(tempNetData,this.netDataQue)
	this.netDataQue = this.netDataQue[0:0]

	this.netMtx.Unlock()

	//网络事件
	for _,data := range tempNetData{
		switch data.dataType {
		case NET_DATATYPE_NEWCONN:
			this.event.OnSocketConn(data.socketID)
		case NET_DATATYPE_NEWMESSAGE:
			this.event.OnSocketMessage(data.socketID,data.data)
		case NET_DATATYPE_CLOSE:
			this.event.OnSocketClose(data.socketID)		
		}
	}
	
	this.tp.Update()
	this.tm.Update()
}
func(this *EventDispatch) Close(){
	this.netDataQue = this.netDataQue[0:0]
	this.tp.Close()
	this.tm.Close()
}

func(this *EventDispatch) GetTimer() *timer.TimerManager{
	return this.tm
}

func(this *EventDispatch) GetThreadPool() *threadpool.ThreadPool{
	return this.tp
}
