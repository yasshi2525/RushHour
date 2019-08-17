import { createAsyncAction, createAction } from "typesafe-actions";
import { Coordinates, GameMap } from "../state";

export interface GameMapRequest extends Coordinates {
    delegate: number
}

export interface PointRequest {
    oid: number,
    x: number,
    y: number
};

export interface ActionPayload {
    status: boolean,
    timestamp: number,
    results: any
};

export interface GameMapResponse extends ActionPayload {
    results: GameMap
};

export interface DepartResponse extends ActionPayload {
    results: {oid: number, id: number, x: number, y: number}
};

export const fetchMap = createAsyncAction("FETCH_MAP_REQUESTED", "FETCH_MAP_SUCCEESSED", "FETCH_MAP_FAILED")<GameMapRequest, GameMapResponse, Error>();
export const diffMap = createAsyncAction("DIFF_MAP_REQUESTED", "DIFF_MAP_SUCCEEDED", "DIFF_MAP_FAILED")<GameMapRequest, GameMapResponse, Error>();
export const cancelEditting = createAction("CANCEL");
export const startDeparture = createAction("START_DEPT");
export const depart = createAsyncAction("DEPART_REQUESTED", "DEPART_SUCCEEDED", "DEPART_FAILED")<PointRequest, DepartResponse, Error>();