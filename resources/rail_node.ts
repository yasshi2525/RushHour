import * as PIXI from "pixi.js";
import { generateOutline, generateFlash } from "./common";

const opts = {
    offset: 2,
    width: 4,
    color: 0x9e9e9e,
    radius: 11
};

export default function(resolution: number, offset: number) {
    let graphics = new PIXI.Graphics();
    graphics.lineStyle(0, 0, 0);
    graphics.drawRect(0, 0, 30 * resolution, 30 * resolution);

    graphics.lineStyle(opts.width * resolution, opts.color);
    let center = (opts.radius + opts.offset) * resolution;
    graphics.drawCircle(center, center, opts.radius * resolution);
    graphics.filters = [
        generateOutline(resolution), 
        generateFlash(resolution, offset)
    ];
    graphics.position = new PIXI.Point(opts.offset * resolution, opts.offset * resolution);
    return graphics;
}

