import { GlowFilter, OutlineFilter, DropShadowFilter, DropShadowFilterOptions } from "pixi-filters";

const round = 240;

function getCurveOffset(frame: number) {
    let ratio = frame / round;
    return Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
}

const flashOpts = { 
    width: { min: 0, max: 3 },
    distance: 6,
    outerStrength: 3, // 初期値
    innerStrength: 0,
    color: 0xaaaaaa,
    quality: 1
};


export function generateFlash(resolution: number, offset: number) {
    offset = getCurveOffset(offset);
    let filter = new GlowFilter(
        flashOpts.distance * resolution, 
        flashOpts.outerStrength * resolution,
        flashOpts.innerStrength * resolution,
        flashOpts.color,
        flashOpts.quality);
    filter.padding = (flashOpts.distance + flashOpts.outerStrength) * 2;
    filter.outerStrength = offset * flashOpts.width.min + (1 - offset) * flashOpts.width.max;
    return filter;
}

const outlineOpts = {
    width: 1,
    padding: 2,
    color: 0xffffff,
    quality: 1
};

export function generateOutline(resolution: number) {
    let filter = new OutlineFilter(
        outlineOpts.width * resolution, 
        outlineOpts.color, 
        outlineOpts.quality
    );
    filter.padding = outlineOpts.padding * resolution;
    return filter;
}

const shadowOpts: DropShadowFilterOptions = {
    color: 0x000000,
    distance: 5
}

export function generateShadow(resolution: number) {
    let filter = new DropShadowFilter(shadowOpts);
    filter.resolution = resolution;
    return filter;
}