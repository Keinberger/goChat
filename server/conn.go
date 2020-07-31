package main

import (
	"bufio"
	"net"

	rsa "github.com/keinberger/goRSA"
)

// writeConn writes msg to the connection of client c
func writeConn(c *client, msg string) {
	var sendMsg []byte
	if c.encryption {
		encrypted := rsa.EncryptBytes([]byte(msg), c.publicKey)
		sendMsg = []byte(string(encrypted) + ".")
	} else {
		sendMsg = []byte(msg + ".")
	}
	c.conn.Write(sendMsg)
}

// handleConnection() creates a new client checks for any messages on connection conn
func handleConnection(conn net.Conn) {
	client := newClient(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		input := scanner.Text()
		go checkInput(input, client)
	}

	conn.Close()
}
