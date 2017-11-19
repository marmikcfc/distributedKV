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



/*var keyValClient=pb.NewKVClient(*grpc.ClientConn)

type keyValClient struct {
	cc *grpc.ClientConn
}*/

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

	println ("")
	println ("In GET")
	c, err := r.getClientForKey(key)


	println ("")
	println (key)
	
	res, er := c.Get(context.Background(), &pb.GetRequest{Key: key}, grpc.FailFast(true))
	if er != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", res.Value)

	//*resp = *res
	return res.Value, nil

/*	item, err := c.Get(key)
	if err != nil {
		return err
	}
	*resp = *item
	return nil*/
}

func (r *Router) Put(item *StoreItem) (bool, error) {
	println ("in router put")
	c, err := r.getClientForKey(item.Key)
	println (c)
	if err != nil {
		println(err)
		return false,err
	}



	response, er := c.Put(context.Background(), &pb.PutRequest{Key: item.Key, Value:item.Value}, grpc.FailFast(true))
	if er != nil {
		log.Fatalf("could not greet: %v", er)
	}

	print ("Server put done")
	fmt.Printf("Greeting: %s", response.Ok)

/*
		res, er := c.Put(context.Background(), &pb.PutRequest{Key: "1233", Value:"absdfc"}, grpc.FailFast(true))

			if er != nil {
		log.Fatalf("could not greet: %v", er)
	}
*/
/*	for _, c := range cs {


		res, er := c.Put(context.Background(), &pb.PutRequest{Key: "1233", Value:"absdfc"}, grpc.FailFast(true))
	if er != nil {
		log.Fatalf("could not greet: %v", er)
	}
*/	

/*		*added, err = c.Put(item)
		if err != nil {
			return err
		}
*/	
	return true,nil
}

/*func (r *Router) Delete(key string, ack *bool) error {
	cs, err := r.getClientsForKey(key)
	if err != nil {
		return err
	}
	for _, c := range cs {
		_, err := c.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}*/

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
/*	c, err := store.NewClient(addr, 500*time.Millisecond)
	if err != nil {
		return err
	}
*/ 
	r.clients[addr] = c
	println (r.clients[addr])
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
	println("In get getClientForKey")
	println(s)
	if err != nil {
		return nil, err
	}
	c, _ := r.clients[s]
	println (c)
	return c, nil
}
