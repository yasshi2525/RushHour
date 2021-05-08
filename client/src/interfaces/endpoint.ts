import {
  FlatObject,
  SerializableObject,
  HashContainer,
  Entity,
  Delegatable,
  EMPTY_DLG
} from "interfaces";

export enum Method {
  GET = "GET",
  PUT = "PUT",
  POST = "POST",
  DELETE = "DELETE"
}

export type HttpMethod = Exclude<Method, Method.GET>;

type _URL = { url: string };
type _HTTP = { method: HttpMethod };
type _AUTH = { auth: true };
type _ARGS<T> = { args: T };

type Primitive = boolean | number | string;

export type SerializableObject = {
  [index: string]: Primitive | Array<Primitive> | Object | Array<Object>;
  [index: number]: Primitive | Array<Primitive> | Object | Array<Object>;
};

type _PAYLOAD<T extends SerializableObject> = { payload: T };

export type FlatObject = {
  [index: string]: Primitive | Array<Primitive>;
  [index: number]: Primitive | Array<Primitive>;
};

export type GetEndpoint<
  I extends FlatObject = {},
  O extends SerializableObject = {}
> = _URL & _ARGS<I> & _PAYLOAD<O>;
export type GetAuthEndpoint<
  I extends FlatObject = {},
  O extends SerializableObject = {}
> = _URL & _AUTH & _ARGS<I> & _PAYLOAD<O>;
export type HttpEndpoint<
  I extends SerializableObject = {},
  O extends SerializableObject = {}
> = _URL & _HTTP & _ARGS<I> & _PAYLOAD<O>;
export type HttpAuthEndpoint<
  I extends SerializableObject = {},
  O extends SerializableObject = {}
> = _URL & _HTTP & _AUTH & _ARGS<I> & _PAYLOAD<O>;

export interface FetchMap {
  [index: string]: number;
  x: number;
  y: number;
  scale: number;
  delegate: number;
}
export type IFetchMapResponseKeys =
  | "residences"
  | "companies"
  | "rail_nodes"
  | "rail_edges";
export const FetchMapResponseKeys: IFetchMapResponseKeys[] = [
  "residences",
  "companies",
  "rail_nodes",
  "rail_edges"
];
export type FetchMapResponse = {
  [index in IFetchMapResponseKeys]: HashContainer<Delegatable>;
} & {
  timestamp: number;
};

export const fetchMap: GetEndpoint<FetchMap, FetchMapResponse> = {
  url: "/api/v1/gamemap",
  args: { x: 0, y: 0, scale: 0, delegate: 0 },
  payload: {
    residences: {},
    companies: {},
    rail_nodes: {},
    rail_edges: {},
    timestamp: 0
  }
};

export interface Player extends Entity {
  [index: string]: string | number;
  image: string;
  hue: number;
}

type IPlayersResponseKeys = "players";

export type PlayersResponse = {
  [index in IPlayersResponseKeys]: HashContainer<Player>;
};

export const players: GetEndpoint<{}, PlayersResponse> = {
  url: "/api/v1/players",
  args: {},
  payload: { players: {} }
};

export const signout: HttpAuthEndpoint = {
  url: "/api/v1/signout",
  method: Method.POST,
  auth: true,
  args: {},
  payload: {}
};

type Depart = {
  x: number;
  y: number;
  scale: number;
};

type DepartResponse = { rn: Delegatable };
type ConnectResponse = { e1: Delegatable; e2: Delegatable };
type ExtendResponse = { rn: Delegatable; e1: Delegatable; e2: Delegatable };

type Extend = {
  x: number;
  y: number;
  rnid: number;
  scale: number;
};

type Connect = {
  from: number;
  to: number;
  scale: number;
};

export const rail: {
  depart: HttpAuthEndpoint<Depart, DepartResponse>;
  extend: HttpAuthEndpoint<Extend, ExtendResponse>;
  connect: HttpAuthEndpoint<Connect, ConnectResponse>;
  destroy: HttpAuthEndpoint<{ id: number }>;
} = {
  depart: {
    url: "/api/v1/rail_nodes",
    method: Method.POST,
    auth: true,
    args: { x: 0, y: 0, scale: 0 },
    payload: { rn: EMPTY_DLG }
  },
  extend: {
    url: "/api/v1/rail_nodes/extend",
    method: Method.POST,
    auth: true,
    args: { rnid: 0, x: 0, y: 0, scale: 0 },
    payload: { rn: EMPTY_DLG, e1: EMPTY_DLG, e2: EMPTY_DLG }
  },
  connect: {
    url: "/api/v1/rail_nodes/connect",
    method: Method.POST,
    auth: true,
    args: { from: 0, to: 0, scale: 0 },
    payload: { e1: EMPTY_DLG, e2: EMPTY_DLG }
  },
  destroy: {
    url: "/api/v1/rail_nodes",
    method: Method.DELETE,
    auth: true,
    args: { id: 0 },
    payload: { id: 0 }
  }
};

export type ConfigResponse = {
  min_scale: number;
  max_scale: number;
};

export const game: {
  status: GetEndpoint<{}, { status: boolean }>;
  const: GetEndpoint<{}, ConfigResponse>;
} = {
  status: {
    url: "/api/v1/game",
    args: {},
    payload: { status: false }
  },
  const: {
    url: "/api/v1/game/const",
    args: {},
    payload: { min_scale: 0, max_scale: 0 }
  }
};
