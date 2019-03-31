import { PointModel, PIXIModel } from  "./geo";
import { Monitorable } from "../interfaces/monitor";
import { LocalableProprty } from "../interfaces/pixi";

const animationOpts = { round: 600 };

export class RailNode extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;
    protected tick: number

    constructor(options: LocalableProprty) {
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
    protected from: RailNode;
    protected to: RailNode;
    protected reverse: RailEdge;
    protected tick: number;

    constructor(options: LocalableProprty) {
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

        let from = this.toView(this.from.props.x, this.from.props.y)

        this.graphics.moveTo()
    }

    protected animate(delta: number) {
        this.tick += delta;
    }
}
