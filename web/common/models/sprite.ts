import * as PIXI from "pixi.js";
import * as Filters from "pixi-filters";
import { Monitorable } from "../interfaces/monitor";
import { SpriteProperty } from "../interfaces/pixi";
import PointModel from "./point";

/**
 * x, y座標はPointModelで初期化済み
 */
const defaultValues: {
    anchor: {x: number, y: number},
    spscale: number,
    alpha: number
} = { 
    anchor: { x: 0.5, y: 0.5 }, 
    spscale: 0.5,
    alpha: 1
};

const outlineOpts = {
    width: {min: 0, max: 4},
    color: 0xeeeeee,
    round: 5000
};

const shadowOpts: PIXI.filters.DropShadowFilterOptions = {
    color: 0x00ffff,
    distance: 5
}

export default class extends PointModel implements Monitorable {
    protected name: string;
    protected sprite: PIXI.Sprite;
    protected outline: PIXI.filters.OutlineFilter;
    protected shadow: PIXI.filters.DropShadowFilter;
    /**
     * インスタンス作成からの累計時間
     */
    protected tick: number;
    /**
     * 明滅率(0-1)
     */
    protected offset: number;

    constructor(options: SpriteProperty) {
        super(options);
        this.name = options.name;
        let resource = this.app.loader.resources[this.name];
        this.sprite = new PIXI.Sprite(resource ? resource.texture : undefined);
        this.outline = new Filters.OutlineFilter(
            outlineOpts.width.min / this.app.renderer.resolution, 
            outlineOpts.color);
        this.shadow = new Filters.DropShadowFilter(shadowOpts);
        this.shadow.distance /= this.app.renderer.resolution;
        this.sprite.filters = [this.outline, this.shadow];
        this.tick = 0;
        this.offset = 0;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.sprite.anchor.set(this.props.anchor.x, this.props.anchor.y);
            this.sprite.alpha = this.props.alpha;
            this.sprite.scale.set(this.props.spscale, this.props.spscale)
            this.container.addChild(this.sprite);
            this.app.ticker.add(() => this.flash())
        });
    }
    
    setupUpdateCallback() {
        super.setupUpdateCallback();

        // 値の更新時、Spriteを更新するように設定
        this.addUpdateCallback("anchor", (value: {x: number, y: number}) => this.sprite.anchor.set(value.x, value.y));
        this.addUpdateCallback("spscale", (value: number) => this.sprite.scale.set(value, value));
        this.addUpdateCallback("alpha", (value: number) => this.sprite.alpha = value);
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.sprite));
    }

    getSprite() {
        return this.sprite;
    }

    beforeRender() {
        super.beforeRender();
        this.sprite.x = this.current.x;
        this.sprite.y = this.current.y;

        this.outline.thickness = (this.offset * outlineOpts.width.min
                                    + (1- this.offset) * outlineOpts.width.max)
                                    / this.app.renderer.resolution;
    }

    protected flash() {
        this.tick += this.app.ticker.elapsedMS;
        let ratio = (this.tick % outlineOpts.round) / outlineOpts.round;
        this.offset = Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
        this.beforeRender();
    }
}