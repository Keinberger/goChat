package main

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	rsa "github.com/philippkeinberger/OwnProjects/goRSA"
)

func executeCommand(cmd string, conn net.Conn) {
	cmd = strings.TrimSpace(cmd)
	switch strings.ToLower(cmd) {
	case "clear":
		clear()
	case "enableencryption":
		if !encryption {
			enableEncryption(conn)
		}
	default:
		fmt.Print("Unknown client command\n")
		fmt.Print(inputPhrase)
	}
	fmt.Print(inputPhrase)
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func setServerKeys(keys string) {
	split := strings.Split(keys, " ")
	nn, _ := strconv.Atoi(split[2])
	aa, _ := strconv.Atoi(split[3])
	serverPublicKey = rsa.PublicKey{
		N: big.NewInt(int64(nn)),
		A: big.NewInt(int64(aa)),
	}
	encryption = true
}

func enableEncryption(conn net.Conn) {
	privateKey = rsa.GeneratePrivateKey()
	publicKey = rsa.GetPublicKey(privateKey)
	fmt.Fprintf(conn, "/enableEncryption "+strconv.Itoa(int(publicKey.N.Int64()))+" "+strconv.Itoa(int(publicKey.A.Int64()))+"\n")
}

func handleReply(conn net.Conn) {
Y:
	for {
		message, _ := bufio.NewReader(conn).ReadString('.')
		message = strings.Trim(message, ".")
		if encryption {
			message = string(rsa.DecryptBytes([]byte(message), privateKey, publicKey))
		}

		switch {
		case strings.Contains(message, "Joined room: "):
			//inputPhrase = "Send text: "
			clear()
		case strings.Contains(message, "You left"):
			//inputPhrase = "Command: "
			clear()
		case strings.Contains(message, "/servermessage keys"):
			setServerKeys(message)
			continue Y
		case strings.Contains(message, "You closed the connectio"):
			inputPhrase = "Press enter to exit the program"
			exit = true
			clear()
			break Y
		}
		if message != "" {
			fmt.Print(message + "\n")
			fmt.Print(inputPhrase)
		}
	}
	fmt.Print(inputPhrase)
}

func sendToServer(conn net.Conn, text string) {
	if encryption {
		bef := strings.Trim(text, "\n")
		encr := rsa.EncryptBytes([]byte(bef), serverPublicKey)
		encrypted := string(encr) + "\n"

		fmt.Fprintf(conn, encrypted)
	} else {
		fmt.Fprintf(conn, text)
	}
}

var (
	exit            bool
	inputPhrase     string
	encryption      bool
	privateKey      rsa.PrivateKey
	publicKey       rsa.PublicKey
	serverPublicKey rsa.PublicKey
)

func main() {
	if len(os.Args) < 2 {
		panic("You have to specify the server and port")
	}
	args := os.Args[1:]
	conn, err := net.Dial("tcp", args[0]+":"+args[1])
	if err != nil {
		panic("Could not connect to server (" + args[0] + ":" + args[1] + ")")
	}
	defer conn.Close()

	inputPhrase = "> "
	clear()

	go handleReply(conn)
	for {
		if exit {
			break
		}
		reader := bufio.NewReader(os.Stdin)
		//clear()
		//fmt.Print(inputPhrase)
		text, _ := reader.ReadString('\n')
		if len(text) > 2 {
			possibleCommand := string(text[:2])
			if possibleCommand == "--" {
				executeCommand(string(text[2:]), conn)
				continue
			}
		}
		// send to server
		sendToServer(conn, text)
	}
}
