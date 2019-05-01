import { createAsyncAction } from "typesafe-actions";
import { Coordinates, GameMap } from "../state";

export interface ActionPayload {
    status: boolean,
    timestamp: number,
    results: any
};

export interface GameMapResponse extends ActionPayload {
    results: GameMap
};

export const fetchMap = createAsyncAction("FETCH_MAP_REQUESTED", "FETCH_MAP_SUCCEESSED", "FETCH_MAP_FAILED")<Coordinates, GameMapResponse ,Error>();
export const diffMap = createAsyncAction("DIFF_MAP_REQUESTED", "DIFF_MAP_SUCCEEDED", "DIFF_MAP_FAILED")<Coordinates, GameMapResponse ,Error>();