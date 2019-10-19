import { createAsyncAction, createAction } from "typesafe-actions";
import GameContainer from "../common";
import { Entry } from "../common/interfaces";
import GameModel from "../common/models";
import { GameMap, MenuStatus } from "../state";

export interface Request {

}

export interface ResourceRequest extends Entry {
}

export interface LoginRequest {
    id: string,
    password: string
}

export interface RegisterRequest extends LoginRequest {
    name: string,
    hue: number
}

export interface ModelRequest {
    model: GameModel
}

export interface UserActionRequest extends ModelRequest {
    scale: number
}

export interface PointRequest extends UserActionRequest {
    x: number,
    y: number
};

export interface ExtendRequest extends PointRequest {
    rnid: number
}

export interface ConnectRequest extends UserActionRequest {
    from: number,
    to: number
}

export interface DestroyRequest extends UserActionRequest{
    resource: string,
    id: number,
    cid: number
}

export interface MenuRequest extends ModelRequest {
    menu: MenuStatus
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

export interface MenuResponse {
    menu: MenuStatus
}

export const initPIXI = createAsyncAction("INIT_PIXI_REQUESTED", "INIT_PIXI_SUCCEEDED", "INIT_PIXI_FAILED")<GameContainer, GameContainer, Error>();
export const fetchMap = createAsyncAction("FETCH_MAP_REQUESTED", "FETCH_MAP_SUCCEEDED", "FETCH_MAP_FAILED")<ModelRequest, GameMapResponse, Error>();
export const login = createAsyncAction("LOGGEDIN_REQUESTED", "LOGGEDIN_SUCCEEDED", "LOGGEDIN_FAILED")<LoginRequest, GameResponse, Error>();
export const resetLoginError = createAction("RESET_LOGIN_ERROR");
export const register = createAsyncAction("REGISTER_REQUESTED", "REGISTER_SUCCEEDED", "REGISTER_FAILED")<RegisterRequest, GameResponse, Error>();
export const settings = createAsyncAction("SETTINGS_REQUESTED", "SETTINGS_SUCCEEDED", "SETTINGS_FAILED")<Request, GameResponse, Error>();
export const editSettings = createAsyncAction("EDIT_SETTINGS_REQUESTED", "EDIT_SETTINGS_SUCCEEDED", "EDIT_SETTINGS_FAILED")<ResourceRequest, GameResponse, Error>();
export const players = createAsyncAction("PLAYERS_REQUESTED", "PLAYERS_SUCCEEDED", "PLAYERS_FAILED")<ModelRequest, GameResponse, Error>();
export const depart = createAsyncAction("DEPART_REQUESTED", "DEPART_SUCCEEDED", "DEPART_FAILED")<PointRequest, GameResponse, Error>();
export const extend = createAsyncAction("EXTEND_REQUESTED", "EXTEND_SUCCEEDED", "EXTEND_FAILED")<ExtendRequest, GameResponse, Error>();
export const connect = createAsyncAction("CONNECT_REQUESTED", "CONNECT_SUCCEEDED", "CONNECT_FAILED")<ConnectRequest, GameResponse, Error>();
export const destroy = createAsyncAction("DESTROY_REQUESTED", "DESTROY_SUCCEEDED", "DESTROY_FAILED")<DestroyRequest, GameResponse, Error>();
export const setMenu = createAsyncAction("MENU_REQUESTED", "MENU_SUCCEEDED", "MENU_FAILED")<MenuRequest, MenuResponse, Error>();
export const gameStatus = createAsyncAction("GAMESTATUS_REQUESTED", "GAMESTATUS_SUCCEEDED", "GAMESTATUS_FAILED")<Request, GameResponse, Error>();
export const inOperation = createAsyncAction("OPERATION_REQUESTED", "OPERATION_SUCCEEDED", "OPERATION_FAILED")<ResourceRequest, GameResponse, Error>();
export const purgeUserData = createAsyncAction("PURGEDATA_REQUESTED", "PURGEDATA_SUCCEEDED", "PURGEDATA_FAILED")<ResourceRequest, GameResponse, Error>();