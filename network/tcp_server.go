package network

import(
	"net"
	"log"
	"time"
	"fmt"
)
type TCPServer struct{
	ln		net.Listener
	conns	[]*TCPConn
}

func NewTCPServer() *TCPServer{

	p := &TCPServer{
		conns : make([]*TCPConn, 0),
	}

	for i:=0;i<512;i++{
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

func(this *TCPServer) Start(port,backlog int) bool{
	ln,err := net.Listen("tcp","127.0.0.1:8888")
	if err != nil{
		log.Printf("Listen:%s",err)
		return false
	}

	this.ln = ln

	go func(){
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
			
			log.Println("new conn")
			tcpConn := this.createConn()

			tcpConn.Attach(conn)
			tcpConn.Recv(this.OnConnRead)
			
			str := "77272367263"
			tcpConn.Send([]byte(str))
			tcpConn.Close()

		}
	}()
	

	return true
}

func(this *TCPServer) OnConnRead(socketID uint32,buf []byte,err error){
	if err != nil{
		fmt.Printf("close")
		return
	}

	fmt.Printf("%s",string(buf))
}