import * as grpcWeb from 'grpc-web';

import * as service_pb from './service_pb';


export class TreeDiagramServiceClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  renderEvent(
    request: service_pb.RenderEventRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderEventResponse) => void
  ): grpcWeb.ClientReadableStream<service_pb.RenderEventResponse>;

  renderPlace(
    request: service_pb.RenderPlaceRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderPlaceResponse) => void
  ): grpcWeb.ClientReadableStream<service_pb.RenderPlaceResponse>;

  renderActors(
    request: service_pb.RenderActorsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderActorsResponse) => void
  ): grpcWeb.ClientReadableStream<service_pb.RenderActorsResponse>;

  queryEvents(
    request: service_pb.QueryEventsRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.QueryEventsResponse) => void
  ): grpcWeb.ClientReadableStream<service_pb.QueryEventsResponse>;

}

export class TreeDiagramServicePromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  renderEvent(
    request: service_pb.RenderEventRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_pb.RenderEventResponse>;

  renderPlace(
    request: service_pb.RenderPlaceRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_pb.RenderPlaceResponse>;

  renderActors(
    request: service_pb.RenderActorsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_pb.RenderActorsResponse>;

  queryEvents(
    request: service_pb.QueryEventsRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<service_pb.QueryEventsResponse>;

}

