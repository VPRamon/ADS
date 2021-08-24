package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
    //"strconv"
    "math/rand"
)
import "fmt"
import "bufio"

var myValue string
var parentExists bool
var parentConn net.Conn
var parentAddress string
var adrSlice []string
var connSlice []net.Conn

func finishProgram(){
	// Waits 3 seconds before finishing the program.
	//we use this fucntion to avoid
	time.Sleep(3 * time.Second)
	os.Exit(0)
}


func askNeighbour(message string, conn net.Conn, c chan int) {
	// In this function we read the information send by the neighbours,that is why we need
	//to read in the connetion the information sent by the neighbour.If the message sent is
	//a 0(we also add a \n to the parsing as in all the message this character is sent after the body of the message)
	// we then sent the int 0 to the channel, if it is 1, we send 1.
	fmt.Fprintf(conn, message+"\n")
	response, _ := bufio.NewReader(conn).ReadString('\n')
	if response == "0\n"{
		c <- 0
	}else if response == "1\n" {
		c <- 1
	}
}

func newConnexion(conn net.Conn){
	/* 	This function is called every time a new client requests a connection.
		Here we receive  every message from this `particular client   			*/
	// run loop forever (or until "stop" message is received)
	for {
		// will listen for message to process ending in newline (\n)

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Print("Error while reading message\n")
			finishProgram()
			break
		}

			// we use then clause/command of stop to finish connection, therefore in every message sent,
			//we check if the mmessage does indeed has the form of "stop"
		if message == "stop\n"{
			fmt.Print("\nClient has decided to Stop communication.\nExiting program ...\n")
			finishProgram()
		//If the message is not stop then it has a real value to the Algorithm and it could be in
		//the case that is the first message, the id of the parent.For this we will have a bool
		//varibale call partenExists that will be set to True once we  meet our parent
		//if we already have a parent, we will print the id of our parent
		}else{

            if parentExists == false {
				parentAddress = message[:len(message)-1]
                parentConn = conn
				parentExists = true
            }else{
                fmt.Fprintf(conn, myValue+"\n")
            }
		}
	}
}

func listener(address string, status *int) {
	//Starts listening in the specified port
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

		//fmt.Printf("New connection request...\n")
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


	adrSlice = []string{}	// Contain all the ConfigFile addresses
	scanner := bufio.NewScanner(file)
	// Reads all the addresses of the file and stores them in the "adrSlice" Slice
	for scanner.Scan(){
		adr := scanner.Text()
		adrSlice = append(adrSlice,adr)
	} // Scans the file and stores the addresses in adrSlice


	info := strings.Split(adrSlice[0], ":")	// Splits the first address (our own address) to get the port
	ip := info[0]
	port := info[1]
	selfAddress := ip+":"+port
	initiator := false
	parentExists = false
    if len(info)>3{
    	if info[3] == "*" {
			initiator = true
		}
	}


	status := 0		// Listener Status {0: inactive, 1: active}
	go listener(port, &status)
	for status != 1{time.Sleep(time.Second)}	// Waits until program starts listening

    connSlice = []net.Conn{} // Stores the active connections
			//NEEEEW//
			// in this we randomly set the value of the node that will take part in the
			//communication and will affect the final desicion of the network

    val := rand.Intn(2)
    if val == 0{
        myValue = "0"
    }else{
        myValue = "1"
    }
    fmt.Printf("My value is %q.\n",myValue)

    fmt.Printf("Everything seams ready. Starting Echo...\n\n")
	time.Sleep(3*time.Second)
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

			//we define the channel size
		channelSize := len(connSlice)
		c := make(chan int, channelSize)
        for i:=0;i<len(connSlice);i++ {		// Sends message to all neighbours
        	fmt.Printf("Asking node %q ...\n",adrSlice[i+1])
						//NEEEW//
						//for each neighbour we have to call a go routine that will be in
						//charge of asking to the corresponding neighbour
        	go askNeighbour(selfAddress, connSlice[i], c)
		}
		counter := 0
		for i:=0;i<channelSize;i++{
			//fmt.Printf("c = %d\n",<-c)
			counter += <-c
		}
		fmt.Printf("I recieved %d ones from %d nodes.\n",counter,cap(c))
        if counter < cap(c)/2{
            fmt.Printf("Final Decission 0\n")
        }else{
            fmt.Printf("Final decission 1\n")
        }
        broadcast("stop\n", connSlice)
		finishProgram()
    }else{
        for parentExists != true {time.Sleep(time.Second)}
		fmt.Printf("My parent is %q\n",parentAddress)

		// Establish all connections and stores them in the Slice var "connSlice"
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
		// asks to all neighbours except parent
        for i:=0;i<len(connSlice);i++ {
			go askNeighbour(selfAddress, connSlice[i],c)
        }

        counter := 0
		for i:=0;i<channelSize;i++{
			//fmt.Printf("c = %d\n",<-c)
			counter += <-c
		}
		fmt.Printf("Sending %d ones to parent\n",counter)
        if counter < cap(c)/2{
            fmt.Fprintf(parentConn, "0\n")
        }else{
            fmt.Fprintf(parentConn, "1\n")
        }
    }
	for status != 0{time.Sleep(time.Second)}	// Waits until program stops listening
}
