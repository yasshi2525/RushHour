import * as PIXI from "pixi.js";
import { generateOutline, generateFlash } from "./common";

const opts = {
    offset: 2,
    width: 1,
    alpha: 0.2,
    color: 0x607d8B,
    radius: 21
};

export default function(resolution: number, offset: number) {
    let graphics = new PIXI.Graphics();
    graphics.lineStyle(0, 0, 0);
    graphics.drawRect(0, 0, 50 * resolution, 50 * resolution);

    graphics.lineStyle(opts.width * resolution, opts.color);
    graphics.beginFill(opts.color, opts.alpha);
    let center = (opts.radius + opts.offset) * resolution;
    graphics.drawCircle(center, center, opts.radius * resolution);
    graphics.endFill();
    graphics.filters = [
        generateOutline(resolution), 
        generateFlash(resolution, offset)
    ];
    graphics.position = new PIXI.Point(opts.offset * resolution, opts.offset * resolution);
    return graphics;
}