import * as PIXI from "pixi.js"

export interface LocalableProprty {
    [key: string]: any,
    app: PIXI.Application,
    // 画面中央に対応するゲームx座標
    cx: number,
    // 画面中央に対応するゲームy座標
    cy: number,
    // 画面幅に対応するゲーム座標系での長さ。2のscale乗
    scale: number
}

export interface SpriteProperty extends LocalableProprty {
    name: string
}