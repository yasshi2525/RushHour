import * as PIXI from "pixi.js"
import { Coordinates } from "./gamemap";

export interface ApplicationProperty {
    [key: string]: any,
    app: PIXI.Application,
};

export interface ContainerProperty extends ApplicationProperty {
    container: PIXI.Container
}

export interface SpriteProperty extends ContainerProperty {
    texture: PIXI.Texture
};

export interface GameModelProperty extends ApplicationProperty, Coordinates {
    textures: {[index: string]: PIXI.Texture}
};

export interface SpriteContainerProperty extends ApplicationProperty {
    texture: PIXI.Texture
}