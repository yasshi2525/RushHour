import * as PIXI from "pixi.js";
import { Coordinates } from "./gamemap";
import GameModel from "../models"

export interface ApplicationProperty {
    [key: string]: any,
    app: PIXI.Application
};

export interface ModelProperty extends ApplicationProperty {
    model: GameModel
}

export interface PIXIProperty extends ModelProperty {
    zIndex: number
}

export interface SpriteProperty extends PIXIProperty {
    texture: PIXI.Texture
};

export interface AnimatedSpriteProperty extends PIXIProperty {
    animation: PIXI.Texture[],
    offset: number
}

export interface GameModelProperty extends ApplicationProperty, Coordinates {
};

export interface SpriteContainerProperty extends PIXIProperty {
    texture: PIXI.Texture
}

export interface AnimatedSpriteContainerProperty extends PIXIProperty {
    animation: PIXI.Texture[]
}

export interface AnimationProperty extends ApplicationProperty {
    filter: PIXI.Filter,
    fn: (filter: PIXI.Filter, offset: number) => void
}

export interface BorderProperty extends PIXIProperty {
    delegate: number
}

export enum ZIndex {
    NORMAL_BORDER,
    WORLD_BORDER,
    RESIDENCE,
    COMPANY,
    RAIL_EDGE,
    RAIL_NODE,
    STATION,
    GATE,
    PLATFORM,
    RAIL_LINE,
    LINE_TASK,
    HUMAN,
    PLAYER,
    ANCHOR,
    CURSOR
};