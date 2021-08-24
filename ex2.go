package main

import "net"
import "fmt"
import "bufio"

import "os"
import "log"

//netstat -tulpn
//
func listening(addr string, c chan string){
      fmt.Println("listening on port",addr)
      ln, _ := net.Listen("tcp",addr)
    
      for {
        conn, _ := ln.Accept()
        fmt.Print(conn.RemoteAddr())
        // will listen for message to process ending in newline (\n)
        message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil { break }
        // output message received
        fmt.Print("Message Received:", string(message))
        // sample process for string received

        // send new string back to client
        c <- message
      }

}


func main(){

arg := os.Args[1]
file, err := os.Open(arg)
var a []string
var message string
c := make(chan string, 400)
if err != nil {
      log.Fatal(err)
      fmt.Println("error")
  }
  defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {

    //we append all the ips to the a array, that will be passed
    // tot the server function
    a = append(a,string(scanner.Text()))}
    go listening(a[0],c)


  for {
    message =  <- c
    fmt.Println("from the local port I have recived",message)
    fmt.Println("Sending now to the others")
    for i:=1;i < len(a);i++{
      conn, _ := net.Dial("tcp", a[i])
      conn.Write([]byte(message + "\n"))


      }

  }
//a[0] has tthe ip of the server
//a[1:] has all the ips of the networks









}
