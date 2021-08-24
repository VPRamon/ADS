package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)
type message2main struct {
	round int
	id int
	subtree int
	address string
}

var networkSize int

var myId int
var parentId int

var myAdr string
var parentAdr string

var adrSlice []string
var connSlice []net.Conn 	// Stores the active connections

var listenerStatus int		// 0:passive  1:active
var mutex = &sync.Mutex{}// Synchronize go routines
var c = make(chan message2main)

var initiator bool

func initialize(fileName string) string{
	// Opens the ConfigFile to read
	file, err := os.Open(fileName)
	if err != nil {	log.Fatal(err) }
	defer file.Close()

	var myInfo string
	var myPort string
	scanner := bufio.NewScanner(file)
	scanner.Scan()		// Avoid connecting to ourselves
	myInfo = scanner.Text()
	myPort = strings.Split(myInfo, ":")[1]

	// Starts listening for connections
	listenerStatus = 0 // Listener Status {0: inactive, 1: active}
	go listener(myPort, &listenerStatus)
	for listenerStatus != 1{time.Sleep(time.Second)} // Waits until program starts listening

	for scanner.Scan(){
		adr := scanner.Text()
		conn, _ := establishConnection(adr)

		adrSlice = append(adrSlice,adr)
		connSlice = append(connSlice,conn)

	} // Scans the file and stores the addresses and connections in adrSlice and connSlice
	return myInfo
}

func finishProgram(){
	// Waits 3 seconds before finishing the program
	//we use this function to avoid
	//time.Sleep(3 * time.Second)
	os.Exit(0)
}

func broadcast(message string) {
	for i:=0;i<len(adrSlice);i++ {		// Sends message to all neighbours
		fmt.Fprintf(connSlice[i], message)
	}
}

func multicast(message string, parentAddress string){
	/* Sends a message to all neighbours except parent (or chosen address) */
	for i:=0;i<len(adrSlice);i++{
		if adrSlice[i] != parentAddress {
			fmt.Fprintf(connSlice[i], message)
		}
	}
}

func unicast(message string, address string){
	/* Sends message to a concrete neighbour address */
	for i:=0;i<len(adrSlice);i++{
		if adrSlice[i] == address {
			fmt.Fprintf(connSlice[i], message)
		}
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

	return conn, nil
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
				broadcast(message)
				fmt.Print("\nClient has decided to Stop communication.\nExiting program ...\n")
			}
			mutex.Unlock()
			finishProgram()

			/* If the message is not stop then it must be the id of the parent. */
		}else {
			message = message[:len(message)-1]
			info := strings.Split(message, ",")
			round, _ := strconv.Atoi(info[0])
			id, _ := strconv.Atoi(info[1])
			subtree, _ := strconv.Atoi(info[2])
			adr := info[3]
			c <- message2main{round, id, subtree, adr}
		}
	}
}

func waitForResponse(actualRound int) (int, int){

	awo := 0
	// Check if I have a parent
	if parentId > 0 {
		//If I do not
		awo = 1			// Avoid Waiting to Ourselves
	}

	//waits for responses
	n:=0	// number of nodes reached

	// We wait for our neighbours to respond
	for i := 0; i < len(connSlice)-awo; i++ {
		p := <-c
		senderRound := p.round
		senderId := p.id
		subtree := p.subtree
		senderAdr := p.address

		//fmt.Printf("received message: %d  With %d subtrees...\n", senderId, subtree)

		//if our id is bigger than the one received , we ignore it.
		if senderRound <= actualRound && ( (parentId > 0 && senderId < parentId) || (parentId == -1 && senderId < myId) ) {
			//fmt.Printf("Ignoring: %d ...\n", senderId)
			i -= 1
			continue

		//if our id is smaller, we change from the wave we were and switch to this new received.
		} else if senderRound >= actualRound && ( (initiator==false && parentId < senderId) || (parentId == -1 && myId < senderId)){
			parentId = senderId
			actualRound = senderRound
			parentAdr = senderAdr
			fmt.Printf("My new parent is: %d with adress %q\n", parentId, parentAdr)
			return 0, 1

		}else if senderRound == actualRound && ( (initiator==false && senderId == parentId) || (parentId == -1 && senderId == myId) ) {
			/* If the ID received is equal to our own ID, then we count the size of the nodes reached
			   and wait for all responses to determine if we are the leader */
			n += subtree
		}
	}
	return n, 0
}

func newRound(actualRound int) {
	/* Generates a random value	from 1 to N (size of the network)	*/
	for i:=0;i<100;i++ { rand.Seed(time.Now().UnixNano()) }
	myId = 1 + rand.Intn(networkSize)
	fmt.Printf("My ID is %d.\n",myId)

	n := 0		// number of nodes reached
	err := 0	// 0:{everything fine}	1:{new parent, need to repeat message broadcast}
	if initiator == true {
		msg := strconv.Itoa(actualRound) + "," + strconv.Itoa(myId) + ",0," + myAdr //message{round, myId, subtree_size?, myAdr}
		broadcast(msg + "\n")                                                       // Sends its Id to all neighbours
		n, err = waitForResponse(actualRound)

	} else {
		/* If non initiator */
		p := <-c // Waits for 1st message
		if p.round >= actualRound && p.id > parentId {
			parentId = p.id
			parentAdr = p.address
			fmt.Printf("My Parent is %d with adress %q\n", parentId, parentAdr)
			msg := strconv.Itoa(actualRound) + "," + strconv.Itoa(parentId) + ",0," + myAdr
			multicast(msg + "\n",parentAdr) // Sends its Id to all neighbours
			n, err = waitForResponse(actualRound)
		}
	}

	/* 2 possibilities: all neighbours respond correctly, or there exists a neighbour with greater Id		*/
	/* If there is a neighbour with greater Id, we switch to this wave and multicast that we made a change 	*/
	for err != 0 {
		fmt.Printf("I changed my parent, repeting broadcast\n")
		msg := strconv.Itoa(actualRound) + "," + strconv.Itoa(parentId) + ",0," + myAdr
		multicast(msg+"\n", parentAdr) // Sends its Id to all neighbours
		n, err = waitForResponse(actualRound)
	}

	n = n+1	// We add ourselves as a reached node
	/* If we have a parent, we send our decision to it, with the total number of nodes reached	*/
	if parentId > 0 {
		msg := strconv.Itoa(actualRound) + "," + strconv.Itoa(parentId) + "," + strconv.Itoa(n) + "," + myAdr
		unicast(msg+"\n", parentAdr)
	}

	/* If at some point, we reached all the nodes, it means that we are the leader	*/
	if n == networkSize{
		fmt.Printf("I AM THE LEADER\n")
		broadcast("stop\n")
		finishProgram()
	}

}

func main() {

	fmt.Println("Launching Program ...")

	networkSize, _ = strconv.Atoi(os.Args[1])
	fmt.Printf("Number of Nodes in the Network : %d\n",networkSize)

	/* First Initialize the listener + all connections */
	fileName := os.Args[2]
	fmt.Printf("Reading ConfigFile : %q ...\n",fileName)
	myInfo := initialize(fileName)


	info := strings.Split(myInfo, ":")
	myIp := info[0]
	myPort := info[1]
	myAdr = myIp+":"+myPort
	initiator = false
	if len(info)>2 && info[2] == "*" {
		fmt.Printf("I am an initiatior.\n")
		initiator = true
	}

	fmt.Printf("Everything seams ready. Starting Anonymous Election echo...\n\n")
	time.Sleep(time.Second)

	round := 0
	for{
		parentId = -1
		fmt.Printf("\nTrying round %d...\n",round)
		newRound(round)
		time.Sleep(time.Second)
		round += 1
	}

}
