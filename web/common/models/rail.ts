import * as PIXI from "pixi.js"
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { AnimatedSpriteProperty, ModelProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { GraphicsAnimationGenerator, GradientAnimationGenerator } from "./animate";
import { config } from "../interfaces/gamemap";

const graphicsOpts = {
    padding: 10,
    width: 4,
    maxWidth: 10,
    color: 0x9e9e9e,
    radius: 10,
    slide: 10
};

const rnDefaultValues: {
    pid: number,
    cid: number,
    color: number,
    mul: number
} = {
    pid: 0,
    cid: 0,
    color: 0,
    mul: 1
};

export class RailNode extends AnimatedSpriteModel implements Monitorable {
    parentRailNode: RailNode | undefined;
    protected edges: {[index: string]: RailEdge};

    constructor(options: AnimatedSpriteProperty) {
        super(options);
        this.parentRailNode = undefined;
        this.edges = {};
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(rnDefaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.sprite.tint = this.props.color;
        });
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("color", (color: number) => this.sprite.tint = color);
        this.addUpdateCallback("visible", () => {
            Object.keys(this.edges).forEach(eid => {
                let re = this.edges[eid];
                if (re.from !== undefined && re.to !== undefined) {
                    re.merge("visible", re.from.get("visible") && re.to.get("visible"));
                }
            });
        });
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            if (this.parentRailNode !== undefined) {
                this.parentRailNode.merge("visible", true);
            }
        })
    }

    beforeRender() {
        super.beforeRender();
        this.sprite.x -= graphicsOpts.padding / 2;
        this.sprite.y -= graphicsOpts.padding / 2;
    }

    resolve(parent: any | undefined) {
        if (parent !== undefined) {
            this.parentRailNode = parent;
            // 拡大時、派生元の座標から移動を開始する
            if (this.props.coord.zoom == 1) {
                this.current = Object.assign({}, parent.current)
                this.latency = config.latency;
            }
            // 縮小時、集約先の座標に向かって移動する
            if (this.props.coord.zoom == -1) {
                this.merge("pos", parent.get("pos"));
            }
            parent.merge("visible", false);
        }
    }
}

export class RailNodeContainer extends AnimatedSpriteContainer<RailNode> implements MonitorContrainer {
    constructor(options: ModelProperty) {
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(graphicsOpts.width, graphicsOpts.color);
        graphics.drawCircle(
            graphicsOpts.padding + graphicsOpts.radius, 
            graphicsOpts.padding + graphicsOpts.radius,
            graphicsOpts.radius);

        let generator = new GraphicsAnimationGenerator(options.app, graphics);
        
        let rect = graphics.getBounds().clone();
        rect.x -= graphicsOpts.padding - 1;
        rect.y -= graphicsOpts.padding - 1;
        rect.width += graphicsOpts.padding * 2;
        rect.height += graphicsOpts.padding * 2;

        let animation = generator.record(rect);
        super({ animation, ...options}, RailNode, {});
    }
}

const reDefaultValues: {
    from: number, 
    to: number, 
    eid: number
} = {
    from: 0, 
    to: 0, 
    eid: 0
};

export class RailEdge extends AnimatedSpriteModel implements Monitorable {
    from: RailNode|undefined;
    to: RailNode|undefined;
    protected reverse: RailEdge|undefined;
    protected theta: number = 0;

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(reDefaultValues);
    }

    resolve(from: any | undefined, to: any | undefined, reverse: any | undefined) {
        if (from !== undefined && to !== undefined) {
            if (this.from !== from && this.to !== to) {
                this.from = from;
                this.to = to;
                if (this.props.coord.zoom == -1) {
                    this.merge("visible", from.get("visible") && to.get("visible"))
                }
                this.sprite.tint = from.get("color");
                from.edges[this.props.id] = this;
                to.edges[this.props.id] = this;
                this.updateDestination();
                this.moveDestination();
            }
        }
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
    }

    beforeRender() {
        if (this.from !== undefined && this.to !== undefined 
            && this.from.current !== undefined && this.to.current !== undefined ) {
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
                x: avg.x + graphicsOpts.slide * Math.cos(theta + Math.PI / 2) - graphicsOpts.width,
                y: avg.y + graphicsOpts.slide * Math.sin(theta + Math.PI / 2) - graphicsOpts.width
            };

            this.sprite.rotation = theta;
            this.sprite.height = Math.min(this.props.mul, graphicsOpts.maxWidth);
            this.sprite.width = Math.sqrt(d.x * d.x + d.y * d.y);
        }
        super.beforeRender();
    }

    shouldEnd() {
        if (this.from !== undefined && this.to !== undefined) {
            // 縮小時、集約先に行き着くまで描画する
            if (this.props.coord.zoom == -1) {
                return this.from.shouldEnd() && this.to.shouldEnd();
            }
        }
        return this.props.outMap;
    }
}

export class RailEdgeContainer extends AnimatedSpriteContainer<RailEdge> implements MonitorContrainer {
    constructor(options: ModelProperty) {
        let generator = new GradientAnimationGenerator(options.app, graphicsOpts.color, 0.25);
        let animation =  generator.record();
        super({ animation, ...options}, RailEdge, {});
    }
}