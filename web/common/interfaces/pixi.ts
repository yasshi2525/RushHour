import { Coordinates } from "./gamemap";

export interface ApplicationProperty {
    [key: string]: any,
    app: PIXI.Application,
};

export interface SpriteProperty extends ApplicationProperty {
    texture: PIXI.Texture
};

export interface AnimatedSpriteProperty extends ApplicationProperty {
    animation: PIXI.Texture[],
    offset: number
}

export interface GameModelProperty extends ApplicationProperty, Coordinates {
};

export interface SpriteContainerProperty extends ApplicationProperty {
    texture: PIXI.Texture
}

export interface AnimatedSpriteContainerProperty extends ApplicationProperty {
    animation: PIXI.Texture[]
}

export interface AnimationProperty extends ApplicationProperty {
    filter: PIXI.Filter,
    fn: (filter: PIXI.Filter, offset: number) => void
}

export interface ResourceAttachable {
    attach(textures: {[index: string]: PIXI.Texture}): void
}