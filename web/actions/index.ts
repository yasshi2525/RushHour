import { createAsyncAction } from "typesafe-actions";
import GameContainer from "../common";
import GameModel from "../common/models";
import { GameMap } from "../state";

export interface ModelRequest {
    model: GameModel,
    dispatch: any
}

export interface GameMapRequest extends ModelRequest {
}

export interface PointRequest extends ModelRequest {
    oid: number,
    x: number,
    y: number,
    scale: number
};

export interface ExtendRequest extends PointRequest {
    rnid: number
}

export interface ActionPayload {
    [index: string]: any,
    status: boolean,
    timestamp: number,
    results: any
};

export interface GameMapResponse extends ActionPayload {
    results: GameMap
};

export interface GameResponse extends ActionPayload {
    results: {[index: string]: any}
};

export const initPIXI = createAsyncAction("INIT_PIXI_REQUESTED", "INIT_PIXI_SUCCEEDED", "INIT_PIXI_FAILED")<GameContainer, GameContainer, Error>();
export const fetchMap = createAsyncAction("FETCH_MAP_REQUESTED", "FETCH_MAP_SUCCEEDED", "FETCH_MAP_FAILED")<GameMapRequest, GameMapResponse, Error>();
export const diffMap = createAsyncAction("DIFF_MAP_REQUESTED", "DIFF_MAP_SUCCEEDED", "DIFF_MAP_FAILED")<GameMapRequest, GameMapResponse, Error>();
export const players = createAsyncAction("PLAYERS_REQUESTED", "PLAYERS_SUCCEEDED", "PLAYERS_FAILED")<ModelRequest, GameResponse, Error>();
export const depart = createAsyncAction("DEPART_REQUESTED", "DEPART_SUCCEEDED", "DEPART_FAILED")<PointRequest, GameResponse, Error>();
export const extend = createAsyncAction("EXTEND_REQUESTED", "EXTEND_SUCCEEDED", "EXTEND_FAILED")<ExtendRequest, GameResponse, Error>();