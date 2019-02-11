import BaseModel from "./base";
import { PointModel } from  "./geo";
import { Monitorable } from "../interfaces/monitor";

const animationOpts = { round: 600 };

export class RailNode extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;
    protected tick: number

    constructor(options: {app: PIXI.Application}) {
        super(options);
        this.graphics = new PIXI.Graphics();
        this.graphics.lineStyle(2, 0x4169e1);
        this.tick = 0;
    }

    setupBeforeCallback() {
        super.setupBeforeCallback()
        this.addBeforeCallback(() => {
            this.container.addChild(this.graphics);
            this.app.ticker.add((d: number) => this.animate(d))
        })
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("x", () => this.render());
        this.addUpdateCallback("y", () => this.render());
    }

    protected render() {
        console.log("render")
        this.graphics.clear()
        this.graphics.arc(this.props.x, this.props.y, 10, 0, Math.PI)
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

export class RailEdge extends BaseModel implements Monitorable {
    
}
