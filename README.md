# goChat

## Requirements

The program requires a Unix System (Linux/macOS) or Go11 or newer.

## Usage

### Setting up a Server

Run 'go run main.go' inside of server/ and the program will automatically start a new server on port :9000 for you. It will log all the commands and/or messages that get sent to the server, same with deletion of chats and other important logs.

You may change the welcome message inside of config.json.

### Using the client

Run 'go run *.go ip port' (e.g. go run *.go localhost 9000) inside of client/. If the client succesfully connected to the server, you can see all the commands one can send to the server or use for the client.

Server commands start with '/' and will get sent to the server and be processed there.

Client commands start with '--'. It is important to note that the command will not get sent to the server, meaning, the input does not leave the client. 
However, '--enableEncryption' is an exception, because the programm will exchange the public-key for the RSA encryption with the server.

## Important

There might still be some bugs in the program.
