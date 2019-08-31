import * as PIXI from "pixi.js";
import { generateOutline } from "./common";

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

const back = rgb2hsl(0x9e9e9e);

export default function(resolution: number, offset: number) {
    let canvas = document.createElement("canvas");
    canvas.width = 0xFF * resolution;
    canvas.height = resolution * 3;
    let ctx = canvas.getContext("2d") as CanvasRenderingContext2D;
    let grd = ctx.createLinearGradient(0, 0, 0xFF * resolution, 0);

    let grads = waveGradation(offset / 240, 0.25);

    grads.forEach(obj => {
        grd.addColorStop(obj.x, back.toBright(obj.v).toString());
    })

    ctx.fillStyle = grd;
    ctx.fillRect(0, 0, 0xFF * resolution, resolution * 3);

    let sprite = new PIXI.Sprite(PIXI.Texture.from(canvas));
    sprite.filters = [
        generateOutline(resolution)
    ];
    sprite.position = new PIXI.Point(resolution, resolution);
    let container = new PIXI.Container();
    let graphics = new PIXI.Graphics();
    graphics.lineStyle(0, 0, 0);
    graphics.drawRect(0, 0, (0xFF + 2) * resolution, 5 * resolution);
    container.addChild(graphics, sprite);
    return container;
} 

