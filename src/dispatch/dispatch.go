package dispatch

import(
	"sync"
)

type NetDataQue struct{
	dataType int
	socketID uint32
	data []byte
}

type EventDispatch struct{

	netDataQue []*NetDataQue
	netMtx		sync.Mutex
}

const(
	NET_DATATYPE_NEWCONN int = iota
	NET_DATATYPE_NEWMESSAGE
	NET_DATATYPE_CLOSE
)

func NewEventDispatch() *EventDispatch{
	return &EventDispatch{
		netDataQue : make([]*NetDataQue,0),
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

	
}

