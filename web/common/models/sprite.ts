import * as PIXI from "pixi.js";
import * as Filters from "pixi-filters";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { SpriteProperty, SpriteContainerProperty } from "../interfaces/pixi";
import { PointModel, PointContainer } from "./point";

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

export abstract class SpriteModel extends PointModel implements Monitorable {
    protected sprite: PIXI.Sprite;
      
    constructor(options: SpriteProperty) {
        super(options);
        this.sprite = new PIXI.Sprite(options.texture);
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
    }
}

export abstract class SpriteContainer<T extends SpriteModel> extends PointContainer<T> implements MonitorContrainer {
    protected outline: PIXI.filters.OutlineFilter;
    protected shadow: PIXI.filters.DropShadowFilter;
    
    constructor(
        options: SpriteContainerProperty,
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: {[index:string]: {}}) {
        super(options, newInstance, newInstanceOptions);

        this.outline = new Filters.OutlineFilter(
            outlineOpts.width.min / this.app.renderer.resolution, 
            outlineOpts.color);
        this.outline.resolution = this.app.renderer.resolution;
        this.shadow = new Filters.DropShadowFilter(shadowOpts);
        this.shadow.resolution = this.app.renderer.resolution;

        this.container.filters = [this.outline, this.shadow];

        this.childOptions.texture = options.texture;
    }

    beforeRender() {
        super.beforeRender();
        this.outline.thickness = (this.offset * outlineOpts.width.min
            + (1- this.offset) * outlineOpts.width.max);
    }
}