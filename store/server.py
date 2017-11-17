import grpc
from concurrent import futures
import time

# import the generated classes
import kv_pb2
import kv_pb2_grpc

# import the original kv_store.py
import kv_store

kvs = kv_store.KVStore()

# create a class to define the server functions, derived from
# calculator_pb2_grpc.CalculatorServicer
class KVServicer(kv_pb2_grpc.KVServicer):

    # kv_store.square_root is exposed here
    # the request and response are of the data type
    # calculator_pb2.Number
    def Get(self, request, context):
        response = kv_pb2.GetResponse()
        response.value = kvs.get(request.key)
        return response

    # kv_store.square_root is exposed here
    # the request and response are of the data type
    # calculator_pb2.Number
    def Delete(self, request, context):
        response = kv_pb2.GetResponse()
        response.value = kvs.delete(request.key)
        return response

    # kv_store.square_root is exposed here
    # the request and response are of the data type
    # calculator_pb2.Number
    def Put(self, request, context):
        response = kv_pb2.PutResponse()
        print "\n"
        print request.key
        print "\n"
        response.ok = kvs.put(request.key,request.value)
        return response


# create a gRPC server
server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

# use the generated function `add_CalculatorServicer_to_server`
# to add the defined class to the server
kv_pb2_grpc.add_KVServicer_to_server(
        KVServicer(), server)

# listen on port 50051
print('Starting server. Listening on port 50051.')
server.add_insecure_port('[::]:50051')
server.start()

# since server.start() will not block,
# a sleep-loop is added to keep alive
try:
    while True:
        time.sleep(86400)
except KeyboardInterrupt:
    server.stop(0)