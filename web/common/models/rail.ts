import { MenuStatus } from "../../state";
import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import { AnimatedSpriteProperty, cursorOpts, AnimatedSpriteContainerProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { config, ResolveError, Point } from "../interfaces/gamemap";
import { PointModel } from "./point";

const graphicsOpts = {
    width: 4,
    color: 0x9e9e9e,
    slide: 12
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
    out: {[index: string]: RailEdge};
    in: {[index: string]: RailEdge};

    constructor(options: AnimatedSpriteProperty) {
        super(options);
        this.parentRailNode = undefined;
        this.out = {};
        this.in = {};
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(rnDefaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("visible", () => {
            Object.keys(this.out).forEach(eid => this.out[eid].updateVisible());
            Object.keys(this.in).forEach(eid => this.in[eid].updateVisible());
        });
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            if (this.parentRailNode !== undefined) {
                this.parentRailNode.merge("visible", true);
            }
            Object.keys(this.out).forEach(eid => this.out[eid].merge("from", undefined));
            Object.keys(this.in).forEach(eid => this.in[eid].merge("to", undefined));
        })
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
    cursor: RailNode | undefined;

    constructor(options: AnimatedSpriteContainerProperty) {
        super(options, RailNode, {});
        if (!this.model.isReadOnly()) {
            this.cursor = this.addChild({ oid: this.model.myid, ...cursorOpts });
        }
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        if (this.cursor !== undefined) {
            this.cursor.current = undefined;
            this.cursor.destination = undefined;
        }
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        if (this.cursor === undefined) {
            return
        }
        this.addUpdateCallback("menu", v => {
            if (this.cursor === undefined) {
                return
            }
            switch(v) {
                case MenuStatus.IDLE:
                    this.cursor.merge("visible", false);
                    break;
                case MenuStatus.SEEK_DEPARTURE:
                case MenuStatus.EXTEND_RAIL:
                    this.cursor.merge("visible", true);
                    break;
            }
        });
        this.addUpdateCallback("cursorClient", (v: Point | undefined) => {
            if (this.cursor === undefined) {
                return
            }
            switch(this.props.menu) {
                case MenuStatus.SEEK_DEPARTURE:
                case MenuStatus.EXTEND_RAIL:
                    this.cursor.merge("visible", v !== undefined);
                    this.cursor.destination = v;
                    this.cursor.moveDestination();
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

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("from", (fromID: string | undefined) => {
            let from = this.model.gamemap.get("rail_nodes", fromID) as RailNode | undefined;
            this.resolveFrom(from);
        });
        this.addUpdateCallback("to", (toID: string | undefined) => {
            let to = this.model.gamemap.get("rail_nodes", toID) as RailNode | undefined;
            this.resolveTo(to);
        });
    }

    resolve(error: ResolveError) {
        let from = this.model.gamemap.get("rail_nodes", this.props.from) as RailNode | undefined;
        let to = this.model.gamemap.get("rail_nodes", this.props.to) as RailNode | undefined;

        this.resolveFrom(from);
        this.resolveTo(to);
        
        let reverse = this.model.gamemap.get("rail_edges", this.props.eid) as RailEdge | undefined;
        if (reverse !== undefined) {
            this.reverse = reverse;
        }
        return error;
    }

    protected resolveFrom(from: RailNode | undefined) {
        if (this.from !== from) {
            this.unlinkFrom();
        }
        this.linkFrom(from);
        this.updateVisible();
    }

    protected resolveTo(to: RailNode | undefined) {
        if (this.to !== to) {
            this.unlinkTo();
        }
        this.linkTo(to);
        this.updateVisible();
    }

    protected unlinkFrom() {
        if (this.from !== undefined) {
            delete this.from.out[this.props.id];
        }
        this.from = undefined;
    }

    protected unlinkTo() {
        if (this.to !== undefined) {
            delete this.to.in[this.props.id];
        }
        this.to = undefined;
    }

    protected linkFrom(from: RailNode | undefined) {
        if (this.from !== from) {
            this.from = from;
            if (from !== undefined) {
                from.out[this.props.id] = this;
                this.sprite.tint = from.get("tint");
            }
        }
    }

    protected linkTo(to: RailNode | undefined) {
        if (this.to !== to) {
            this.to = to;
            if (to !== undefined) {
                to.in[this.props.id] = this;
                this.sprite.tint = to.get("tint");
            }
        }
    }

    updateVisible() {
        if (this.from !== undefined && this.to !== undefined) {
            this.merge("visible", this.from.get("visible") && this.to.get("visible"))
        } else {
            this.merge("visible", false);
        }
        this.updateDisplayInfo();
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
                x: avg.x + graphicsOpts.slide * Math.cos(theta + Math.PI / 2),
                y: avg.y + graphicsOpts.slide * Math.sin(theta + Math.PI / 2)
            };

            this.sprite.rotation = theta;
            this.sprite.width = Math.sqrt(d.x * d.x + d.y * d.y);
            this.sprite.visible = true;
        } else {
            this.sprite.visible = false;
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
    deamon: RailNode | undefined;
    cursorOut: RailEdge | undefined;
    cursorIn: RailEdge | undefined;

    constructor(options: AnimatedSpriteContainerProperty) {
        super(options, RailEdge, {});
        if (!this.model.isReadOnly()) {
            this.cursorOut = this.addChild({ ...cursorOpts, oid: this.model.myid, id: "cursorOut", from: "cursor", reverse: "cursorIn" });
            this.cursorIn = this.addChild({ ...cursorOpts, oid: this.model.myid, id: "cursorIn", to: "cursor", reverse: "cursorOut" });
            this.cursorOut.resolve({});
            this.cursorIn.resolve({});
            this.deamon = this.cursorOut.from as RailNode
        }
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        if (this.model.isReadOnly()) {
            return
        }

        this.addUpdateCallback("anchorObj", (v: PointModel | undefined) => {
            if (v instanceof RailNode || v === undefined) {
                this.setAnchor(v)
            }
        });
        this.addUpdateCallback("cursorObj", (v: PointModel | undefined) => {
            if (v instanceof RailNode) {
                this.setCursor(v);
            } else {
                this.setCursor(this.deamon);
            }
        });
    }

    protected setAnchor(anchor: RailNode | undefined) {
        let id = anchor !== undefined ? anchor.get("id") : undefined;
        if (this.cursorOut !== undefined && this.cursorIn !== undefined) {
            this.cursorOut.merge("to", id);
            this.cursorIn.merge("from", id);
        }
    }

    protected setCursor(cursor: RailNode | undefined) {
        let id = cursor !== undefined ? cursor.get("id") : undefined;
        if (this.cursorOut !== undefined && this.cursorIn !== undefined) {
            this.cursorOut.merge("from", id);
            this.cursorIn.merge("to", id);
        }
    }
}