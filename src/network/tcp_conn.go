package network

import(
	"net"
	"bufio"
	_"log"
	"sync"
	"encoding/binary"
	"bytes"
)

type ReadCallback func(socketID uint32, buf []byte,err error)

type TCPConn struct{
	conn net.Conn

	roundIndex,bindIndex uint16

	readBuf,writeBuf []byte

	readFunc ReadCallback

	reader *bufio.Reader
	writer *bufio.Writer

	active bool

	sendding,readding bool

	mtx sync.Mutex

	wg sync.WaitGroup
}

func newTCPConn(bindIndex uint16) *TCPConn{
	return &TCPConn{
		bindIndex : bindIndex,
		readBuf : make([]byte, 0),
		writeBuf : make([]byte, 0),
		active : false,
		sendding : false,
		readding : false,
		roundIndex : 1,
	}
}

func (this *TCPConn) Attach(conn net.Conn) {
	
	this.conn = conn
	this.reader = bufio.NewReader(this.conn)
	this.writer = bufio.NewWriter(this.conn)
	this.active = true
}
func (this *TCPConn) Detach(){

	this.roundIndex++
	if this.roundIndex == 0{
		this.roundIndex++
	}

	this.active = false
	this.reader = nil
	this.writer = nil
	this.conn = nil
	this.readFunc = nil
	this.sendding = false
	this.readding = false
}

func (this *TCPConn) IsActive() bool{
	return this.active
}

func (this *TCPConn) GetSocketID() uint32{
	return uint32(this.bindIndex) << 16 | uint32(this.roundIndex)
}
func (this *TCPConn) GetRoundIndex() uint16{
	return this.roundIndex
}
func(this *TCPConn) Recv(cb ReadCallback){

	this.readFunc = cb

	go func(){

		this.wg.Add(1)
		defer this.wg.Done()

		for{

			buf := make([]byte,1024)		
			this.readding = true
			n,err := this.reader.Read(buf)
			this.readding = false
			buf = buf[0:n]

			this.readBuf = append(this.readBuf,buf...)
			
			if err == nil{
				this.readFunc(this.GetSocketID(), nil,err)
				return
			}

			
			for len(this.readBuf) >= 4{

				var dataLen uint32
				err = binary.Read(bytes.NewReader(this.readBuf[0:4]),binary.LittleEndian,&dataLen)
				if err != nil{
					this.readFunc(this.GetSocketID(), nil,err)
					return
				}
			
				if len(this.readBuf) < int(dataLen + 4){
					continue
				}

				data := this.readBuf[4:4 + dataLen]
				this.readFunc(this.GetSocketID(), data,nil)
				this.readBuf = this.readBuf[4 + dataLen:]
			}

		}
		
	}()	
}
func(this* TCPConn) Send(buf []byte) bool{

	if buf != nil{

		lenBuf := new(bytes.Buffer)
		binary.Write(lenBuf,binary.LittleEndian,uint32(len(buf)))

		sendBuf := lenBuf.Bytes()
		sendBuf = append(sendBuf,buf...)

		this.mtx.Lock()
		this.writeBuf = append(this.writeBuf,sendBuf...)
		this.mtx.Unlock()
	}
	
	if this.sendding || len(this.writeBuf) == 0{
		return true
	}

	this.sendding = true

	go func(){

		this.wg.Add(1)
		defer this.wg.Done()

		this.mtx.Lock()
		tmpBuf := make([]byte,len(this.writeBuf))
		copy(tmpBuf,this.writeBuf)
		this.writeBuf = this.writeBuf[0:0]
		this.mtx.Unlock()

		this.writer.Write(tmpBuf)
		this.sendding = false

		this.Send(nil)
	}()

	return true
}

func(this *TCPConn) Close(){
	if this.IsActive() == false{
		return
	}
	
	this.conn.Close()
	this.wg.Wait()
}