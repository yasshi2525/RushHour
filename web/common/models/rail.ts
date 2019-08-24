import * as PIXI from "pixi.js"
import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import { AnimatedSpriteProperty, PIXIProperty, cursorOpts } from "../interfaces/pixi";
import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { GraphicsAnimationGenerator, GradientAnimationGenerator } from "./animate";
import { config, ResolveError, Point } from "../interfaces/gamemap";
import { MenuStatus } from "@/state";

const graphicsOpts = {
    padding: 10,
    offset: 1,
    width: 4,
    maxWidth: 10,
    color: 0x9e9e9e,
    radius: 10,
    slide: 10
};

const rnDefaultValues: {
    pid: number,
    cid: number,
    mul: number
} = {
    pid: 0,
    cid: 0,
    mul: 1
};

export class RailNode extends AnimatedSpriteModel implements Monitorable {
    parentRailNode: RailNode | undefined;
    edges: {[index: string]: RailEdge};

    constructor(options: AnimatedSpriteProperty) {
        super(options);
        this.parentRailNode = undefined;
        this.edges = {};
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(rnDefaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
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

    updateDisplayInfo() {
        super.updateDisplayInfo();
        this.sprite.x -= graphicsOpts.offset;
        this.sprite.y -= graphicsOpts.offset;
    }

    resolve(error: ResolveError) {
        let owner = this.resolveOwner(this.props.oid);
        if (owner !== undefined) {
            this.merge("tint", owner.get("color"));
        }
        let parent = this.model.gamemap.get("rail_nodes", this.props.pid) as RailNode | undefined;
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
        let hasUnresolvedOwner = error.hasUnresolvedOwner || owner === undefined
        error.hasUnresolvedOwner = hasUnresolvedOwner;
        return error;
    }
}

export class RailNodeContainer extends AnimatedSpriteContainer<RailNode> implements MonitorContainer {
    cursor: RailNode;

    constructor(options: PIXIProperty) {
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
        this.cursor = this.addChild(cursorOpts);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("menu", v => {
            switch(v) {
                case MenuStatus.IDLE:
                    this.cursor.merge("visible", false);
                    break;
            }

        });
        this.addUpdateCallback("cursorClient", (v: Point | undefined) => {
            switch(this.props.menu) {
                case MenuStatus.SEEK_DEPARTURE:
                case MenuStatus.EXTEND_RAIL:
                    if (v !== undefined && this.model.controllers.getCursor().selected === undefined) {
                        this.cursor.merge("visible", true);
                        this.cursor.destination = v;
                        this.cursor.moveDestination();
                    } else {
                        this.cursor.merge("visible", false);
                    }
                    break;
                default:
                    this.cursor.merge("visible", false);
            }
        });
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

    resolve(error: ResolveError) {
        let from = this.model.gamemap.get("rail_nodes", this.props.from) as RailNode | undefined;
        let to = this.model.gamemap.get("rail_nodes", this.props.to) as RailNode | undefined;
        let reverse = this.model.gamemap.get("rail_edges", this.props.eid) as RailEdge | undefined;

        if (from !== undefined && to !== undefined) {
            if (this.from !== from && this.to !== to) {
                this.from = from;
                this.to = to;
                if (this.props.coord.zoom == -1) {
                    this.merge("visible", from.get("visible") && to.get("visible"))
                }
                this.sprite.tint = from.get("tint");
                from.edges[this.props.id] = this;
                to.edges[this.props.id] = this;
                this.updateDestination();
                this.moveDestination();
            }
        }
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
        return error;
    }

    updateDisplayInfo() {
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
            this.sprite.height = graphicsOpts.width;
            this.sprite.width = Math.sqrt(d.x * d.x + d.y * d.y);
        }
        super.updateDisplayInfo();
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

export class RailEdgeContainer extends AnimatedSpriteContainer<RailEdge> implements MonitorContainer {
    cursorFrom: RailEdge;
    cursorTo: RailEdge;

    constructor(options: PIXIProperty) {
        let generator = new GradientAnimationGenerator(options.app, graphicsOpts.color, 0.25);
        let animation =  generator.record();
        super({ animation, ...options}, RailEdge, {});
        this.cursorFrom = this.addChild({ ...cursorOpts, id: "cursorFrom" });
        this.cursorTo = this.addChild({ ...cursorOpts, id: "cursorTo" });
    }
}