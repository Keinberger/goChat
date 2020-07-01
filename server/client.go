package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/philippkeinberger/OwnProjects/rsa/rsa"
)

type client struct {
	nick          string
	currentRoomID int
	encryption    bool
	publicKey     rsa.PublicKey
	conn          net.Conn
}

func newClient(c net.Conn) *client {
	userName := "anonymous" + strconv.Itoa(userCount)
	fmt.Printf("New User(%v) joined\n", userName)
	cl := &client{
		nick:          userName,
		currentRoomID: 0,
		conn:          c,
	}
	writeConn(cl, "Succesfully connected to Safe Messanger\n"+commandMessage)
	clients = append(clients, cl)
	return cl
}

func (c *client) setNick(name string) {
	for _, v := range rooms[c.currentRoomID].clients {
		if v.nick == name {
			writeConn(c, "Nick already in use")
			return
		}
	}
	c.nick = name
	writeConn(c, "Changed nick to: "+name)
}

func (c *client) join(name string) {
	id := getRoomID(name)
	if id == 0 {
		writeConn(c, "Channel not found")
		return
	}
	c.currentRoomID = id
	rooms[c.currentRoomID].clients = append(rooms[c.currentRoomID].clients, c)

	i := 0
	users := ""
	for _, v := range rooms[c.currentRoomID].clients {
		if v.nick != "" {
			users += v.nick + ", "
			i++
		}
	}
	users = string(users[:len(users)-2])

	// send message to channel (analog)
	sendRoomMessageAnalog(c, c.nick+" joined the room", c.currentRoomID)

	// send message to channel
	// sendRoomMessage(rooms[c.currentRoomID].chat, c.nick + " joined the room")

	writeConn(c, "Joined room: "+name+"\nCurrent users ("+strconv.Itoa(i)+"): "+users)

	go checkRoomMessages(c)
}

func (c *client) leave() {
	if c.isInsideRoom() {
		var key int
		for k, v := range rooms[c.currentRoomID].clients {
			if v == c {
				key = k
			}
		}
		rooms[c.currentRoomID].clients[key] = &client{}

		writeConn(c, "You left: "+rooms[c.currentRoomID].name)

		// send message to channel (analog)
		sendRoomMessageAnalog(c, c.nick+" left the room", c.currentRoomID)

		// send message to channel
		// sendRoomMessage(rooms[c.currentRoomID].chat, c.nick + " left the channel")

		c.currentRoomID = 0
	} else {
		writeConn(c, "You are not inside a room")
	}
}

func (c *client) isInsideRoom() bool {
	return c.currentRoomID != 0
}

func (c *client) enableEncryption(nn, aa int) {
	c.publicKey = rsa.PublicKey{
		N: nn,
		A: aa,
	}
	writeConn(c, "/servermessage keys "+strconv.Itoa(publicKey.N)+" "+strconv.Itoa(publicKey.A))
	c.encryption = true
	writeConn(c, "Traffic is now RSA encrypted")
}

func (c *client) close() {
	if c.isInsideRoom() {
		c.leave()
	}
	fmt.Println(c.nick + " closed the connection")
	writeConn(c, "You closed the connection")
	c.conn.Close()
}
