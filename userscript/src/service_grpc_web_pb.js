/**
 * @fileoverview gRPC-Web generated client stub for treediagram.pb
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
const proto = {};
proto.treediagram = {};
proto.treediagram.pb = require('./service_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.treediagram.pb.TreeDiagramServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.treediagram.pb.TreeDiagramServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.treediagram.pb.RenderEventRequest,
 *   !proto.treediagram.pb.RenderEventResponse>}
 */
const methodDescriptor_TreeDiagramService_RenderEvent = new grpc.web.MethodDescriptor(
  '/treediagram.pb.TreeDiagramService/RenderEvent',
  grpc.web.MethodType.UNARY,
  proto.treediagram.pb.RenderEventRequest,
  proto.treediagram.pb.RenderEventResponse,
  /**
   * @param {!proto.treediagram.pb.RenderEventRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderEventResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.treediagram.pb.RenderEventRequest,
 *   !proto.treediagram.pb.RenderEventResponse>}
 */
const methodInfo_TreeDiagramService_RenderEvent = new grpc.web.AbstractClientBase.MethodInfo(
  proto.treediagram.pb.RenderEventResponse,
  /**
   * @param {!proto.treediagram.pb.RenderEventRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderEventResponse.deserializeBinary
);


/**
 * @param {!proto.treediagram.pb.RenderEventRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.treediagram.pb.RenderEventResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.treediagram.pb.RenderEventResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.treediagram.pb.TreeDiagramServiceClient.prototype.renderEvent =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderEvent',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderEvent,
      callback);
};


/**
 * @param {!proto.treediagram.pb.RenderEventRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.treediagram.pb.RenderEventResponse>}
 *     Promise that resolves to the response
 */
proto.treediagram.pb.TreeDiagramServicePromiseClient.prototype.renderEvent =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderEvent',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderEvent);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.treediagram.pb.RenderPlaceRequest,
 *   !proto.treediagram.pb.RenderPlaceResponse>}
 */
const methodDescriptor_TreeDiagramService_RenderPlace = new grpc.web.MethodDescriptor(
  '/treediagram.pb.TreeDiagramService/RenderPlace',
  grpc.web.MethodType.UNARY,
  proto.treediagram.pb.RenderPlaceRequest,
  proto.treediagram.pb.RenderPlaceResponse,
  /**
   * @param {!proto.treediagram.pb.RenderPlaceRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderPlaceResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.treediagram.pb.RenderPlaceRequest,
 *   !proto.treediagram.pb.RenderPlaceResponse>}
 */
const methodInfo_TreeDiagramService_RenderPlace = new grpc.web.AbstractClientBase.MethodInfo(
  proto.treediagram.pb.RenderPlaceResponse,
  /**
   * @param {!proto.treediagram.pb.RenderPlaceRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderPlaceResponse.deserializeBinary
);


/**
 * @param {!proto.treediagram.pb.RenderPlaceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.treediagram.pb.RenderPlaceResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.treediagram.pb.RenderPlaceResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.treediagram.pb.TreeDiagramServiceClient.prototype.renderPlace =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderPlace',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderPlace,
      callback);
};


/**
 * @param {!proto.treediagram.pb.RenderPlaceRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.treediagram.pb.RenderPlaceResponse>}
 *     Promise that resolves to the response
 */
proto.treediagram.pb.TreeDiagramServicePromiseClient.prototype.renderPlace =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderPlace',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderPlace);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.treediagram.pb.RenderActorsRequest,
 *   !proto.treediagram.pb.RenderActorsResponse>}
 */
const methodDescriptor_TreeDiagramService_RenderActors = new grpc.web.MethodDescriptor(
  '/treediagram.pb.TreeDiagramService/RenderActors',
  grpc.web.MethodType.UNARY,
  proto.treediagram.pb.RenderActorsRequest,
  proto.treediagram.pb.RenderActorsResponse,
  /**
   * @param {!proto.treediagram.pb.RenderActorsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderActorsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.treediagram.pb.RenderActorsRequest,
 *   !proto.treediagram.pb.RenderActorsResponse>}
 */
const methodInfo_TreeDiagramService_RenderActors = new grpc.web.AbstractClientBase.MethodInfo(
  proto.treediagram.pb.RenderActorsResponse,
  /**
   * @param {!proto.treediagram.pb.RenderActorsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.RenderActorsResponse.deserializeBinary
);


/**
 * @param {!proto.treediagram.pb.RenderActorsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.treediagram.pb.RenderActorsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.treediagram.pb.RenderActorsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.treediagram.pb.TreeDiagramServiceClient.prototype.renderActors =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderActors',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderActors,
      callback);
};


/**
 * @param {!proto.treediagram.pb.RenderActorsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.treediagram.pb.RenderActorsResponse>}
 *     Promise that resolves to the response
 */
proto.treediagram.pb.TreeDiagramServicePromiseClient.prototype.renderActors =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/RenderActors',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_RenderActors);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.treediagram.pb.QueryEventsRequest,
 *   !proto.treediagram.pb.QueryEventsResponse>}
 */
const methodDescriptor_TreeDiagramService_QueryEvents = new grpc.web.MethodDescriptor(
  '/treediagram.pb.TreeDiagramService/QueryEvents',
  grpc.web.MethodType.UNARY,
  proto.treediagram.pb.QueryEventsRequest,
  proto.treediagram.pb.QueryEventsResponse,
  /**
   * @param {!proto.treediagram.pb.QueryEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.QueryEventsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.treediagram.pb.QueryEventsRequest,
 *   !proto.treediagram.pb.QueryEventsResponse>}
 */
const methodInfo_TreeDiagramService_QueryEvents = new grpc.web.AbstractClientBase.MethodInfo(
  proto.treediagram.pb.QueryEventsResponse,
  /**
   * @param {!proto.treediagram.pb.QueryEventsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.treediagram.pb.QueryEventsResponse.deserializeBinary
);


/**
 * @param {!proto.treediagram.pb.QueryEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.treediagram.pb.QueryEventsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.treediagram.pb.QueryEventsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.treediagram.pb.TreeDiagramServiceClient.prototype.queryEvents =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/QueryEvents',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_QueryEvents,
      callback);
};


/**
 * @param {!proto.treediagram.pb.QueryEventsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.treediagram.pb.QueryEventsResponse>}
 *     Promise that resolves to the response
 */
proto.treediagram.pb.TreeDiagramServicePromiseClient.prototype.queryEvents =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/treediagram.pb.TreeDiagramService/QueryEvents',
      request,
      metadata || {},
      methodDescriptor_TreeDiagramService_QueryEvents);
};


module.exports = proto.treediagram.pb;

