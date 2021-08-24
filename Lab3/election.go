package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)
type message2main struct {
		Id int
		ip string
		port string
}

var myId int
var parentId int

var myAdr string
var parentAdr string

var adrSlice []string
var connSlice []net.Conn 	// Stores the active connections

var listenerStatus int		// 0:passive  1:active
var mutex = &sync.Mutex{}// Synchronize go routines
var c = make(chan message2main)


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
	time.Sleep(3 * time.Second)
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
	//fmt.Printf("%q is now ready. Establishing connection\n", adr)
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
		//fmt.Print(message,conn)
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
			info := strings.Split(message, ":")
			//fmt.Println("Scanning")
			ip := info[0]
			port := info[1]
			idStr := info[2]
			idStr = idStr[:len(idStr)-1];

			idInt, _ := strconv.Atoi(idStr)
			//fmt.Println("Scanning")
			//fmt.Println(message2main{idInt,ip,port})
			//fmt.Println("Closing the channel")
			c <- message2main{idInt, ip, port}
			//fmt.Println("Closing the channel")
			//close(c)
		}
	}
}


func main() {
	//c := make(chan message2main)
	fileName := os.Args[1]
	fmt.Println("Launching Program ...")
	fmt.Printf("ConfigFile name: %q\n",fileName)

	/* Firs Initialize the listener + all connections */
	myInfo := initialize(fileName)	// The function also returns our own information

	info := strings.Split(myInfo, ":")	// Splits the first address (our own address) to get the information
	myIp := info[0]
	myPort := info[1]
	myId, _ := strconv.Atoi(info[2])
	myAdr = myIp+":"+myPort
	initiator := false
	parentId = -1
	if len(info)>3 && info[3] == "*" {
		initiator = true
		parentId = myId
	}


	fmt.Printf("Everything seams ready. Starting Election echo...\n\n")
	time.Sleep(3*time.Second)

	fmt.Printf("My ID is %d.\n",myId)

	if initiator==true {
		fmt.Printf("I am an initiatior.\n")
		broadcast(myAdr+":"+strconv.Itoa(myId)+"\n")	// Sends its Id to all neighbours
	}

	consensus := 0;
	//fmt.Printf("reciving from one of the listening functions\n")
	var p message2main;
	for{
		if consensus == len(connSlice){
			fmt.Printf("I am the leader\n")
			broadcast("stop\n")
			finishProgram()
		}
		p = <- c
		//CODE Here
		// AQUI TINDRIEM QUE FER TOTES LES OPCIONS,SI ES MES GRAN, MES PETIT O IGUAL, I PRENDRE DECISIONS.
		idInt := p.Id;
		ip := p.ip;
		port := p.port;

		//fmt.Printf("Info received: %d  %q  %q\n",idInt, ip, port);


		if parentId > idInt{
			/* If the received ID is smaller than our parent ID, then we ignore it */
			fmt.Printf("Ignoring: %d ...\n",idInt);
			continue
		}else if parentId < idInt{
			/* If the received ID is greater than our parent ID, we switch to this wave */
			parentId = idInt;
			parentAdr = ip+":"+port;
			fmt.Printf("My new parent is: %d\n",parentId);
			multicast(myAdr+":"+strconv.Itoa(idInt)+"\n", parentAdr);

			l := 0;
			if(myId != parentId){
				l = 1;
			}
			for i:=0;i<len(connSlice)-l;i++ {
				p = <- c;
				idInt := p.Id;
				//ip := p_aux.ip;
				//port := p_aux.port;
				if(idInt > parentId){
					parentId = idInt;
				}
			}
			// SI TOTS ELS NEIGHBOURS RESPONEN AMB EL MATEIX ID, ALESHORES LI RESPONEM AL PARE AMB AQUEST ID
			unicast(myAdr+":"+strconv.Itoa(idInt)+"\n", parentAdr);
		}else if myId == parentId{
			/* If the ID received is equal to our own ID, then we wait for all responses
			   to determine if we are the leader */
			consensus += 1;
			// wait for consensus
			// Write code here ...
		}else{
			/* Else, If the  received ID is equal to our parents ID, but we are not an initiator,
			   then we respond to our parent wih his own ID (similar to an ACK)   		 */
			/* Notice that to indicate that we are initiators we pointed ourselves to be our own parent,
			   thus we would fall in the previous condition 							 */
			unicast(myAdr+":"+strconv.Itoa(idInt)+"\n", parentAdr);
		}

	}
	//fmt.Printf("reciving from one of the listening functions\n")
	//fmt.Println(p)


}
