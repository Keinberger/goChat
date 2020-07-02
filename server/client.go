package main

import (
	"fmt"
	"math/big"
	"net"
	"strconv"

	rsa "github.com/philippkeinberger/OwnProjects/goRSA"
)

// client defines the data of a client
type client struct {
	nick          string
	currentRoomID int
	encryption    bool
	publicKey     rsa.PublicKey
	conn          net.Conn
}

// newClient() returns a new client having the standard data types inside
// and sends a welcome message to the clientconnection
func newClient(c net.Conn) *client {
	userName := "anonymous" + strconv.Itoa(userCount)
	fmt.Printf("New User(%v) joined\n", userName)
	cl := &client{
		nick:          userName,
		currentRoomID: 0,
		conn:          c,
	}
	writeConn(cl, conf.CommandMSG.Title+conf.CommandMSG.Text)
	clients = append(clients, cl)
	return cl
}

// setNick() method changes the nick of a client to string name
func (c *client) setNick(name string) {
	if rooms[c.currentRoomID].nickTaken(name) {
		writeConn(c, "Nick already in use")
		return
	}
	c.nick = name
	writeConn(c, "Changed nick to: "+name)
}

// join() method lets a client join chatRoom name
func (c *client) join(name string) {
	id := getRoomID(name)
	if id == 0 {
		writeConn(c, "Room not found")
		return
	}
	if rooms[id].nickTaken(c.nick) {
		writeConn(c, "Nick already in use")
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

// leave() method lets the client leave the current chatRoom
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

// isInsideRoom() method returns true if the client is inside of a chatRoom
func (c *client) isInsideRoom() bool {
	return c.currentRoomID != 0
}

// enableEncryption() method enables the RSA encryption for a client
func (c *client) enableEncryption(nn, aa int) {
	c.publicKey = rsa.PublicKey{
		N: big.NewInt(int64(nn)),
		A: big.NewInt(int64(aa)),
	}
	writeConn(c, "/servermessage keys "+strconv.Itoa(int(publicKey.N.Int64()))+" "+strconv.Itoa(int(publicKey.A.Int64())))
	c.encryption = true
	writeConn(c, "Traffic is now RSA encrypted")
}

// close() method closes the current client connection to the server
func (c *client) close() {
	if c.isInsideRoom() {
		c.leave()
	}
	fmt.Println(c.nick + " closed the connection")
	writeConn(c, "You closed the connection")
	c.conn.Close()
}
