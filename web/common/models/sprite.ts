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

export default class extends PointModel implements Monitorable {
    protected name: string;
    protected container: PIXI.Container;
    protected loader: PIXI.loaders.Loader;
    protected sprite?: PIXI.Sprite;

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

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            let resource = this.loader.resources[this.name];
            this.sprite = new PIXI.Sprite(resource ? resource.texture : undefined);
            this.setupSprite();
            this.container.addChild(this.sprite);
        });
    }

    protected setupSprite() {
        if (this.sprite !== undefined) {
            this.sprite.setTransform(this.props.x, this.props.y, this.props.scale, this.props.scale);
            this.sprite.anchor.set(this.props.anchor.x, this.props.anchor.y);
            this.sprite.alpha = this.props.alpha;
        }
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => ifExists<PIXI.Sprite>(this.sprite, sprite => {
            this.container.removeChild(sprite);
            this.sprite = undefined;
        }));
    }

    getSprite() {
        return this.sprite;
    }
}