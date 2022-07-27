/**
 * @fileoverview gRPC-Web generated client stub for treediagram.pb
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as service_pb from './service_pb';


export class TreeDiagramServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'binary';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorRenderEvent = new grpcWeb.MethodDescriptor(
    '/treediagram.pb.TreeDiagramService/RenderEvent',
    grpcWeb.MethodType.UNARY,
    service_pb.RenderEventRequest,
    service_pb.RenderEventResponse,
    (request: service_pb.RenderEventRequest) => {
      return request.serializeBinary();
    },
    service_pb.RenderEventResponse.deserializeBinary
  );

  renderEvent(
    request: service_pb.RenderEventRequest,
    metadata: grpcWeb.Metadata | null): Promise<service_pb.RenderEventResponse>;

  renderEvent(
    request: service_pb.RenderEventRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderEventResponse) => void): grpcWeb.ClientReadableStream<service_pb.RenderEventResponse>;

  renderEvent(
    request: service_pb.RenderEventRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: service_pb.RenderEventResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/treediagram.pb.TreeDiagramService/RenderEvent',
        request,
        metadata || {},
        this.methodDescriptorRenderEvent,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderEvent',
    request,
    metadata || {},
    this.methodDescriptorRenderEvent);
  }

  methodDescriptorRenderPlace = new grpcWeb.MethodDescriptor(
    '/treediagram.pb.TreeDiagramService/RenderPlace',
    grpcWeb.MethodType.UNARY,
    service_pb.RenderPlaceRequest,
    service_pb.RenderPlaceResponse,
    (request: service_pb.RenderPlaceRequest) => {
      return request.serializeBinary();
    },
    service_pb.RenderPlaceResponse.deserializeBinary
  );

  renderPlace(
    request: service_pb.RenderPlaceRequest,
    metadata: grpcWeb.Metadata | null): Promise<service_pb.RenderPlaceResponse>;

  renderPlace(
    request: service_pb.RenderPlaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderPlaceResponse) => void): grpcWeb.ClientReadableStream<service_pb.RenderPlaceResponse>;

  renderPlace(
    request: service_pb.RenderPlaceRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: service_pb.RenderPlaceResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/treediagram.pb.TreeDiagramService/RenderPlace',
        request,
        metadata || {},
        this.methodDescriptorRenderPlace,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderPlace',
    request,
    metadata || {},
    this.methodDescriptorRenderPlace);
  }

  methodDescriptorRenderActors = new grpcWeb.MethodDescriptor(
    '/treediagram.pb.TreeDiagramService/RenderActors',
    grpcWeb.MethodType.UNARY,
    service_pb.RenderActorsRequest,
    service_pb.RenderActorsResponse,
    (request: service_pb.RenderActorsRequest) => {
      return request.serializeBinary();
    },
    service_pb.RenderActorsResponse.deserializeBinary
  );

  renderActors(
    request: service_pb.RenderActorsRequest,
    metadata: grpcWeb.Metadata | null): Promise<service_pb.RenderActorsResponse>;

  renderActors(
    request: service_pb.RenderActorsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.RenderActorsResponse) => void): grpcWeb.ClientReadableStream<service_pb.RenderActorsResponse>;

  renderActors(
    request: service_pb.RenderActorsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: service_pb.RenderActorsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/treediagram.pb.TreeDiagramService/RenderActors',
        request,
        metadata || {},
        this.methodDescriptorRenderActors,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderActors',
    request,
    metadata || {},
    this.methodDescriptorRenderActors);
  }

  methodDescriptorQueryEvents = new grpcWeb.MethodDescriptor(
    '/treediagram.pb.TreeDiagramService/QueryEvents',
    grpcWeb.MethodType.UNARY,
    service_pb.QueryEventsRequest,
    service_pb.QueryEventsResponse,
    (request: service_pb.QueryEventsRequest) => {
      return request.serializeBinary();
    },
    service_pb.QueryEventsResponse.deserializeBinary
  );

  queryEvents(
    request: service_pb.QueryEventsRequest,
    metadata: grpcWeb.Metadata | null): Promise<service_pb.QueryEventsResponse>;

  queryEvents(
    request: service_pb.QueryEventsRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: service_pb.QueryEventsResponse) => void): grpcWeb.ClientReadableStream<service_pb.QueryEventsResponse>;

  queryEvents(
    request: service_pb.QueryEventsRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: service_pb.QueryEventsResponse) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/treediagram.pb.TreeDiagramService/QueryEvents',
        request,
        metadata || {},
        this.methodDescriptorQueryEvents,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/treediagram.pb.TreeDiagramService/QueryEvents',
    request,
    metadata || {},
    this.methodDescriptorQueryEvents);
  }

}

