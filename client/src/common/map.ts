import { useContext, useCallback } from "react";
import { Entity, GameMap } from "interfaces";
import { players } from "interfaces/endpoint";
import { ServerErrors } from "interfaces/error";
import GameModel from "models";
import { Http, http, useHttpTask } from "./http";
import ModelContext from "./model";

const makeArray = (err: ServerErrors | undefined): Http =>
  err ? [err, undefined] : [undefined, undefined];

export const usePlayers = (model: GameModel) => {
  return useHttpTask(players, (d: Entity[]) =>
    model.gamemap.mergeChildren("players", d)
  );
};

interface FetchMapRequest {
  cx: number;
  cy: number;
  scale: number;
  delegate: number;
}

export const fetchMap = (model: GameModel) => async (
  signal: AbortSignal
): Promise<Http> => {
  const [err, data] = await http<FetchMapRequest, GameMap>({
    url: "/api/v1/gamemap",
    args: { ...model.coord, delegate: model.delegate },
    signal
  });
  if (data) {
    model.gamemap.mergeAll(data);
  }
  return makeArray(err);
};

export const useMapReload = (props: { watch: any; errRef: any }) => {
  const [watch, setWatch] = useState(watch);
  const model = useContext(ModelContext);
  useCallback();
};
