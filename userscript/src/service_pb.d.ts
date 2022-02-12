import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class Date extends jspb.Message {
  getYear(): number;
  setYear(value: number): Date;
  hasYear(): boolean;
  clearYear(): Date;

  getMonth(): number;
  setMonth(value: number): Date;
  hasMonth(): boolean;
  clearMonth(): Date;

  getDay(): number;
  setDay(value: number): Date;
  hasDay(): boolean;
  clearDay(): Date;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Date.AsObject;
  static toObject(includeInstance: boolean, msg: Date): Date.AsObject;
  static serializeBinaryToWriter(message: Date, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Date;
  static deserializeBinaryFromReader(message: Date, reader: jspb.BinaryReader): Date;
}

export namespace Date {
  export type AsObject = {
    year?: number,
    month?: number,
    day?: number,
  }
}

export class RenderEventRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RenderEventRequest;
  hasId(): boolean;
  clearId(): RenderEventRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderEventRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RenderEventRequest): RenderEventRequest.AsObject;
  static serializeBinaryToWriter(message: RenderEventRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderEventRequest;
  static deserializeBinaryFromReader(message: RenderEventRequest, reader: jspb.BinaryReader): RenderEventRequest;
}

export namespace RenderEventRequest {
  export type AsObject = {
    id?: string,
  }
}

export class RenderEventResponse extends jspb.Message {
  getDateStr(): string;
  setDateStr(value: string): RenderEventResponse;
  hasDateStr(): boolean;
  clearDateStr(): RenderEventResponse;

  getDate(): Date | undefined;
  setDate(value?: Date): RenderEventResponse;
  hasDate(): boolean;
  clearDate(): RenderEventResponse;

  getSnapshotsList(): Array<RenderEventResponse.Snapshot>;
  setSnapshotsList(value: Array<RenderEventResponse.Snapshot>): RenderEventResponse;
  clearSnapshotsList(): RenderEventResponse;
  addSnapshots(value?: RenderEventResponse.Snapshot, index?: number): RenderEventResponse.Snapshot;

  getPlaceStatsTotal(): RenderEventResponse.PlaceNoteCountStats | undefined;
  setPlaceStatsTotal(value?: RenderEventResponse.PlaceNoteCountStats): RenderEventResponse;
  hasPlaceStatsTotal(): boolean;
  clearPlaceStatsTotal(): RenderEventResponse;

  getPlaceStatsFinished(): RenderEventResponse.PlaceNoteCountStats | undefined;
  setPlaceStatsFinished(value?: RenderEventResponse.PlaceNoteCountStats): RenderEventResponse;
  hasPlaceStatsFinished(): boolean;
  clearPlaceStatsFinished(): RenderEventResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderEventResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RenderEventResponse): RenderEventResponse.AsObject;
  static serializeBinaryToWriter(message: RenderEventResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderEventResponse;
  static deserializeBinaryFromReader(message: RenderEventResponse, reader: jspb.BinaryReader): RenderEventResponse;
}

export namespace RenderEventResponse {
  export type AsObject = {
    dateStr?: string,
    date?: Date.AsObject,
    snapshotsList: Array<RenderEventResponse.Snapshot.AsObject>,
    placeStatsTotal?: RenderEventResponse.PlaceNoteCountStats.AsObject,
    placeStatsFinished?: RenderEventResponse.PlaceNoteCountStats.AsObject,
  }

  export class Snapshot extends jspb.Message {
    getTimestamp(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTimestamp(value?: google_protobuf_timestamp_pb.Timestamp): Snapshot;
    hasTimestamp(): boolean;
    clearTimestamp(): Snapshot;

    getNoteCount(): number;
    setNoteCount(value: number): Snapshot;
    hasNoteCount(): boolean;
    clearNoteCount(): Snapshot;

    getAddedActorsList(): Array<string>;
    setAddedActorsList(value: Array<string>): Snapshot;
    clearAddedActorsList(): Snapshot;
    addAddedActors(value: string, index?: number): Snapshot;

    getRemovedActorsList(): Array<string>;
    setRemovedActorsList(value: Array<string>): Snapshot;
    clearRemovedActorsList(): Snapshot;
    addRemovedActors(value: string, index?: number): Snapshot;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Snapshot.AsObject;
    static toObject(includeInstance: boolean, msg: Snapshot): Snapshot.AsObject;
    static serializeBinaryToWriter(message: Snapshot, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Snapshot;
    static deserializeBinaryFromReader(message: Snapshot, reader: jspb.BinaryReader): Snapshot;
  }

  export namespace Snapshot {
    export type AsObject = {
      timestamp?: google_protobuf_timestamp_pb.Timestamp.AsObject,
      noteCount?: number,
      addedActorsList: Array<string>,
      removedActorsList: Array<string>,
    }
  }


  export class PlaceNoteCountStats extends jspb.Message {
    getTotal(): number;
    setTotal(value: number): PlaceNoteCountStats;
    hasTotal(): boolean;
    clearTotal(): PlaceNoteCountStats;

    getRank(): number;
    setRank(value: number): PlaceNoteCountStats;
    hasRank(): boolean;
    clearRank(): PlaceNoteCountStats;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PlaceNoteCountStats.AsObject;
    static toObject(includeInstance: boolean, msg: PlaceNoteCountStats): PlaceNoteCountStats.AsObject;
    static serializeBinaryToWriter(message: PlaceNoteCountStats, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PlaceNoteCountStats;
    static deserializeBinaryFromReader(message: PlaceNoteCountStats, reader: jspb.BinaryReader): PlaceNoteCountStats;
  }

  export namespace PlaceNoteCountStats {
    export type AsObject = {
      total?: number,
      rank?: number,
    }
  }

}

export class RenderPlaceRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RenderPlaceRequest;
  hasId(): boolean;
  clearId(): RenderPlaceRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderPlaceRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RenderPlaceRequest): RenderPlaceRequest.AsObject;
  static serializeBinaryToWriter(message: RenderPlaceRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderPlaceRequest;
  static deserializeBinaryFromReader(message: RenderPlaceRequest, reader: jspb.BinaryReader): RenderPlaceRequest;
}

export namespace RenderPlaceRequest {
  export type AsObject = {
    id?: string,
  }
}

export class RenderPlaceResponse extends jspb.Message {
  getKnownEventCount(): number;
  setKnownEventCount(value: number): RenderPlaceResponse;
  hasKnownEventCount(): boolean;
  clearKnownEventCount(): RenderPlaceResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderPlaceResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RenderPlaceResponse): RenderPlaceResponse.AsObject;
  static serializeBinaryToWriter(message: RenderPlaceResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderPlaceResponse;
  static deserializeBinaryFromReader(message: RenderPlaceResponse, reader: jspb.BinaryReader): RenderPlaceResponse;
}

export namespace RenderPlaceResponse {
  export type AsObject = {
    knownEventCount?: number,
  }
}

export class RenderActorsRequest extends jspb.Message {
  getIdList(): Array<string>;
  setIdList(value: Array<string>): RenderActorsRequest;
  clearIdList(): RenderActorsRequest;
  addId(value: string, index?: number): RenderActorsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderActorsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RenderActorsRequest): RenderActorsRequest.AsObject;
  static serializeBinaryToWriter(message: RenderActorsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderActorsRequest;
  static deserializeBinaryFromReader(message: RenderActorsRequest, reader: jspb.BinaryReader): RenderActorsRequest;
}

export namespace RenderActorsRequest {
  export type AsObject = {
    idList: Array<string>,
  }
}

export class RenderActorsResponse extends jspb.Message {
  getItemsMap(): jspb.Map<string, RenderActorsResponse.ResponseItem>;
  clearItemsMap(): RenderActorsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RenderActorsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RenderActorsResponse): RenderActorsResponse.AsObject;
  static serializeBinaryToWriter(message: RenderActorsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RenderActorsResponse;
  static deserializeBinaryFromReader(message: RenderActorsResponse, reader: jspb.BinaryReader): RenderActorsResponse;
}

export namespace RenderActorsResponse {
  export type AsObject = {
    itemsMap: Array<[string, RenderActorsResponse.ResponseItem.AsObject]>,
  }

  export class ResponseItem extends jspb.Message {
    getKnownEventCount(): number;
    setKnownEventCount(value: number): ResponseItem;
    hasKnownEventCount(): boolean;
    clearKnownEventCount(): ResponseItem;

    getSnapshotsList(): Array<RenderActorsResponse.ResponseItem.Snapshot>;
    setSnapshotsList(value: Array<RenderActorsResponse.ResponseItem.Snapshot>): ResponseItem;
    clearSnapshotsList(): ResponseItem;
    addSnapshots(value?: RenderActorsResponse.ResponseItem.Snapshot, index?: number): RenderActorsResponse.ResponseItem.Snapshot;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ResponseItem.AsObject;
    static toObject(includeInstance: boolean, msg: ResponseItem): ResponseItem.AsObject;
    static serializeBinaryToWriter(message: ResponseItem, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ResponseItem;
    static deserializeBinaryFromReader(message: ResponseItem, reader: jspb.BinaryReader): ResponseItem;
  }

  export namespace ResponseItem {
    export type AsObject = {
      knownEventCount?: number,
      snapshotsList: Array<RenderActorsResponse.ResponseItem.Snapshot.AsObject>,
    }

    export class Snapshot extends jspb.Message {
      getDateStr(): string;
      setDateStr(value: string): Snapshot;
      hasDateStr(): boolean;
      clearDateStr(): Snapshot;

      getDate(): Date | undefined;
      setDate(value?: Date): Snapshot;
      hasDate(): boolean;
      clearDate(): Snapshot;

      getFavoriteCount(): number;
      setFavoriteCount(value: number): Snapshot;
      hasFavoriteCount(): boolean;
      clearFavoriteCount(): Snapshot;

      serializeBinary(): Uint8Array;
      toObject(includeInstance?: boolean): Snapshot.AsObject;
      static toObject(includeInstance: boolean, msg: Snapshot): Snapshot.AsObject;
      static serializeBinaryToWriter(message: Snapshot, writer: jspb.BinaryWriter): void;
      static deserializeBinary(bytes: Uint8Array): Snapshot;
      static deserializeBinaryFromReader(message: Snapshot, reader: jspb.BinaryReader): Snapshot;
    }

    export namespace Snapshot {
      export type AsObject = {
        dateStr?: string,
        date?: Date.AsObject,
        favoriteCount?: number,
      }
    }

  }

}

export class QueryEventsRequest extends jspb.Message {
  getOffset(): number;
  setOffset(value: number): QueryEventsRequest;
  hasOffset(): boolean;
  clearOffset(): QueryEventsRequest;

  getFilter(): QueryEventsRequest.EventFilter | undefined;
  setFilter(value?: QueryEventsRequest.EventFilter): QueryEventsRequest;
  hasFilter(): boolean;
  clearFilter(): QueryEventsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: QueryEventsRequest): QueryEventsRequest.AsObject;
  static serializeBinaryToWriter(message: QueryEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryEventsRequest;
  static deserializeBinaryFromReader(message: QueryEventsRequest, reader: jspb.BinaryReader): QueryEventsRequest;
}

export namespace QueryEventsRequest {
  export type AsObject = {
    offset?: number,
    filter?: QueryEventsRequest.EventFilter.AsObject,
  }

  export class EventFilter extends jspb.Message {
    getActorIdsList(): Array<string>;
    setActorIdsList(value: Array<string>): EventFilter;
    clearActorIdsList(): EventFilter;
    addActorIds(value: string, index?: number): EventFilter;

    getPlaceId(): string;
    setPlaceId(value: string): EventFilter;
    hasPlaceId(): boolean;
    clearPlaceId(): EventFilter;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EventFilter.AsObject;
    static toObject(includeInstance: boolean, msg: EventFilter): EventFilter.AsObject;
    static serializeBinaryToWriter(message: EventFilter, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EventFilter;
    static deserializeBinaryFromReader(message: EventFilter, reader: jspb.BinaryReader): EventFilter;
  }

  export namespace EventFilter {
    export type AsObject = {
      actorIdsList: Array<string>,
      placeId?: string,
    }
  }

}

export class QueryEventsResponse extends jspb.Message {
  getEventsList(): Array<QueryEventsResponse.Event>;
  setEventsList(value: Array<QueryEventsResponse.Event>): QueryEventsResponse;
  clearEventsList(): QueryEventsResponse;
  addEvents(value?: QueryEventsResponse.Event, index?: number): QueryEventsResponse.Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: QueryEventsResponse): QueryEventsResponse.AsObject;
  static serializeBinaryToWriter(message: QueryEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryEventsResponse;
  static deserializeBinaryFromReader(message: QueryEventsResponse, reader: jspb.BinaryReader): QueryEventsResponse;
}

export namespace QueryEventsResponse {
  export type AsObject = {
    eventsList: Array<QueryEventsResponse.Event.AsObject>,
  }

  export class Event extends jspb.Message {
    getId(): string;
    setId(value: string): Event;
    hasId(): boolean;
    clearId(): Event;

    getName(): string;
    setName(value: string): Event;
    hasName(): boolean;
    clearName(): Event;

    getDateStr(): string;
    setDateStr(value: string): Event;
    hasDateStr(): boolean;
    clearDateStr(): Event;

    getDate(): Date | undefined;
    setDate(value?: Date): Event;
    hasDate(): boolean;
    clearDate(): Event;

    getFinished(): boolean;
    setFinished(value: boolean): Event;
    hasFinished(): boolean;
    clearFinished(): Event;

    getLastNoteCount(): number;
    setLastNoteCount(value: number): Event;
    hasLastNoteCount(): boolean;
    clearLastNoteCount(): Event;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Event.AsObject;
    static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
    static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Event;
    static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
  }

  export namespace Event {
    export type AsObject = {
      id?: string,
      name?: string,
      dateStr?: string,
      date?: Date.AsObject,
      finished?: boolean,
      lastNoteCount?: number,
    }
  }

}

