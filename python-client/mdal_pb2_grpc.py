# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import mdal_pb2 as mdal__pb2


class MDALStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.DataQuery = channel.unary_stream(
        '/mdalgrpc.MDAL/DataQuery',
        request_serializer=mdal__pb2.DataQueryRequest.SerializeToString,
        response_deserializer=mdal__pb2.DataQueryResponse.FromString,
        )


class MDALServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def DataQuery(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_MDALServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'DataQuery': grpc.unary_stream_rpc_method_handler(
          servicer.DataQuery,
          request_deserializer=mdal__pb2.DataQueryRequest.FromString,
          response_serializer=mdal__pb2.DataQueryResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'mdalgrpc.MDAL', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))