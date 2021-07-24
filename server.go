package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	Message   chan string
	mapLock   sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (this *Server) BroadCast(user *User, msg string) {
	message := "[" + user.Addr + "]" + user.Name + ": " + msg

	this.Message <- message
}

func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("connected")
	// conn.Write([]byte("You connected to the server"))
	user := NewUser(conn, this)

	user.UserOnline()

	isAlive := make(chan bool)

	// receive from cli and broadcast
	go func() {
		buffer := make([]byte, 4096)
		for {
			n, err := conn.Read(buffer)

			if n == 0 {
				user.UserOffline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn.Read error", err)
				return
			}
			msg := string(buffer[:n-1])

			user.MessageHandler(msg)

			isAlive <- true
		}
	}()

	for {
		select {
		case <-isAlive:

		case <-time.After(time.Minute * 10):
			user.SendToUser("No activity for more than 10 minutes, you've been removed\n")
			user.UserOffline()

			close(user.C)
			conn.Close()

			return //exit this handler
		}
	}
}

func (this *Server) ListenMessage() {
	for {
		Message := <-this.Message
		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- Message
		}
		this.mapLock.Unlock()
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error: ", err)
		return
	}

	defer listener.Close()

	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listner Accept error: ", err)
			continue
		}

		go this.Handler(conn)
	}

}
