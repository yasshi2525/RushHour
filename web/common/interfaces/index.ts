import GameModel from "../models";
import { MenuStatus } from "../../state";

export interface ActionPayload {
    [index: string]: any,
    status: boolean,
    timestamp: number,
    results: any
};

export interface AppBarProperty {
    [key: string]: any,
    readOnly: boolean,
    displayName: string | undefined,
    image: string | undefined
    hue: number | undefined,
};

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