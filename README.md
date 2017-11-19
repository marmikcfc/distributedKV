PUT
    curl -X PUT --data "bar" localhost:8080/default/foo/0

GET
    curl -X GET localhost:8080/default/foo/0

Python
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. kv.proto

protoc --go_out=plugins=grpc:../router/genFile kv.proto