import { Monitorable } from "../interfaces/monitor";
import { ApplicationProperty } from "../interfaces/pixi";
import GraphicsModel from  "./graphics";

const graphicsOpts = {
    width: 2,
    color: 0x4169e1,
    radius: 10
};

export class RailNode extends GraphicsModel implements Monitorable {

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.graphics.lineStyle(
                graphicsOpts.width, 
                graphicsOpts.color);
            this.graphics.arc(0, 0, graphicsOpts.radius, 0, Math.PI * 2)
        });
    }
}

const reDefaultValues: {from: number, to: number, eid: number} = {from: 0, to: 0, eid: 0};

export class RailEdge extends GraphicsModel implements Monitorable {
    protected from: RailNode|undefined;
    protected to: RailNode|undefined;
    protected reverse: RailEdge|undefined;

    constructor(options: ApplicationProperty) {
        super(options);
        this.graphics = new PIXI.Graphics();
        this.tick = 0;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(reDefaultValues);
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
        this.graphics.lineStyle(graphicsOpts.width / this.app.renderer.resolution, 
            graphicsOpts.color);

        if (this.from !== undefined && this.to !== undefined) {
            // 中心がcurrentなので、相対座標を求める
            this.graphics.moveTo(
                this.from.get("x") - this.props.x, 
                this.from.get("y") - this.props.y)
            this.graphics.lineTo(
                this.to.get("x") - this.props.x, 
                this.to.get("y") - this.props.y)
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
}
