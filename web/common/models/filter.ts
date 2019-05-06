import { GlowFilter, OutlineFilter, DropShadowFilter } from "pixi-filters";
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
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: (filter: GlowFilter, offset: number) => {
        filter.outerStrength = offset * flashOpts.width.min 
                                + (1 - offset) * flashOpts.width.max;
    }};
}

const outlineOpts = {
    width: 1,
    color: 0xffffff
};

export function generateOutline(app: PIXI.Application) {
    let filter = new OutlineFilter(
        outlineOpts.width, outlineOpts.color
    );
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: () => {}};
}

const shadowOpts: PIXI.filters.DropShadowFilterOptions = {
    color: 0x000000,
    distance: 5
}

export function generateShadow(app: PIXI.Application) {
    let filter = new DropShadowFilter(shadowOpts);
    filter.resolution = app.renderer.resolution;
    return { app, filter, fn: () => {}}
}