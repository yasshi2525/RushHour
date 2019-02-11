import * as PIXI from "pixi.js";
import * as Filters from "pixi-filters";
import { Monitorable } from "../interfaces/monitor";
import { PointModel } from "./geo";

/**
 * x, y座標はPointModelで初期化済み
 */
const defaultValues: {
    anchor: {x: number, y: number},
    scale: number,
    alpha: number
} = { 
    anchor: { x: 0.5, y: 0.5 }, 
    scale: 0.5 * window.devicePixelRatio,
    alpha: 1
};

const outlineOpts = {
    thickness: {max: 4, min: 2},
    color: 0xeeeeee,
    round: 240
};

const shadowOpts: PIXI.filters.DropShadowFilterOptions = {
    color: 0x00ffff
}

export default class extends PointModel implements Monitorable {
    protected name: string;
    protected sprite: PIXI.Sprite;
    protected outline: PIXI.filters.OutlineFilter;
    protected shadow: PIXI.filters.DropShadowFilter;
    protected tick: number

    constructor(options: { name: string, app: PIXI.Application }) {
        super(options);
        this.name = options.name;
        let resource = this.app.loader.resources[this.name];
        this.sprite = new PIXI.Sprite(resource ? resource.texture : undefined);
        this.outline = new Filters.OutlineFilter(outlineOpts.thickness.min, outlineOpts.color);
        this.shadow = new Filters.DropShadowFilter(shadowOpts);
        this.sprite.filters = [this.outline, this.shadow];
        this.tick = 0
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();

        // 値の更新時、Spriteを更新するように設定
        this.addUpdateCallback("x", (value: number) => this.sprite.x = value);
        this.addUpdateCallback("y", (value: number) => this.sprite.y = value);
        this.addUpdateCallback("anchor", (value: {x: number, y: number}) => this.sprite.anchor.set(value.x, value.y));
        this.addUpdateCallback("scale", (value: number) => this.sprite.scale.set(value, value));
        this.addUpdateCallback("alpha", (value: number) => this.sprite.alpha = value);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.setupSprite();
            this.container.addChild(this.sprite);
            this.app.ticker.add((d: number) => this.animate(d))
        });
    }

    protected setupSprite() {
        this.sprite.setTransform(this.props.x, this.props.y, this.props.scale, this.props.scale);
        this.sprite.anchor.set(this.props.anchor.x, this.props.anchor.y);
        this.sprite.alpha = this.props.alpha;
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.sprite));
    }

    getSprite() {
        return this.sprite;
    }

    protected animate(delta: number) {
        this.tick += delta;

        let dx = (this.tick / outlineOpts.round) - Math.floor(this.tick / outlineOpts.round);
        if (dx > 0.5) {
            dx = 1 - dx;
        }
        this.outline.thickness = dx * outlineOpts.thickness.min + (1 - dx) * outlineOpts.thickness.max;
        this.outline.thickness *= window.devicePixelRatio;
        this.app.renderer.render(this.sprite);
    }
}