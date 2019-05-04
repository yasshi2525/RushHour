import * as Filters from "pixi-filters";
import { ApplicationProperty } from "../interfaces/pixi";
import { Monitorable } from "../interfaces/monitor";
import PointModel from "./point"

const animationOpts = { 
    width: { min: 2, max: 4 },
    round: 5000,
    distance: 8,
    outerStrength: 2, // 初期値
    innerStrength: 0,
    color: 0xcccccc,
    quality: 0.5 
};

export default abstract class extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;
    /**
     * インスタンス作成からの累計時間
     */
    protected tick: number;

    /**
     * 明滅率(0-1)
     */
    protected offset: number;

    protected glow: PIXI.filters.GlowFilter;

    constructor(options: ApplicationProperty) {
        super(options);
        this.graphics = new PIXI.Graphics();
        this.glow = new Filters.GlowFilter(
            animationOpts.distance, 
            animationOpts.outerStrength * this.app.renderer.resolution,
            animationOpts.innerStrength,
            animationOpts.color,
            animationOpts.quality
        );
        this.glow.resolution = this.app.renderer.resolution;
        this.graphics.filters = [this.glow];
        this.tick = 0;
        this.offset = 0;
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.container.addChild(this.graphics);
            this.app.ticker.add(() => this.flash())
        });
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.graphics));
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.x = this.current.x;
        this.graphics.y = this.current.y;
        this.glow.outerStrength = (this.offset * animationOpts.width.min 
                                + (1 - this.offset) * animationOpts.width.max);
    }

    protected flash() {
        this.tick += this.app.ticker.elapsedMS;
        let ratio = (this.tick % animationOpts.round) / animationOpts.round;
        this.offset = Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
        this.beforeRender();
    }
}