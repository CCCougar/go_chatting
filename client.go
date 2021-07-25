package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	conn       net.Conn
	Name       string
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// create object
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// connnect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}
	client.conn = conn

	//return the created obj
	return client
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP setting")
	flag.IntVar(&serverPort, "port", 8888, "Server Port setting")
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println(">>>>>>>>>>>Menu<<<<<<<<<<<")
	fmt.Println("1. Public Chatting")
	fmt.Println("2. Private Chatting")
	fmt.Println("3. Name Change")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Please input valid number")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println(">>>>>>Public Chatting<<<<<<")
			client.PublicChatting()
			break
		case 2:
			fmt.Println(">>>>>>Private Chatting<<<<<<")
			client.PrivateChatting()
			break
		case 3:
			fmt.Println(">>>>>>Name Change<<<<<<")
			client.ChangeName()
			break
		}
	}
}

func (client *Client) PublicChatting() {
	var clientMsg string

	fmt.Println("Please input your message(WITH NO **SPACE** IN YOUR MESSAGE, input \"exit\" to quit)")
	fmt.Scanln(&clientMsg)

	// fmt.Println("clientMsg: " + clientMsg)

	for clientMsg != "exit" {
		if len(clientMsg) > 0 {
			// fmt.Println("clientMsg: " + clientMsg)
			sendMsg := clientMsg + "\n"
			// fmt.Println("sendMsg:" + sendMsg)
			_, err := client.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("conn.Write error: ", err)
				break
			}
		}

		clientMsg = ""

		fmt.Println("Please input your message(end by input \"exit\")")
		fmt.Scanln(&clientMsg)
		// fmt.Println("clientMsg: " + clientMsg)
	}
}

func (client *Client) ShowOnlineUsers() {
	client.conn.Write([]byte("who is online\n"))
}

func (client *Client) PrivateChatting() {
	client.ShowOnlineUsers()

	var targetUser string

	fmt.Println(">>>>>>>>>Please choose which one to chat:<<<<<<<<<")
	fmt.Scanln(&targetUser)
	var clientMsg string

	fmt.Println("Please input your message(WITH NO **SPACE** IN YOUR MESSAGE, input \"exit\" to quit)")
	fmt.Scanln(&clientMsg)

	for clientMsg != "exit" {
		if len(clientMsg) > 0 {
			sendMsg := "to|" + targetUser + "|" + clientMsg + "\n"
			// fmt.Println(sendMsg) 		// for debug
			_, err := client.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("conn.Write error: ", err)
				break
			}
		}

		clientMsg = ""

		fmt.Println("Please input your message(WITH NO **SPACE** IN YOUR MESSAGE, input \"exit\" to quit)")
		fmt.Scanln(&clientMsg)
	}

}

func (client *Client) ChangeName() {
	fmt.Println("Please input your newname:")
	fmt.Scanln(&client.Name)

	client.conn.Write([]byte("change my name to " + client.Name + "\n"))
}

// Deal the message returned from server
func (this *Client) Dealresponse() {
	io.Copy(os.Stdout, this.conn) // Permanent block and copy&update
}

func main() {
	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("***Client create faild***")
		return
	}

	fmt.Println("***Client create succeed***")

	go client.Dealresponse()

	client.Run()
}
