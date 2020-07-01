package main

import (
	"fmt"
	"strings"
	"time"
)

// room defines the data structure for a room
type room struct {
	name    string
	clients []*client
	chat    chan string
}

// getRoomID returns the ID of room name
func getRoomID(name string) int {
	var key int
	for k, v := range rooms {
		if v.name == name {
			key = k
			break
		}
	}
	return key
}

// createRoom creates a new room named nam
func createRoom(c *client, nam string) {
	if strings.Contains(listRooms(), nam) {
		writeConn(c, "Room already exists")
		return
	}
	roo := room{
		name: nam,
		chat: make(chan string),
	}
	rooms = append(rooms, roo)
	writeConn(c, "Created room: "+nam)

	go deleteRoomIfEmpty(getRoomID(nam))
}

// deleteRoomIfEmpty() deletes a room with id roomID if it is empty
func deleteRoomIfEmpty(roomID int) {
	time.Sleep(time.Second * 30)
	userInside := false
	for {
		for _, v := range rooms[roomID].clients {
			if v.nick != "" {
				userInside = true
				break
			}
		}
		if !userInside {
			break
		}
	}
	fmt.Println(rooms[roomID].name + " got deleted")
	rooms[roomID] = room{}
}

// listRoom() returns a list of all rooms
func listRooms() string {
	var roo string
	for k, v := range rooms {
		if k == 0 {
			continue
		}
		roo += v.name + ", "
	}
	if roo == "" {
		return "No rooms found"
	}
	roo = strings.TrimSpace(roo)
	roo = strings.Trim(roo, ",")
	return roo
}

// checkRoomMessages() checks the channel of a room for messages and sends the message to the client.conn
func checkRoomMessages(client *client) {
Y:
	for {
		select {
		case msg := <-rooms[client.currentRoomID].chat:
			writeConn(client, msg)
			if msg == client.nick+conf.LeaveMessage {
				break Y
			}
		}
	}
}

// sendRoomMessage() send a message msg to channel ch
func sendRoomMessage(ch chan string, msg string) {
	ch <- msg
}

// sendRoomMessageAnalog() sends a mesage msg of room with id roomID to all clients in the room
func sendRoomMessageAnalog(sender *client, msg string, roomID int) {
	for _, c := range rooms[roomID].clients {
		if c.nick == "" || sender == c {
			continue
		}
		go writeConn(c, msg)
	}
}
