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

export interface SpriteProperty extends ModelProperty {
    texture: PIXI.Texture
};

export interface AnimatedSpriteProperty extends ModelProperty {
    animation: PIXI.Texture[],
    offset: number
}

export interface GameModelProperty extends ApplicationProperty, Coordinates {
};

export interface SpriteContainerProperty extends ModelProperty {
    texture: PIXI.Texture
}

export interface AnimatedSpriteContainerProperty extends ModelProperty {
    animation: PIXI.Texture[]
}

export interface AnimationProperty extends ApplicationProperty {
    filter: PIXI.Filter,
    fn: (filter: PIXI.Filter, offset: number) => void
}
