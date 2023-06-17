// @generated by protoc-gen-connect-web v0.10.0 with parameter "target=ts,import_extension=none"
// @generated from file service.proto (package treediagram.pb, syntax proto2)
/* eslint-disable */
// @ts-nocheck

import { QueryEventsRequest, QueryEventsResponse, RenderActorsRequest, RenderActorsResponse, RenderEventRequest, RenderEventResponse, RenderPlaceRequest, RenderPlaceResponse } from "./service_pb";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service treediagram.pb.TreeDiagramService
 */
export const TreeDiagramService = {
  typeName: "treediagram.pb.TreeDiagramService",
  methods: {
    /**
     * @generated from rpc treediagram.pb.TreeDiagramService.RenderEvent
     */
    renderEvent: {
      name: "RenderEvent",
      I: RenderEventRequest,
      O: RenderEventResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc treediagram.pb.TreeDiagramService.RenderPlace
     */
    renderPlace: {
      name: "RenderPlace",
      I: RenderPlaceRequest,
      O: RenderPlaceResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc treediagram.pb.TreeDiagramService.RenderActors
     */
    renderActors: {
      name: "RenderActors",
      I: RenderActorsRequest,
      O: RenderActorsResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc treediagram.pb.TreeDiagramService.QueryEvents
     */
    queryEvents: {
      name: "QueryEvents",
      I: QueryEventsRequest,
      O: QueryEventsResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

