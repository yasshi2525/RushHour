import { AnimationProperty } from "../interfaces/pixi";
import { generateFlash, generateOutline, generateShadow } from "./filter";
import { config } from "../interfaces/gamemap";

function getCurveOffset(frame: number) {
    let ratio = frame / config.round;
    return Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
}

export abstract class AnimationGenerator {
    app: PIXI.Application;
    object: PIXI.DisplayObject;
    filters: AnimationProperty[];
    textures: PIXI.Texture[];
    
    constructor(app: PIXI.Application, object: PIXI.DisplayObject) {
        this.app = app;
        this.object = object;
        this.filters = [];
        this.textures = [];
    }

    protected applyFilter() {
        this.object.filters = this.filters.map(v => v.filter);
    }

    record(rect: PIXI.Rectangle) {
        for(let i = 0; i < config.round; i++) {
            let offset = getCurveOffset(i);
            this.filters.forEach(v => v.fn(v.filter, offset));
            this.textures.push(this.app.renderer.generateTexture(
                this.object, PIXI.SCALE_MODES.LINEAR, 
                this.app.renderer.resolution, rect));
        }
        return this.textures;
    }
}

export class GraphicsAnimationGenerator extends AnimationGenerator {
    constructor(app: PIXI.Application, obj: PIXI.Graphics) {
        super(app, obj);
        this.filters.push(generateOutline(app));
        this.filters.push(generateFlash(app));
        this.applyFilter();
    }
} 

export class ImageAnimationGenerator extends AnimationGenerator {
    constructor(app: PIXI.Application, texture: PIXI.Texture) {
        super(app, new PIXI.Sprite(texture));
        this.filters.push(generateOutline(app));
        this.filters.push(generateShadow(app));
        this.applyFilter();
    }
}