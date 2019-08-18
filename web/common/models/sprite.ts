import * as PIXI from "pixi.js";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { SpriteProperty, SpriteContainerProperty, AnimatedSpriteProperty, AnimatedSpriteContainerProperty } from "../interfaces/pixi";
import { PointModel, PointContainer } from "./point";

/**
 * x, y座標はPointModelで初期化済み
 */
const defaultValues: {
    spscale: number,
    alpha: number
} = { 
    spscale: 0.5,
    alpha: 1
};

export abstract class SpriteModel extends PointModel implements Monitorable {
    sprite: PIXI.Sprite;
      
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
            this.sprite.anchor.set(0.5, 0.5);
            this.sprite.alpha = this.props.alpha;
            this.sprite.scale.set(this.props.spscale, this.props.spscale)
            this.container.addChild(this.sprite);
        });
    }
    
    setupUpdateCallback() {
        super.setupUpdateCallback();

        // 値の更新時、Spriteを更新するように設定
        this.addUpdateCallback("spscale", (value: number) => this.sprite.scale.set(value, value));
        this.addUpdateCallback("alpha", (value: number) => this.sprite.alpha = value);
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.sprite));
    }

    beforeRender() {
        super.beforeRender();
        if (this.current === undefined) {
            this.sprite.visible = false;
        } else {
            this.sprite.visible = true;
            this.sprite.x = this.current.x;
            this.sprite.y = this.current.y;
        }
    }
}

export abstract class SpriteContainer<T extends SpriteModel> extends PointContainer<T> implements MonitorContrainer {
    
    constructor(
        options: SpriteContainerProperty,
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: {[index:string]: {}}) {
        super(options, newInstance, newInstanceOptions);
        
        this.childOptions.texture = options.texture;
    }
}

export abstract class AnimatedSpriteModel extends SpriteModel implements Monitorable {
    constructor(options: AnimatedSpriteProperty) {
        super({ texture: PIXI.Texture.EMPTY, ...options });
        
        let sprite = new PIXI.AnimatedSprite(options.animation);
        sprite.gotoAndPlay(options.offset);
        this.sprite = sprite;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues({ spscale: 1.0 });
    }
}

export abstract class AnimatedSpriteContainer<T extends AnimatedSpriteModel> extends SpriteContainer<T> implements MonitorContrainer {
    constructor(
        options: AnimatedSpriteContainerProperty,
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: {[index:string]: {}}) {
        newInstanceOptions.offset = 0;
        super(
            { texture: PIXI.Texture.EMPTY, ...options }, 
            newInstance,
            { ...newInstanceOptions, animation: options.animation});
    }
}