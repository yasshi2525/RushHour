import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { ApplicationProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { GraphicsAnimationGenerator, GradientAnimationGenerator } from "./animate";

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
        let animation = generator.record(rect);
        super({ animation, ...options}, RailNode, {});
    }
}

const reDefaultValues: {from: number, to: number, eid: number} = {from: 0, to: 0, eid: 0};

export class RailEdge extends AnimatedSpriteModel implements Monitorable {
    protected from: RailNode|undefined;
    protected to: RailNode|undefined;
    protected reverse: RailEdge|undefined;
    protected theta: number = 0;

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(reDefaultValues);
    }

    resolve(from: any | undefined, to: any | undefined, reverse: any | undefined) {
        if (from !== undefined && to !== undefined) {
            if (this.from !== from || this.to !== to ) {
                this.from = from;
                this.to = to;
                this.updateDestination(true);
            }
        }
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
    }

    protected calcDestination() {
        if (this.from !== undefined && this.to !== undefined) {
            let theta = Math.atan2(
                this.to.destination.y - this.from.destination.y, 
                this.to.destination.x - this.from.destination.x);
            return {
                x: (this.from.destination.x + this.to.destination.x) / 2
                    + graphicsOpts.slide * Math.cos(theta + Math.PI / 2),
                y: (this.from.destination.y + this.to.destination.y) / 2
                    + graphicsOpts.slide * Math.sin(theta + Math.PI / 2)
            };
        }
        return {x: 0, y: 0};
    }

    protected smoothMove() {
        if (this.from !== undefined && this.to !== undefined) {
            let d = { 
                x: this.to.current.x - this.from.current.x,
                y: this.to.current.y - this.from.current.y
            };
            let avg = {
                x: (this.from.current.x + this.to.current.x) / 2,
                y: (this.from.current.y + this.to.current.y) / 2
            };
            let theta = Math.atan2(d.y, d.x);
            this.current = {
                x: avg.x + graphicsOpts.slide * Math.cos(theta + Math.PI / 2),
                y: avg.y + graphicsOpts.slide * Math.sin(theta + Math.PI / 2)
            };

            this.sprite.rotation = theta;
            this.sprite.height = graphicsOpts.width;
            this.sprite.width = Math.sqrt(d.x * d.x + d.y * d.y);
        }
        this.beforeRender();
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

export class RailEdgeContainer extends AnimatedSpriteContainer<RailEdge> implements MonitorContrainer {
    constructor(options: ApplicationProperty) {
        let generator = new GradientAnimationGenerator(options.app, graphicsOpts.color, 0.25);
        let animation =  generator.record();
        super({ animation, ...options}, RailEdge, {});
    }
}