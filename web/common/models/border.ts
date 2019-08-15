import { GraphicsModel } from "./graphics";
import { Monitorable } from "../interfaces/monitor";
import { ApplicationProperty } from "../interfaces/pixi";
import { config } from "../interfaces/gamemap";

const graphicsOpts = {
    world: 0xf44336,
    normal: 0x9e9e9e,
    width: 1
};

export class WorldBorder extends GraphicsModel implements Monitorable {
    protected radius: number;
    protected destRadius: number;

    constructor(props: ApplicationProperty) {
        super(props);
        this.radius = this.calcRadius(config.scale.default);
        this.destRadius = this.radius;
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.graphics.zIndex = -1;
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.clear();
        this.graphics.lineStyle(graphicsOpts.width, graphicsOpts.world);
        this.graphics.drawRect(-this.radius/2, -this.radius/2, this.radius, this.radius);
    }

    updateDestination() {
        super.updateDestination();
        this.destRadius = this.calcRadius(this.props.coord.scale);
    }

    moveDestination() {
        super.moveDestination();
        this.radius = this.destRadius;
    }

    protected smoothMove() {
        super.smoothMove()
        if (this.latency > 0) {
            let ratio = this.latency / config.latency;
            if (ratio < 0.5) {
                ratio = 1.0 - ratio;
            }
            this.radius = this.radius * ratio + this.destRadius * (1 - ratio);
        } else {
            this.radius = this.destRadius;
        }
    }

    protected calcRadius(scale: number) {
        return Math.pow(2, config.scale.max - scale) * Math.max(this.app.renderer.width, this.app.renderer.height); 
    }
}