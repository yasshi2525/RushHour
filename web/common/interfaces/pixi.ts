import * as PIXI from "pixi.js"

export interface ApplicationProperty {
    [key: string]: any,
    app: PIXI.Application,
}

export interface SpriteProperty extends ApplicationProperty {
    name: string
}