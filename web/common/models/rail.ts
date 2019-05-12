import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { GraphicsModel, GraphicsContainer } from  "./graphics";
import { ApplicationProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { GraphicsAnimationGenerator } from "./animate";

const graphicsOpts = {
    width: 4,
    color: 0x4169e1,
    radius: 10,
    slide: 10
};

export class RailNode extends AnimatedSpriteModel implements Monitorable {
}

export class RailNodeContainer extends AnimatedSpriteContainer<RailNode> implements MonitorContrainer {
    constructor(options: ApplicationProperty) {
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(graphicsOpts.width, graphicsOpts.color);
        graphics.arc(0, 0, graphicsOpts.radius, 0, Math.PI * 2);

        let generator = new GraphicsAnimationGenerator(options.app, graphics);

        let rect = graphics.getBounds().clone();
        rect.x -= 3 * options.app.renderer.resolution;
        rect.y -= 3 * options.app.renderer.resolution;
        rect.width += 6 * options.app.renderer.resolution;
        rect.height += 6 * options.app.renderer.resolution;
        let animation =  generator.record(rect);
        super({ animation, ...options}, RailNode, {});
    }
}

const reDefaultValues: {from: number, to: number, eid: number} = {from: 0, to: 0, eid: 0};

export class RailEdge extends GraphicsModel implements Monitorable {
    protected from: RailNode|undefined;
    protected to: RailNode|undefined;
    protected reverse: RailEdge|undefined;

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(reDefaultValues);
    }

    resolve(from: any | undefined, to: any | undefined, reverse: any | undefined) {
        if (from !== undefined && to !== undefined) {
            this.from = from;
            this.to = to;
            this.props.x = (from.get("x") + to.get("x")) / 2;
            this.props.y = (from.get("y") + to.get("y")) / 2;
        }
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.clear();
        this.graphics.lineStyle(graphicsOpts.width, graphicsOpts.color);

        if (this.from !== undefined && this.to !== undefined) {
            // 中心がcurrentなので、相対座標を求める
            let from = this.from.current;
            let to = this.to.current;

            var theta = Math.atan2(to.y - from.y, to.x - from.x) - Math.PI / 2;

            this.graphics.moveTo(
                from.x + Math.cos(theta) * graphicsOpts.slide - this.current.x, 
                from.y + Math.sin(theta) * graphicsOpts.slide - this.current.y);
            this.graphics.lineTo(
                to.x + Math.cos(theta) * graphicsOpts.slide - this.current.x, 
                to.y + Math.sin(theta) * graphicsOpts.slide - this.current.y);
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

export class RailEdgeContainer extends GraphicsContainer<RailEdge> implements MonitorContrainer {
    constructor(options: ApplicationProperty) {
        super(options, RailEdge, {});
    }
}