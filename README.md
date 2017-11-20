##Distributed Key Value Store

###Running Instructions
First run `make` to install all the go and python requirements.

Next, you can run python server by python `.\store\server.py` and go router by `go run kvstore.go`


### Implementation

Router is implemented by go. By default it starts with 1 node but you can specify the the number of nodes using `-n` flag. Router uses consistent hashing(third party package) and hence data won't be affected by changing the number of nodes.

Communication between client and the router is via HTTP and communication between server and router is via gRPC. 


### Testing

endpoint /set can be tested by :-
    curl -H "Content-Type: application/json" -X PUT -d "[{\"key\":\"1\",\"value\": \"pikachu\" }, {\"key\":\"2\",\"value\":\"Raichu\"}]" http://localhost:8080/set

endpoint /fetch an be tested by 
    curl -X GET --data "[\"1\",\"2\"]" localhost:8080/fetch

endpoint /query can be tested by
	curl -X GET --data "[\"1\",\"2\"]" localhost:8080/query
