package network

import(
	"net"
	"log"
	"time"
	"sync"
	"fmt"
)

type ITCPServer interface{
	OnNewConn(socketID uint32)
	OnNewMessage(socketID uint32,data []byte)
	OnConnClose(socketID uint32)
}
type TCPServer struct{
	ln		net.Listener
	conns	[]*TCPConn
	handler ITCPServer

	wg		sync.WaitGroup
}

func NewTCPServer(handler ITCPServer,maxConn int) *TCPServer{

	p := &TCPServer{
		conns : make([]*TCPConn, 0),
		handler : handler,
	}

	for i:=0;i < maxConn;i++{
		p.conns = append(p.conns,newTCPConn(uint16(i)))
	}

	return p;
}

func (this *TCPServer) createConn() *TCPConn{
	for _,conn := range this.conns{
		if conn.IsActive() {
			continue
		}
		return conn
	}
	return nil
}

func(this *TCPServer) Start(port int) bool{

	host := fmt.Sprintf(":%d",port)
	ln,err := net.Listen("tcp",host)
	if err != nil{
		log.Printf("Listen:%s",err)
		return false
	}

	this.ln = ln

	go func(){

		this.wg.Add(1)
		defer this.wg.Done()

		for{
			conn,err := this.ln.Accept()
			var tempDelay time.Duration
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					log.Printf("accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				return
			}

			tempDelay = 0	
			tcpConn := this.createConn()

			if tcpConn == nil{
				conn.Close()
				continue
			}


			log.Println("new conn")

			tcpConn.Attach(conn)
			tcpConn.Recv(this.OnConnRead)

			this.handler.OnNewConn(tcpConn.GetSocketID())
		}
	}()
	

	return true
}

func(this *TCPServer) OnConnRead(socketID uint32,buf []byte,err error){
	if err != nil{
		this.handler.OnConnClose(socketID)
		conn := this.getConn(socketID)
		if conn != nil{
			conn.Detach()
		}
		return
	}
	this.handler.OnNewMessage(socketID,buf)
}

func(this *TCPServer)getConn(socketID uint32) *TCPConn{

	bindIndex  := (socketID & 0xffff0000) >> 16
	roundIndex := socketID & 0x0000ffff

	if bindIndex > uint32(len(this.conns) - 1){
		return nil
	}
	conn := this.conns[bindIndex]
	if conn.GetRoundIndex() != uint16(roundIndex){
		return nil
	}

	return conn
}

func(this *TCPServer)Close(){

	for _,conn := range this.conns{
		conn.Close()
	}

	this.ln.Close()
	this.wg.Wait()
}

func(this *TCPServer)SendData(socketID uint32, data []byte){

	conn := this.getConn(socketID)

	if conn == nil{
		return
	}

	conn.Send(data)
}

func (this *TCPServer) CloseSocket(socketID uint32)  {

	conn := this.getConn(socketID)

	if conn == nil{
		return
	}
	conn.Close();
}