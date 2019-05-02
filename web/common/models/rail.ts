import { PointModel, PIXIModel } from  "./geo";
import { Monitorable } from "../interfaces/monitor";
import { ApplicationProperty } from "../interfaces/pixi";

const animationOpts = { round: 600 };

export class RailNode extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;
    protected tick: number

    constructor(options: ApplicationProperty) {
        super(options);
        this.graphics = new PIXI.Graphics();
        this.tick = 0;
    }

    setupBeforeCallback() {
        super.setupBeforeCallback()
        this.addBeforeCallback(() => {
            this.container.addChild(this.graphics);
            this.app.ticker.add((d: number) => this.animate(d))
        })
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.clear();
        this.graphics.lineStyle(2, 0x4169e1);
        this.graphics.arc(this.vx, this.vy, 10, 0, Math.PI * 2)
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.graphics));
    }

    protected animate(delta: number) {
        this.tick += delta;
        let mod = this.tick % animationOpts.round;
        this.graphics.lineWidth = (mod < animationOpts.round / 2) ? mod : (animationOpts.round - mod);
        this.app.renderer.render(this.graphics);
    }
}

const reDefaultValues: {from: number, to: number, eid: number} = {from: 0, to: 0, eid: 0};

export class RailEdge extends PIXIModel implements Monitorable {
    protected graphics: PIXI.Graphics;
    protected from: RailNode|undefined;
    protected to: RailNode|undefined;
    protected reverse: RailEdge|undefined;
    protected tick: number;

    constructor(options: ApplicationProperty) {
        super(options);
        this.graphics = new PIXI.Graphics();
        this.tick = 0;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(reDefaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback()
        this.addBeforeCallback(() => {
            this.container.addChild(this.graphics);
            this.app.ticker.add((d: number) => this.animate(d))
        })
    }

    resolve(from: any | undefined, to: any | undefined, reverse: any | undefined) {
        if (from !== undefined) {
            this.from = from;
        }
        if (to !== undefined) {
            this.to = to;
        }
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.clear();
        this.graphics.lineStyle(2, 0x4169e1);

        if (this.from !== undefined && this.to !== undefined) {
            let from = this.toView(this.from.get("x"), this.from.get("y"))
            let to = this.toView(this.to.get("x"), this.to.get("y"))
            this.graphics.moveTo(from.x, from.y)
            this.graphics.lineTo(to.x, to.y)
        }
    }

    shouldEnd() {
        if (this.from !== undefined && this.to !== undefined) {
            return super.shouldEnd()
                && this.isOut(this.from.get("x"), this.from.get("y"))
                && this.isOut(this.to.get("x"), this.to.get("y"));
        } else {
            return super.shouldEnd();
        }
    }

    protected animate(delta: number) {
        this.tick += delta;
    }
}
