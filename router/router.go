package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "./genFile"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewKVClient(conn)

	// // Contact the server and print out its response.
	
	r, err := c.Put(context.Background(), &pb.PutRequest{Key: "12433", Value:"absdfc"}, grpc.FailFast(true))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Ok)


	r1, err1 := c.Put(context.Background(), &pb.PutRequest{Key: "1233", Value:"abc"}, grpc.FailFast(true))
	if err1 != nil {
		log.Fatalf("could not greet: %v", err1)
	}
	log.Printf("Greeting: %s", r1.Ok)


		r2, err2 := c.Get(context.Background(), &pb.GetRequest{Key: "1233"}, grpc.FailFast(true))
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Printf("Greeting: %s", r2.Value)

/*

		r3, err3 := c.Delete(context.Background(), &pb.DeleteRequest{Key: "1233"}, grpc.FailFast(true))
	if err3 != nil {
		log.Fatalf("could not greet: %v", err3)
	}
	log.Printf("Greeting: %s", r3.Ok)*/
}