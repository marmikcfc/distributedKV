package main

import (
	"flag"
	"fmt"
	"./endpoint"
	rt "./router"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)
const(
	address     = "localhost:50051"
)
var (
	tcpAddr     = flag.String("t", ":8081", "The tcp address to bind to for the internal RPC.")
	httpAddr    = flag.String("h", ":8080", "The http address of which to serve the REST API.")
	numberOfNodes = flag.String("n", "1", "Number of nodes to be started with")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	fmt.Println("Started Main")

	r := rt.New()
	e := endpoint.New()

	rpc.Register(r)
	e.RegisterInternalRPC()
	go e.Listen(*httpAddr)

	l, err := net.Listen("tcp", *tcpAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() {
		for {
			conn, _ := l.Accept()
			go rpc.ServeConn(conn)
		}
	}()

	var ok bool
	err = r.AddStore(address, &ok)
	fmt.Println("Added Store")
	if err != nil {
		log.Fatal("AddStore error:", err)
	}
	  i, err := strconv.Atoi(*numberOfNodes)
	err = e.AddRouter(*tcpAddr,i)
	if err != nil {
		log.Fatal("AddRouter error:", err)
	}

	quit := make(chan int)
	<-quit // Wait to be told to exit.
}
