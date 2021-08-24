package main

import (
	"net"
	"os"
	"time"
)
import "fmt"
import "bufio"
//import "strings" // only needed below for sample processing


func finishProgram(){
	// Waits 3 seconds before finishing the program
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func listener(port string, status *int) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Unavailable port %q\nFinishing program...",port)
		os.Exit(1)
	}
	fmt.Println("Program Listening in port", port)
	*status = 1

	// accept connection on port
	conn, _ := ln.Accept()

	// run loop forever (or until message stop is received)
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			finishProgram()
			break
		}
		// output message received
		if message == "stop\n"{
			fmt.Print("\nClient has decided to Stop communication.\nExiting program ...\n:")
			finishProgram()
		} else {
			fmt.Printf("\nMessage Received: %q\n", message[:len(message)-1])
		}
	}
}

func main() {

	fmt.Println("Launching Program ...")
	fmt.Print("Listening Port: ")

	// Reads from the stdin buffer the port where the program will be listening
	reader := bufio.NewReader(os.Stdin)
	listeningPort, _ := reader.ReadString('\n')
	listeningPort = listeningPort[:len(listeningPort)-2]

	// Listener Status {0: inactive, 1: active}
	status := 0 // The status variable let us know if the program is properly listening.
	go listener(listeningPort, &status)	// Program starts listening
	for status != 1{time.Sleep(time.Second)} // wait until the port is accepted and listening

	// Reads from the stdin buffer the port and ip of the target we want to connect
	fmt.Print("Server IP: ")
	serverIp, _ := reader.ReadString('\n')
	serverIp = serverIp[:len(serverIp)-2]
	fmt.Print("Server Port: ")
	serverPort, _ := reader.ReadString('\n')
	serverPort = serverPort[:len(serverPort)-2]

	// connect to this socket
	conn, err := net.Dial("tcp", serverIp+":"+serverPort)
	if err != nil {
		fmt.Printf("Server %q not ready. Waiting...\n", serverIp+":"+serverPort)
		timeOut := 0
		for err != nil{
			if timeOut > 10 {
				fmt.Printf("Wating Time exceeded. Ending Program ...\n")
				finishProgram()
			}	// If the program can't establish connection in an interval of 30 seconds, then ends the program
			time.Sleep(3 * time.Second)
			conn, err = net.Dial("tcp", serverIp+":"+serverPort)
			timeOut = timeOut + 1
		}
	}
	fmt.Printf("Program ready to stablish conversation with %q\n\n", serverIp+":"+serverPort)

	// If the connection is properly established, we get into a infinite loop where we can send a message for each iteration
	for {
		// read in input from stdin
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-2]

		// if the input string is "stop", then we notify to the server and finish the program
		if text == "stop"{
			fmt.Printf("You have decided to stop communication.\n Exiting program...\n")
			fmt.Fprintf(conn, "stop\n")
			finishProgram()
		}

		// Only those messages that start with "send:" will be send
		if len(text)>5{
			if text[:5] == "send:"{
				fmt.Fprintf(conn, text[5:len(text)]+"\n")	// send to socket
			} else{
				fmt.Printf("wrong message syntax. Try \"send:text\":")
			}
		} else{
			fmt.Printf("wrong message syntax. Try \"send:text\":")
		}

	}
}

