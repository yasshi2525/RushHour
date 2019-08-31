import * as PIXI from "pixi.js";
import { generateOutline, generateFlash } from "./common";

const opts = {
    offset: 4,
    width: 4,
    color: 0x607d8B,
    radius: 21,
    slice: 8
};

export default function(resolution: number, offset: number) {
    let graphics = new PIXI.Graphics();
    graphics.lineStyle(0, 0, 0);
    graphics.drawRect(0, 0, 50 * resolution, 50 * resolution);

    graphics.lineStyle(opts.width * resolution, opts.color);
    let center = (opts.radius + opts.offset) * resolution;

    for (var i = 0; i < opts.slice; i++) {
        let start = i / opts.slice * Math.PI * 2;
        let end = (i + 0.5) / opts.slice * Math.PI * 2;
        let next = (i + 1) / opts.slice * Math.PI * 2;

        graphics.lineStyle(opts.width * resolution, opts.color, 1);
        graphics.arc(center, center, opts.radius * resolution, start, end);
        graphics.lineStyle(opts.width * resolution, opts.color, 0);
        graphics.arc(center, center, opts.radius * resolution, end, next);
    }
    graphics.pivot = new PIXI.Point(center, center);
    graphics.rotation = offset / 240 * Math.PI * 2;
    graphics.filters = [
        generateOutline(resolution), 
        generateFlash(resolution, offset)
    ];
    graphics.position = new PIXI.Point(center, center);
    return graphics;
}

