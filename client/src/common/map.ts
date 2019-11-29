import { useContext, useCallback } from "react";
import { Entity, GameMap } from "interfaces";
import { FetchMap, players, fetchMap } from "interfaces/endpoint";
import { useHttpTask } from "./http";
import ModelContext from "./model";

export const usePlayers = () => {
  const model = useContext(ModelContext);
  const [_players] = useHttpTask(players, (d: Entity[]) =>
    model.gamemap.mergeChildren("players", d)
  );
  return useCallback(() => _players(undefined), []);
};

/**
 * ```
 * const fire = useReload();
 * fire();
 * ```
 */
export const useReload = () => {
  const model = useContext(ModelContext);
  const [fire] = useHttpTask<FetchMap, GameMap>(fetchMap, d =>
    model.gamemap.mergeAll(d)
  );
  return () => fire({ ...model.coord, delegate: model.delegate });
};
