package handle_net

import (
	"bufio"
	"bytes"
	"emperror.dev/errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type OpFunc func(*Message) (*Message, error)

func NewHandleServer(prefixes []string) (*HandleServer, error) {
	return &HandleServer{
		Prefixes:  prefixes,
		quit:      make(chan bool),
		exited:    make(chan bool),
		functions: make(map[OpCode]OpFunc),
	}, nil
}

type HandleServer struct {
	listener  net.Listener
	quit      chan bool
	exited    chan bool
	Prefixes  []string
	functions map[OpCode]OpFunc
}

func (srv *HandleServer) RegisterFunction(op OpCode, f OpFunc) {
	srv.functions[op] = f
}

func (srv *HandleServer) Start(protocol, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr(protocol, addr)
	if err != nil {
		return errors.Wrapf(err, "cannot resolve address %s %s", protocol, addr)
	}
	srv.listener, err = net.Listen(protocol, tcpAddr.String())
	if err != nil {
		return errors.Wrapf(err, "cannot listen on %s %s", protocol, tcpAddr.String())
	}
	go srv.Serve()
	return nil
}

func (srv *HandleServer) Serve() {
	var handlers sync.WaitGroup
	for {
		select {
		case <-srv.quit:
			fmt.Println("Shutting down...")
			srv.listener.Close()
			handlers.Wait()
			close(srv.exited)
			return
		default:
			//srv.listener.SetDeadline(time.Now().Add(1e9))
			conn, err := srv.listener.Accept()
			if err != nil {
				var opErr *net.OpError
				if errors.As(err, &opErr) && opErr.Timeout() {
					continue
				}
				fmt.Println("Failed to accept connection:", err.Error())
			}
			handlers.Add(1)
			go func() {
				defer handlers.Done()
				srv.handleConnection(conn)
			}()
		}
	}
}

func (srv *HandleServer) Stop() error {
	fmt.Println("Stop requested")
	close(srv.quit)
	<-srv.exited
	fmt.Println("Stopped successfully")
	return nil
}

func (srv *HandleServer) handleConnection(conn net.Conn) {
	log.Printf("Connection from %s", conn.RemoteAddr().String())
	defer log.Printf("Connection from %s closed", conn.RemoteAddr().String())
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buffer = bytes.NewBuffer(nil)
		var envelopeData = make([]byte, 20)
		var envelope *MessageEnvelope
		var i = 0
		var maxLength = 0
		for {
			b, err := reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println("read error:", err)
				return
			}
			if i < 20 {
				envelopeData[i] = b
			} else if i == 20 {
				envelope = &MessageEnvelope{}
				if err := envelope.UnmarshalHandleBinary(envelopeData); err != nil {
					log.Println("envelope unmarshal error:", err)
					return
				}
				maxLength = int(envelope.MessageLength) + 20
			}
			buffer.WriteByte(b)
			i++
			if i >= maxLength && maxLength > 0 {
				break
			}
		}
		if buffer.Len() == 0 {
			continue
		}
		message := &Message{}
		if err := message.UnmarshalHandleBinary(buffer.Bytes()); err != nil {
			log.Println("unmarshal error:", err)
		}
		if f, ok := srv.functions[message.Header.OpCode]; ok {
			resp, err := f(message)
			if err != nil {
				log.Println("error in function:", err)
			}
			if resp != nil {
				/*
					data, err := resp.MarshalHandleBinary()
					if err != nil {
						log.Println("error in marshal:", err)
					}
					conn.Write(data)

				*/
			}
		} else {
			log.Println("no function for opcode:", message.Header.OpCode)
		}
	}
}
