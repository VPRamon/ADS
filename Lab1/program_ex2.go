package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
)
import "fmt"
import "bufio"



func finishProgram(){
	// Waits 3 seconds before finishing the program
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func newConnexion(conn net.Conn){
	/* 	This function is called every time a new client requests a connection.
		Here we receive  every message from this `particular client   			*/

	// run loop forever (or until "stop" message is received)
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
			fmt.Printf("Message Received: %q\n", message[:len(message)-1])
		}
	}
}

func listener(port string, status *int) {
	//Starts listening in the specified port
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Unavailable port %q\nFinishing program...",port)
		finishProgram()
	}
	fmt.Println("Program Listening in port", port)
	*status = 1
	// For each connection request we create a new connection to listen from a different buffer
	for {
		conn, _ := ln.Accept()
		go newConnexion(conn)
	}
}

func broadcast(message string, adrSlice []net.Conn) {
	for i:=0;i<len(adrSlice);i++ {		// Sends message to all neighbours
		fmt.Fprintf(adrSlice[i], message)
	}
}

func establishConnection(adr string) (net.Conn, error) {
	// connect to this socket
	conn, err := net.Dial("tcp", adr)
	if err != nil {
		fmt.Printf("Server %q seems down. Waiting...\n", adr)
		timeOut := 0
		for err != nil{
			if timeOut > 30 {
				fmt.Printf("Wating Time exceeded for %q. stopping request...\n", adr)
				return nil, err
			}
			time.Sleep(time.Second)
			conn, err = net.Dial("tcp", adr)
			timeOut = timeOut + 1
		}
	}
	fmt.Printf("%q is now ready. Establishing connection\n", adr)
	return conn, nil
}

func main() {

	fileName := os.Args[1]

	fmt.Println("Launching Program ...")
	fmt.Printf("ConfigFile name: %q\n",fileName)

	// Opens the ConfigFile to read
	file, err := os.Open(fileName)
	if err != nil {	log.Fatal(err) }
	defer file.Close()


	adrSlice := []string{}	// Contain all the ConfigFile addresses
	scanner := bufio.NewScanner(file)
	// Reads all the addresses of the file and sotores them in the "adrSlice" Slice
	for scanner.Scan(){
		adr := scanner.Text()
		adrSlice = append(adrSlice,adr)
	} // Scans the file and stores the addresses in adrSlice


	selfAddress := strings.Split(adrSlice[0], ":")	// Splits the first address (our own address) to get the port
	port := selfAddress[1]
	status := 0		// Listener Status {0: inactive, 1: active}
	go listener(port, &status)
	for status != 1{time.Sleep(time.Second)}	// Waits until program starts listening

	connSlice := []net.Conn{}	// Stores the active connections
	// Establishing all connections and stores the Slice var "connSlice"
	for i:=1;i<len(adrSlice);i++ {
		conn, err := establishConnection(adrSlice[i])
		if err!=nil{ continue }
		connSlice = append(connSlice,conn)
	}

	reader := bufio.NewReader(os.Stdin)
	// Infinite loop where we can broadcast a message for each iteration
	for {
		// read in input from stdin
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-2]

		// if the input string is "stop", then we notify to the server and finish the program
		if text == "stop"{
			fmt.Printf("You have decided to stop communication.\n Exiting program...\n")
			broadcast("stop\n",connSlice)
			finishProgram()
		}

		// Only those messages that start with "send:" will be send
		if len(text)>5{
			if text[:5] == "send:"{
				broadcast(text[5:len(text)]+"\n",connSlice)
			} else{
				fmt.Printf("wrong message syntax. Try \"send:text\":")
			}
		} else{
			fmt.Printf("wrong message syntax. Try \"send:text\":")
		}

	}
}


