package main

import (
	"bufio"
	"net"

	rsa "github.com/philippkeinberger/OwnProjects/goRSA"
)

// writeConn writes msg to the connection of client c
func writeConn(c *client, msg string) {
	if c.encryption {
		encrypted := rsa.EncryptBytes([]byte(msg), c.publicKey)
		send := []byte(string(encrypted) + ".")
		c.conn.Write(send)
	} else {
		c.conn.Write([]byte(msg + "."))
	}
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
