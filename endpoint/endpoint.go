package endpoint

import (
	"errors"
	rt "../router"
	"stathat.com/c/consistent"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

var _ = log.Printf

type (
	Endpoint struct {
		internal *EndpointInternal
	}
	EndpointInternal struct {
		routers  map[string]*rt.Client
		hashRing *consistent.Consistent
		mu       *sync.RWMutex
	}
)

func New() *Endpoint { 
	i := &EndpointInternal{
		routers:  make(map[string]*rt.Client),
		hashRing: consistent.New(),
		mu:       &sync.RWMutex{},
	}
	return &Endpoint{internal: i}
}

func (e *Endpoint) RegisterInternalRPC() {
	rpc.Register(e.internal)
}

func (e *Endpoint) Listen(httpAddr string) {
	http.HandleFunc("/", e.StoreHandler)
	log.Println(http.ListenAndServe(httpAddr, nil))
}

func (e *Endpoint) StoreHandler(w http.ResponseWriter, req *http.Request) {
	namespace, group, id, err := parseURI(req.URL.RequestURI())
	if err != nil {
		log.Printf("HTTP Action returned error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URI: " + err.Error() + "\n"))
		return
	}

	// Temporal hack, will be replaced when store can properly handle gruops
	key := namespace + "/" + group + "/" + id

	var resp string
	switch req.Method {
	case "GET":
		buf, err := ioutil.ReadAll(req.Body)
		data := string(buf[:])
		resp, err = e.Get(data)
		if err != nil {
			break
		}
		arr := []byte(resp)
		w.Write(arr)
	case "PUT":
		println("IN PUT")
		buf, err := ioutil.ReadAll(req.Body)
		data := string(buf[:])
		if err != nil {
			break
		}
			s := strings.Split(data, " ")
    	key, value := s[0], s[1]

		added, err := e.Put(key, value)
		if added {
			w.WriteHeader(http.StatusCreated)
		}
	case "DELETE":
		_, err = e.Delete(key)
	}

	if err != nil {
		log.Printf("HTTP Action returned error: %s", err)
		if err.Error() == "Key not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("Error: " + err.Error() + "\n"))
	}
}
 
func (e *Endpoint) Get(data string) (string, error) {
	r, err := e.internal.getRouterForKey(data)
	println ("GET Router")
	println (r)
	println ("##################################")
	if err != nil {
		return "", err
	}
	item, err := r.Route.Get(data)
	if err != nil {
		return "", err
	}
	return string(item), err
}

func (e *Endpoint) Put(key string, value string) (bool, error) {

    	println (key)
    	println (value)
	r, err := e.internal.getRouterForKey(key)

/*	println ("Router")
	println (r)
	println ("##################################")
*/
	if err != nil {
		return false, err
	}
	
/*	println("")
	println(r.Route)
*/
	added, err := r.Route.Put(&rt.StoreItem{Key: key, Value: value})
/*	println("put done")
	println(added)
*/
	return added, err
}

func (e *Endpoint) Delete(key string) ([]byte, error) {
	r, err := e.internal.getRouterForKey(key)
	if err != nil {
		return nil, err
	}
	_, err = r.Delete(key)
	return nil, err
}

func (e *Endpoint) AddRouter(addr string) error {
	var ok bool
	return e.internal.AddRouter(addr, &ok)
}

func (e *EndpointInternal) AddRouter(addr string, ok *bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	c, err := rt.NewClient(addr, 500*time.Millisecond)
	if err != nil {
		return err
	}
	e.routers[addr] = c
	e.hashRing.Add(addr)
	return nil
}

func (e *EndpointInternal) getRouterForKey(key string) (*rt.Client, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	s, err := e.hashRing.Get(key)
	if err != nil {
		return nil, err
	}
	c, _ := e.routers[s]
	return c, nil
}

func parseURI(uri string) (string, string, string, error) {
	s := strings.Split(uri, "/")
	if len(s) != 4 {
		return "", "", "", errors.New("URI " + uri + " does not match /[namespace]/[key]/[id]")
	}
	return s[0], s[1], s[2], nil
}
