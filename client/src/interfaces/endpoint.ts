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

export const game: { [index: string]: Endpoint } = {
  status: {
    url: "/api/v1/game"
  }
};
