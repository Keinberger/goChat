package main

import (
	"fmt"
	"strings"
	"time"
)

type room struct {
	name    string
	clients []*client
	chat    chan string
}

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
	return roo
}

func checkRoomMessages(client *client) {
Y:
	for {
		select {
		case msg := <-rooms[client.currentRoomID].chat:
			writeConn(client, msg)
			if msg == client.nick+leaveMessage {
				break Y
			}
		}
	}
}

func sendRoomMessage(ch chan string, msg string) {
	ch <- msg
}

func sendRoomMessageAnalog(sender *client, msg string, roomID int) {
	for _, c := range rooms[roomID].clients {
		if c.nick == "" || sender == c {
			continue
		}
		go writeConn(c, msg)
	}
}
