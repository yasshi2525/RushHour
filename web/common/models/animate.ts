import { AnimationProperty } from "../interfaces/pixi";
import { generateFlash, generateOutline, generateShadow } from "./filter";
import { config } from "../interfaces/gamemap";

/**
 * cosカーブで値を返す
 * @param frame フレーム番号
 */
function getCurveOffset(frame: number) {
    let ratio = frame / config.round;
    return Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
}

class HSL { 
    h = 0;
    s = 0;
    l = 0;

    constructor(h: number, s: number, l: number) {
        this.h = h;
        this.s = s;
        this.l = l;
    }

    /**
     * 明るくする。0指定で現状、1指定で白
     * @param v 明度アップ率
     */
    toBright(v: number) {
        return new HSL(this.h, this.s, Math.round(this.l + (100 - this.l) * v));
    }

    toString() {
        return "hsl(" + this.h + "," + this.s + "%," + this.l + "%)";
    }
};

function rgb2hsl(color: number): HSL {
    let r = (color >> 16) & 0xFF;
    let g = (color >> 8) & 0xFF;
    let b = color & 0xFF;
    let max = Math.max(r, g, b);
    let min = Math.min(r, g, b);
    let hsl = new HSL(0, 0, (max + min) / 2);
    
    if (max != min) {
        switch(max) {
            case r:
                hsl.h = 60 * (g - b) / (max - min);
                break;
            case g:
                hsl.h = 60 * (b - r) / (max - min) + 120;
                break;
            case b:
                hsl.h = 60 * (r - g) / (max - min) + 240;
                break;
        } 
    }

    if (hsl.l <= 127) {
        hsl.s = (max - min) / (max + min);
    } else {
        hsl.s = (max - min) / (510 - max - min);
    }

    if (hsl.h < 0) {
        hsl.h += 360;
    }

    hsl.h = Math.round(hsl.h);
    hsl.s = Math.round(hsl.s * 100);
    hsl.l = Math.round(hsl.l / 0xFF * 100);
    return hsl;
}

abstract class Generator {
    app: PIXI.Application;
    textures: PIXI.Texture[];

    constructor(app: PIXI.Application) {
        this.app = app;
        this.textures = [];
    }
}

/**
 * グラディエーションなしのアニメーション生成器
 */
abstract class MonoGenerator extends Generator {
    object: PIXI.DisplayObject;
    filters: AnimationProperty[];
    
    constructor(app: PIXI.Application, object: PIXI.DisplayObject) {
        super(app);
        this.object = object;
        this.filters = [];
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

function waveGradation(p: number, r: number) {
    let grads = []

    if (p < r) {
        grads.push({x: 0,         v: (r - p) / r});
        grads.push({x: p,         v: 1});
        grads.push({x: p + r,     v: 0});
        grads.push({x: p - r + 1, v: 0});
        grads.push({x: 1,         v: (r - p) / r});
    } else if (p < 1 - r) {
        grads.push({x: 0,     v: 0});
        grads.push({x: p - r, v: 0});
        grads.push({x: p,     v: 1});
        grads.push({x: p + r, v: 0});
        grads.push({x: 1,     v: 0});
    } else if (p <= 1) {
        grads.push({x: 0,         v: (p - 1 + r) / r});
        grads.push({x: p - 1 + r, v: 0});
        grads.push({x: p - r,     v: 0});
        grads.push({x: p,         v: 1});
        grads.push({x: 1,         v: (p - 1 + r) / r});
    } else {
        grads.push({x: 0, v: 0});
        grads.push({x: 1, v: 0});
    }
    return grads;
}


export class GradientAnimationGenerator extends Generator {
    quality: number;
    back: HSL;
    front: HSL;
    width: number;

    constructor(app: PIXI.Application, color: number, width: number) {
        super(app);
        this.quality = 0xFF;
        this.back = rgb2hsl(color);
        this.front = new HSL(this.back.h, this.back.s, 100);
        this.width = width;
    }

    record() {
        for(let i = 0; i < config.round; i++) {
            let x = i / config.round;
            let radius = this.width / 2;
            
            let canvas = document.createElement("canvas");
            canvas.width = this.quality;
            canvas.height = 1;
            let ctx = canvas.getContext("2d");
    
            if (ctx != null) {
                let grd = ctx.createLinearGradient(0, 0, this.quality, 0);

                let grads = waveGradation(x, radius);

                grads.forEach(obj => {
                    grd.addColorStop(obj.x, this.back.toBright(obj.v).toString());
                })
                
                ctx.fillStyle = grd;
                ctx.fillRect(0, 0, this.quality, 1);

                this.textures.push(PIXI.Texture.from(canvas));
            }
        }

        return this.textures;
    }
}

export class GraphicsAnimationGenerator extends MonoGenerator {
    constructor(app: PIXI.Application, obj: PIXI.Graphics) {
        super(app, obj);
        this.filters.push(generateOutline(app));
        this.filters.push(generateFlash(app));
        this.applyFilter();
    }
}


export class ImageAnimationGenerator extends MonoGenerator {
    constructor(app: PIXI.Application, texture: PIXI.Texture) {
        super(app, new PIXI.Sprite(texture));
        this.filters.push(generateOutline(app));
        this.filters.push(generateShadow(app));
        this.applyFilter();
    }
}