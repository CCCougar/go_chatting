package main

import (
	"net"
	"strings"
)

type User struct {
	Conn   net.Conn
	Name   string
	Addr   string
	C      chan string
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	Address := conn.RemoteAddr().String()

	user := &User{
		Conn:   conn,
		Name:   Address,
		Addr:   Address,
		C:      make(chan string),
		server: server,
	}

	go user.ListenC()

	return user
}

func (this *User) ListenC() {
	for {
		mess := <-this.C
		this.Conn.Write([]byte(mess + "\n"))
	}
}

//User Online
func (this *User) UserOnline() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "***LOGGED***")
}

//User Offline
func (this *User) UserOffline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "***LOGGED OUT***")
}

// Handle user messages
func (this *User) MessageHandler(msg string) {
	// Check which users are online
	if msg == "who is online" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + " ***is online***\n"
			this.SendToUser(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 17 && msg[:17] == "change my name to" {
		newName := strings.Fields(msg)[4]
		// fmt.Println(newName)

		// whether the name is being used
		_, ok := this.server.OnlineMap[newName]
		if ok == false {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendToUser("You've change your name to " + newName + "\n")
		} else {
			nameExistMsg := newName + "exists \n"
			this.SendToUser(nameExistMsg)
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		// 1. get the message to send
		msgToSend := strings.Split(msg, "|")[2]

		// 2. get the target user pointer
		targetUserName := strings.Split(msg, "|")[1]
		targetUser, ok := this.server.OnlineMap[targetUserName]
		if ok == false {
			this.SendToUser("User: " + targetUserName + " does't online!")
			return
		}

		// 3. send message
		targetUser.SendToUser(this.Name + " said to you: " + msgToSend + "\n")
		this.SendToUser("Your message to " + targetUserName + " is been send\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) SendToUser(msg string) {
	this.Conn.Write([]byte(msg))
}
