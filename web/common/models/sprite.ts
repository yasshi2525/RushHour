import * as PIXI from "pixi.js"; 
import { Monitorable } from "../interfaces/monitor";
import PointModel from "./point";

/**
 * x, y座標はPointModelで初期化済み
 */
const defaultValues: {
    anchor: {x: number, y: number},
    scale: number,
    alpha: number
} = { 
    anchor: { x: 0.5, y: 0.5 }, 
    scale: 1,
    alpha: 1
};

const ifExists = <T>(target: any, callback: (target: T) => void) =>
    (target !== undefined) ? { execute: callback(target)} : function(){};

export default abstract class extends PointModel implements Monitorable {
    name: string;
    container: PIXI.Container;
    loader: PIXI.loaders.Loader;
    sprite?: PIXI.Sprite;

    constructor(options: { name: string, container: PIXI.Container, loader: PIXI.loaders.Loader }) {
        super();
        this.name = options.name;
        this.container = options.container;
        this.loader = options.loader;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();

        // 値の更新時、Spriteを更新するように設定
        this.addUpdateCallback("x", (value: number) => ifExists<PIXI.Sprite>(this.sprite, 
            sprite => sprite.x = value));
        this.addUpdateCallback("y", (value: number) => ifExists<PIXI.Sprite>(this.sprite, 
            sprite => sprite.y = value));
        this.addUpdateCallback("anchor", (value: {x: number, y: number}) => ifExists<PIXI.Sprite>(this.sprite, 
            sprite => sprite.anchor.set(value.x, value.y)));
        this.addUpdateCallback("scale", (value: number) => ifExists<PIXI.Sprite>(this.sprite,
            sprite => sprite.scale.set(value, value)));
        this.addUpdateCallback("alpha", (value: number) => ifExists<PIXI.Sprite>(this.sprite,
            sprite => sprite.alpha = value));
    }

    protected setupSprite() {
        if (this.sprite !== undefined) {
            this.sprite.setTransform(this.props.x, this.props.y, this.props.scale, this.props.scale);
            this.sprite.anchor.set(this.props.anchor.x, this.props.anchor.y);
            this.sprite.alpha = this.props.alpha;
        }
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.beforeCallbacks.push(() => {
            this.sprite = new PIXI.Sprite(this.loader.resources[this.name].texture);
            this.setupSprite();
            this.container.addChild(this.sprite);
        });
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.afterCallbacks.push(() => ifExists<PIXI.Sprite>(this.sprite, sprite => this.container.removeChild(sprite)));
    }
}