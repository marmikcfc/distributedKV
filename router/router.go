package router

import (
	"stathat.com/c/consistent"
	"log"
	"sync"
	"fmt"
	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"

	pb "./genFile"
//	"reflect"
)

var _ = log.Printf

type (

	StoreItem struct {
		Key   string
		Value string
	}
)


type (
	Router struct {
		clients  map[string]pb.KVClient
		hashRing *consistent.Consistent
		mu       *sync.RWMutex
		Replicas int
	}
)

func New() *Router {

	println ("creating new router")
	r := &Router{
		clients:  make(map[string]pb.KVClient),
		hashRing: consistent.New(),
		mu:       &sync.RWMutex{},
		Replicas: 2,
	}
	return r
}

func (r *Router) Get(key string) (string, error) {
	// c is client

	c, err := r.getClientForKey(key)


	println ("")
	println (key)
	
	res, er := c.Get(context.Background(), &pb.GetRequest{Key: key}, grpc.FailFast(true))
	if er != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	log.Printf("Server's response: %s", res.Value)
	return res.Value, nil

}

func (r *Router) Put(item *StoreItem) (bool, error) {
	c, err := r.getClientForKey(item.Key)
	println (c)
	if err != nil {
		println(err)
		return false,err
	}



	response, er := c.Put(context.Background(), &pb.PutRequest{Key: item.Key, Value:item.Value}, grpc.FailFast(true))
	if er != nil {
		log.Fatalf("could not connect to server: %v", er)
	}

	fmt.Printf("Server's response %s", response.Ok)

	return true,nil
}


func (r *Router) AddStore(addr string, ok *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	c := pb.NewKVClient(conn)
	println(c)
	r.clients[addr] = c
	r.hashRing.Add(addr)
	return nil
}

func (r *Router) getClientsForKey(key string) ([]pb.KVClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names, err := r.hashRing.GetN(key, r.Replicas)

	if err != nil {
		return nil, err
	}
	cs := make([]pb.KVClient, 0, len(names))
	for _, value := range names {
		c, _ := r.clients[value]
		cs = append(cs, c)
	}
	return cs, nil
}

func (r *Router) getClientForKey(key string) (pb.KVClient, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, err := r.hashRing.Get(key)
	if err != nil {
		return nil, err
	}
	c, _ := r.clients[s]
	return c, nil
}
