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
	"encoding/json"
	"fmt"
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
	 group, err := parseURI(req.URL.RequestURI())
	if err != nil {
		log.Printf("HTTP Action returned error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URI: " + err.Error() + "\n"))
		return
	}

	var resp string

	switch group{
	case "set":

	type KeyVal struct{
			Key string
			Value string
		}

	type rd struct{
		Keys_added int32
		Keys_failed []string
	}

		var resp_data rd
		var kv []KeyVal
		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			break
		}

		er :=  json.Unmarshal(buf, &kv)
		if er != nil {
			break
		}

		var add int32
		var fail []string
		add =0 
    	for l := range kv {
    	added, err := e.Put(kv[l].Key, kv[l].Value)
    	fmt.Printf("Key = %v, Val= %v", kv[l].Key, kv[l].Value)
		if added {
			add++
			w.WriteHeader(http.StatusCreated)
		}
		if err != nil {
			fail = append(fail,kv[l].Key)
		}
    	}

    	resp_data.Keys_added = add
    	resp_data.Keys_failed = fail
    	respData, err := json.Marshal(&resp_data)
    	arr := []byte(respData)
    	w.Write(arr)

	case "fetch":
		switch req.Method {
			case "GET":
				buf, err := ioutil.ReadAll(req.Body)
				data := string(buf[:])
				data = strings.Replace(data, "[", "", -1)
				data = strings.Replace(data, "]", "", -1)
				s := strings.Split(data, ",")
				type KeyVal struct{
			Key string
			Value string
		}

		var kv []KeyVal
				for i := range s {
				var kvItem KeyVal
				resp, err = e.Get(strings.Replace(s[i], "\"", "", -1))
				kvItem.Key = s[i]
				kvItem.Value = resp
				kv = append(kv,kvItem)
				if err != nil {
					break
				}

				}
				respData, err := json.Marshal(&kv)
				arr := []byte(respData)
				w.Write(arr)

			case "POST":
				buf, err := ioutil.ReadAll(req.Body)
				data := string(buf[:])
				data = strings.Replace(data, "[", "", -1)
				data = strings.Replace(data, "]", "", -1)
				s := strings.Split(data, ",")
				type KeyVal struct{
			Key string
			Value string
		}

		var kv []KeyVal
				for i := range s {
				var kvItem KeyVal
				resp, err = e.Get(strings.Replace(s[i], "\"", "", -1))
				kvItem.Key = s[i]
				kvItem.Value = resp
				kv = append(kv,kvItem)
				if err != nil {
					break
				}

				}
				respData, err := json.Marshal(&kv)
				arr := []byte(respData)
				w.Write(arr)		}


	case "query":
		switch req.Method {
			case "GET":
				buf, err := ioutil.ReadAll(req.Body)
				data := string(buf[:])
				data = strings.Replace(data, "[", "", -1)
				data = strings.Replace(data, "]", "", -1)
				s := strings.Split(data, ",")
				type KeyVal struct{
			Key string
			Value bool
		}

		var kv []KeyVal
				for i := range s {
				var kvItem KeyVal
				resp, err = e.Get(strings.Replace(s[i], "\"", "", -1))
				kvItem.Key = s[i]
				println (resp)
				if (resp == "Key does not exist"){
					kvItem.Value = false
				} else{
					kvItem.Value = true
				}
				kv = append(kv,kvItem)
				if err != nil {
					break
				}

				}
				respData, err := json.Marshal(&kv)
				arr := []byte(respData)
				w.Write(arr)

			case "POST":
				buf, err := ioutil.ReadAll(req.Body)
				data := string(buf[:])
				data = strings.Replace(data, "[", "", -1)
				data = strings.Replace(data, "]", "", -1)
				s := strings.Split(data, ",")
				type KeyVal struct{
			Key string
			Value bool
		}

		var kv []KeyVal
				for i := range s {
				var kvItem KeyVal
				resp, err = e.Get(strings.Replace(s[i], "\"", "", -1))
				kvItem.Key = s[i]
				println (resp)
				if (resp == "Key does not exist"){
					kvItem.Value = false
				} else{
					kvItem.Value = true
				}
				kv = append(kv,kvItem)
				if err != nil {
					break
				}

				}
				respData, err := json.Marshal(&kv)
				arr := []byte(respData)
				w.Write(arr)		}

	default:
		w.Write([]byte("404 - Not a valid endpoint!"))

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

	if err != nil {
		return false, err
	}

	added, err := r.Route.Put(&rt.StoreItem{Key: key, Value: value})
	return added, err
}


func (e *Endpoint) AddRouter(addr string, number_of_nodes int) error {
	var ok bool
	return e.internal.AddRouter(addr, &ok,number_of_nodes)
}

func (e *EndpointInternal) AddRouter(addr string, ok *bool, number_of_nodes int) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	c, err := rt.NewClient(addr, 500*time.Millisecond, number_of_nodes)
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

func parseURI(uri string) (string, error) {
	s := strings.Split(uri, "/")
	if len(s) != 2 {
		return "", errors.New("URI " + uri + " does not match the format")
	}
	return s[1], nil
}
