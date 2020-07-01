package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/philippkeinberger/OwnProjects/rsa/rsa"
)

func initialise() {
	leaveMessage = " left the chat"
	commandMessage = "List of commands:\n/nick <name> Change your nick\n/join <room> Join a room\n/room See current room\n/rooms List all rooms\n/leave Leave current room\n/createRoom <name> Create a new room\n/commands List all commands\n/encrypted Shows if conenction is encrypted\n"
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

func writeConn(c *client, msg string) {
	if c.encryption {
		encrypted := rsa.EncryptBytes([]byte(msg), c.publicKey)
		send := []byte(string(encrypted) + ".")
		c.conn.Write(send)
	} else {
		c.conn.Write([]byte(msg + "."))
	}
}

func handleConnection(conn net.Conn) {
	client := newClient(conn)

	// listen for user input sent to the server
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		go checkInput(input, client)
	}

	// encrypt traffic with RSA

	conn.Close()
}

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
			createRoom(client, split[1])
		case "leave":
			client.leave()
		case "close":
			client.close()
		case "commands":
			writeConn(client, commandMessage)
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
	leaveMessage   string
	commandMessage string
	wg             *sync.WaitGroup
	port           string = "9000"
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
