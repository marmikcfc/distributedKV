package main

import (
	rt "./router"
)



func main() {
	

	r :=rt.New()
	var ok bool
	r.AddStore(address,&ok)
	println("added Store")
	r.Put(&rt.StoreItem{Key: "1234", Value: "data"})	
	r.Get("123")	
}
