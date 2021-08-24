package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"strconv"
	"sync"
	"time"
)

var myStringValue string
var myValue int
var parentExists bool
var parentConn net.Conn
var parentAddress string
var parentId string
var connSlice []net.Conn 	// Stores the active connections
var listenerStatus int		// 0:passive  1:active
var mutex = &sync.Mutex{}	// Synchronize go routines

func initialize(fileName string) []string{
	// Opens the ConfigFile to read
	file, err := os.Open(fileName)
	if err != nil {	log.Fatal(err) }
	defer file.Close()

	adrSlice := []string{}	// Contain all the ConfigFile addresses
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		adr := scanner.Text()
		adrSlice = append(adrSlice,adr)
	} // Scans the file and stores the addresses in adrSlice
	return adrSlice
}

func finishProgram(){
	// Waits 3 seconds before finishing the program
	//we use this function to avoid
	time.Sleep(3 * time.Second)
	os.Exit(0)
}

func askNeighbour(message string, conn net.Conn, c chan int) {
	/* In this function we read the information send by the neighbours,that is why we need
	   to read in the connection the information sent by the neighbour. If the message sent is
	   a 0(we also add a \n to the parsing as in all the message this character is sent after the body of the message)
	   we then sent the int 0 to the channel, if it is 1, we send 1.	*/

	fmt.Fprintf(conn, message+"\n")
	response, _ := bufio.NewReader(conn).ReadString('\n')
	val, _ := strconv.Atoi(response[:len(response)-1])
	c <- val
}

func newConnexion(conn net.Conn){
	/* 	This function is called every time a new client requests a connection.
		Here we receive  every message from this particular client   			*/

	// run loop forever (or until "stop" message is received)
	for {
		// Will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Print("Error while reading message\n")
			finishProgram()
			break
		}

		/* We use then clause/command of stop to finish connection, therefore in every message sent,
		   we check if the message does indeed has the form of "stop"								*/
		if message == "stop\n"{
			mutex.Lock()
			if listenerStatus==1 {
				listenerStatus = 0
				broadcast(message, connSlice)
				fmt.Print("\nClient has decided to Stop communication.\nExiting program ...\n")
			}
			mutex.Unlock()
			finishProgram()

		/* If the message is not stop then it has a real value to the Algorithm and it could be in
		   the case that is the first message, the id of the parent. For this we will have a bool
		   global variable called 'parentExists' that will be set to True once we meet our parent.
		   Then the process becomes active.															*/
		// If we already have a parent, it means that some neighbour is requesting our value.
		}else{
			info := strings.Split(message, ":")
			ip := info[0]
			port := info[1]
			id := info[2]
			//fmt.Printf("Request value from %q",id)
            if parentExists == false {
				parentAddress = ip+":"+port
				parentId = id[:len(id)-1]
                parentConn = conn
				parentExists = true
            }else{
                fmt.Fprintf(conn, myStringValue+"\n")
            }
		}
	}
}

func listener(address string, status *int) {
	/* Starts listening in the specified port */

	ln, err := net.Listen("tcp", ":"+address)
	if err != nil {
		fmt.Printf("Unavailable port %q\nFinishing program...\n",address)
		finishProgram()
	}
	fmt.Printf("Program Listening in port %q\n", address)
	*status = 1
	// For each connection request we create a new connection to listen from a different buffer
	for{
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
			if timeOut > 5 {
				fmt.Printf("Wating Time exceeded for %q. stopping request...\n", adr)
				return nil, err
			}
			time.Sleep(time.Second)
			conn, err = net.Dial("tcp", adr)
			timeOut = timeOut + 1
		}
	}
	//fmt.Printf("%q is now ready. Establishing connection\n", adr)
	return conn, nil
}


func decision(c chan int, channelSize int) int{
	if channelSize < 1{
		fmt.Printf("I have no neighbours to ask.\n")
		return myValue
	}
	counter := 0
	for i:=0;i<channelSize;i++{
		counter += <-c
	}
	fmt.Printf("I recieved %d ones from %d nodes.\n",counter,channelSize)
	counter += myValue
	var myDecision int
	if float32(counter)/float32(channelSize+1) < 0.5 {
		myDecision = 0
	}else {
		myDecision = 1
	}
	return myDecision
}

func main() {

	fileName := os.Args[1]
	fmt.Println("Launching Program ...")
	fmt.Printf("ConfigFile name: %q\n",fileName)

	adrSlice := initialize(fileName)	// Contain all the ConfigFile addresses

	info := strings.Split(adrSlice[0], ":")	// Splits the first address (our own address) to get the information
	ip := info[0]
	port := info[1]
	id := info[2]
	selfAddress := ip+":"+port
	initiator := false
	parentExists = false
    if len(info)>3 && info[3] == "*" {
		initiator = true
	}


	// In this we randomly set the value of the node that will take part in the
	// communication and will affect the final decision of the network
	rand.Seed(time.Now().UnixNano())
	myValue = rand.Intn(2) // Generates a random binary value
	myStringValue = strconv.Itoa(myValue)



	// Starts listening for connections
	listenerStatus = 0 // Listener Status {0: inactive, 1: active}
	go listener(port, &listenerStatus)
	for listenerStatus != 1{time.Sleep(time.Second)} // Waits until program starts listening


    fmt.Printf("Everything seams ready. Starting Echo...\n\n")
	time.Sleep(3*time.Second)

	fmt.Printf("My value is %q.\n", myStringValue)
	fmt.Printf("My ID is %q.\n",id)


	//connSlice := []net.Conn{} // Stores the active connections

	if initiator==true {
        fmt.Printf("I am the initiatior.\n")

		// Establish all connections and stores them in the Slice var "connSlice"
		for i:=1;i<len(adrSlice);i++ {
			conn, err := establishConnection(adrSlice[i])
			if err!=nil{
				fmt.Printf("An error occurred while trying to connect %q !.\n",adrSlice[i])
				continue
			}
			connSlice = append(connSlice,conn)
		}

		channelSize := len(connSlice)
		c := make(chan int, channelSize)
		/* For each neighbour we have to call a go routine that will be in charge of asking to the corresponding neighbour for their value.*/
        for i:=0;i<len(connSlice);i++ {		// Sends message to all neighbours
        	fmt.Printf("Asking node %q ...\n",adrSlice[i+1])
        	go askNeighbour(selfAddress+":"+id, connSlice[i], c)
		}

		finalDecision := decision(c, channelSize)
		fmt.Printf("Final Decission %d\n", finalDecision)

        // Finishing Program
        broadcast("stop\n", connSlice)
		finishProgram()

    }else{
        for parentExists != true {time.Sleep(time.Second)}	// Waits until it receives a first message (then this client is its parent)
		fmt.Printf("My parent ID is %q\n",parentId)

		// Establish all connections and stores them in the Slice variable "connSlice"
		for i:=1;i<len(adrSlice);i++ {
			if adrSlice[i] != parentAddress {
				conn, err := establishConnection(adrSlice[i])
				if err != nil {
					fmt.Printf("An error occurred while trying to connect %q !.\n", adrSlice[i])
					continue
				}
				connSlice = append(connSlice, conn)
			}
		}

		channelSize := len(connSlice)
		c := make(chan int, channelSize)
		// asks to all neighbours with active connections (that excludes the parent)
        for i:=0;i<len(connSlice);i++ {
			go askNeighbour(selfAddress+":"+id, connSlice[i],c)
        }

        myDecision := decision(c,channelSize)
		fmt.Printf("Sending %d to parent\n",myDecision)
		fmt.Fprintf(parentConn, strconv.Itoa(myDecision)+"\n")
    }
	for listenerStatus != 0{time.Sleep(time.Second)} // Waits until program stops listening
}
