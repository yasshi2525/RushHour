import { GlowFilter, OutlineFilter, DropShadowFilter, DropShadowFilterOptions } from "pixi-filters";
import { AnimationProperty } from "../interfaces/pixi";

const flashOpts = { 
    width: { min: 0, max: 3 },
    distance: 6,
    outerStrength: 3, // 初期値
    innerStrength: 0,
    color: 0xaaaaaa,
    quality: 1
};


export function generateFlash(app: PIXI.Application): AnimationProperty {
    let filter = new GlowFilter(
        flashOpts.distance, 
        flashOpts.outerStrength,
        flashOpts.innerStrength,
        flashOpts.color,
        flashOpts.quality);
    filter.padding = flashOpts.distance + flashOpts.outerStrength * 2;
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: (filter: GlowFilter, offset: number) => {
        filter.outerStrength = offset * flashOpts.width.min 
                                + (1 - offset) * flashOpts.width.max;
    }};
}

const outlineOpts = {
    width: 1,
    color: 0xffffff,
    quality: 1
};

export function generateOutline(app: PIXI.Application) {
    let filter = new OutlineFilter(
        outlineOpts.width, outlineOpts.color, outlineOpts.quality
    );
    filter.padding = outlineOpts.width;
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: () => {}};
}

const shadowOpts: DropShadowFilterOptions = {
    color: 0x000000,
    distance: 5
}

export function generateShadow(app: PIXI.Application) {
    let filter = new DropShadowFilter(shadowOpts);
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: () => {}}
}