# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import kv_pb2 as kv__pb2


class KVStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Get = channel.unary_unary(
        '/kv.KV/Get',
        request_serializer=kv__pb2.GetRequest.SerializeToString,
        response_deserializer=kv__pb2.GetResponse.FromString,
        )
    self.Put = channel.unary_unary(
        '/kv.KV/Put',
        request_serializer=kv__pb2.PutRequest.SerializeToString,
        response_deserializer=kv__pb2.PutResponse.FromString,
        )
    self.Delete = channel.unary_unary(
        '/kv.KV/Delete',
        request_serializer=kv__pb2.DeleteRequest.SerializeToString,
        response_deserializer=kv__pb2.DeleteResponse.FromString,
        )


class KVServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Get(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Put(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Delete(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_KVServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Get': grpc.unary_unary_rpc_method_handler(
          servicer.Get,
          request_deserializer=kv__pb2.GetRequest.FromString,
          response_serializer=kv__pb2.GetResponse.SerializeToString,
      ),
      'Put': grpc.unary_unary_rpc_method_handler(
          servicer.Put,
          request_deserializer=kv__pb2.PutRequest.FromString,
          response_serializer=kv__pb2.PutResponse.SerializeToString,
      ),
      'Delete': grpc.unary_unary_rpc_method_handler(
          servicer.Delete,
          request_deserializer=kv__pb2.DeleteRequest.FromString,
          response_serializer=kv__pb2.DeleteResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'kv.KV', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))