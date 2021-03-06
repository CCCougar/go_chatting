package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	ServerIp   string
	ServerPort int
	conn       net.Conn
	Name       string
	flag       int
}

// create a client object
func NewClient(serverIp string, serverPort int) *Client {
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
	client.ChangeName()
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
	// var clientMsg string

	// fmt.Println("Please input your message(WITH NO **SPACE** IN YOUR MESSAGE, input \"exit\" to quit)")
	// fmt.Scanln(&clientMsg)
	fmt.Print("Please input your message(input \"exit\" to quit) >>> ")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
	// fmt.Println("clientMsg: " + input)

	for input != "exit\x0a" {
		if len(input) > 0 {
			// sendMsg := clientMsg + "\n"
			sendMsg := input + "\n"
			// fmt.Println("sendMsg:" + sendMsg)
			_, err := client.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("conn.Write error: ", err)
				break
			}
		}

		// clientMsg = ""
		input = ""

		fmt.Print("Please input your message(input \"exit\" to quit) >>> ")
		// fmt.Scanln(&clientMsg)
		input, _ = inputReader.ReadString('\n')
		// fmt.Printf("input: %x", input)
		// time.Sleep(time.Millisecond * 100)
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
	fmt.Println("Please give yourself a name:")
	fmt.Scanln(&client.Name)

	client.conn.Write([]byte("change my name to " + client.Name + "\n"))
}

// Deal the message returned from server
func (this *Client) Dealresponse() {
	file, err := os.OpenFile("Response_from_server.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	// f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)
	file.Write([]byte("------------" + time.Now().Local().Format("Mon Jan 2 15:04:05 -0700 MST 2006") + "--------------\n"))
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Printf("os.Create error: %s", err)
	}

	defer file.Close()

	io.Copy(file, this.conn) // Permanent block and copy&update
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
