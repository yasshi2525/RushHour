import { ReactNode } from "react";
import GameModel from "../models";
import { MenuStatus } from "../../state";

export interface ActionPayload {
    [index: string]: any,
    status: boolean,
    timestamp: number,
    results: any
};

export interface BaseProperty {
    [key: string]: any,
    children: ReactNode | undefined
}

export interface UserInfo {
    [key: string]: string | number | boolean,
    id: number,
    name: string,
    image: string,
    hue: number,
    admin: boolean
}

export interface GameComponentProperty {
    model: GameModel,
};

export interface CanvasProperty extends GameComponentProperty {
    [key: string]: any,
    isFetchRequired: boolean,
    isPlayerFetched: boolean,
    dispatch: any
};

export interface MenuProperty {
    [key: string]: any,
    model: GameModel,
    menu: MenuStatus,
    setMenu: (opts: {model: GameModel, menu: MenuStatus}) => void
};

export interface Entry {
    key: string
    value: any
}

export interface AsyncStatus {
    waiting: boolean
    value: any
}