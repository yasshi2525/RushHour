export enum Method {
  GET = "GET",
  PUT = "PUT",
  POST = "POST",
  DELETE = "DELETE"
}

export interface Endpoint {
  method?: Method;
  url: string;
  auth?: boolean;
}

export interface FetchMap {
  cx: number;
  cy: number;
  scale: number;
  delegate: number;
}
export const fetchMap: Endpoint = {
  url: "/api/v1/gamemap"
};

export const players: Endpoint = {
  url: "/api/v1/players"
};

export const signout: Endpoint = {
  url: "/api/v1/signout",
  method: Method.POST,
  auth: true
};

export interface Depart {
  x: number;
  y: number;
  scale: number;
}

export interface Extend {
  x: number;
  y: number;
  rnid: number;
  scale: number;
}

export interface Connect {
  from: number;
  to: number;
  scale: number;
}

export const rail: { [index: string]: Endpoint } = {
  depart: { url: "/api/v1/rail_nodes", method: Method.POST, auth: true },
  extend: { url: "/api/v1/rail_nodes/extend", method: Method.POST, auth: true },
  connect: {
    url: "/api/v1/rail_nodes/connect",
    method: Method.POST,
    auth: true
  },
  destroy: {
    url: "/api/v1/rail_nodes",
    method: Method.DELETE,
    auth: true
  }
};

export const game: { [index: string]: Endpoint } = {
  status: {
    url: "/api/v1/game"
  }
};
