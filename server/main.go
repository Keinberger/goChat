package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"

	rsa "github.com/keinberger/goRSA"
)

// Config defines the structure of the config.json file
type Config struct {
	LeaveMessage string `json:"leaveMessage"`
	CommandMSG   `json:"commandMessage"`
}

// CommandMSG defines the structure of the welcome message, including the title and commandMessage
type CommandMSG struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// initialise() sets all default variables for the server
func initialise() {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		panic(err)
	}

	wg = &sync.WaitGroup{}
	room := room{
		name: "Home",
		chat: make(chan string),
	}
	rooms = append(rooms, room)

	privateKey = rsa.GeneratePrivateKey()
	publicKey = rsa.GetPublicKey(privateKey)

	fmt.Println("Succesfully started server on port: " + port)
}

// checkInput() checks the input string for client and handles it
func checkInput(input string, client *client) {
	if len(input) < 1 {
		return
	}

	fmt.Println(client.nick + " sent to server: " + input)
	if client.encryption {
		decrypted := rsa.DecryptBytes([]byte(input), privateKey, publicKey)
		input = string(decrypted)
	}

	if string(input[0]) != "/" {
		if client.isInsideRoom() {
			msg := client.nick + ": " + input
			writeConn(client, msg)
			// go sendRoomMessage(rooms[client.currentRoomID].chat, msg)
			go sendRoomMessageAnalog(client, msg, client.currentRoomID)
		} else {
			writeConn(client, "Unknown command")
		}
	} else {
		input = strings.Trim(input, "/")
		split := strings.Split(input, " ")
		switch strings.ToLower(split[0]) {
		case "enableencryption":
			if len(split) > 2 {
				n, _ := strconv.Atoi(split[1])
				a, _ := strconv.Atoi(split[2])
				client.enableEncryption(n, a)
			}
		case "nick":
			if len(split) == 1 {
				writeConn(client, client.nick)
			} else {
				client.setNick(split[1])
			}
		case "join":
			if len(split) == 1 {
				writeConn(client, "Choose channel: "+listRooms())
			} else {
				client.join(split[1])
			}
		case "room":
			writeConn(client, rooms[client.currentRoomID].name)
		case "rooms":
			writeConn(client, listRooms())
		case "createroom":
			if len(split) > 1 {
				createRoom(client, split[1])
			}
		case "leave":
			client.leave()
		case "close":
			client.close()
		case "commands":
			writeConn(client, conf.CommandMSG.Text)
		case "encrypted":
			writeConn(client, strconv.FormatBool(client.encryption))
		default:
			writeConn(client, "Unknown command")
		}
	}
}

var (
	rooms          []room
	clients        []*client
	publicKey      rsa.PublicKey
	privateKey     rsa.PrivateKey
	userCount      int
	wg             *sync.WaitGroup
	port           string = "9000"
	configFilePath string = "config.json"
	conf           Config
)

func main() {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	initialise()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		userCount++
		go handleConnection(conn)
	}
}
