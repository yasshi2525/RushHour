import { GraphicsModel, GraphicsContainer } from "./graphics";
import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import { ModelProperty } from "../interfaces/pixi";
import { config, Coordinates, Chunk, getChunk } from "../interfaces/gamemap";

const graphicsOpts = {
    world: 0xf44336,
    normal: 0x9e9e9e,
    width: 1
};

export class WorldBorder extends GraphicsModel implements Monitorable {
    protected radius: number;
    protected destRadius: number;

    constructor(props: ModelProperty) {
        super(props);
        this.radius = this.calcRadius(config.scale.default);
        this.destRadius = this.radius;
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.graphics.zIndex = -1;
    }

    updateDisplayInfo() {
        super.updateDisplayInfo();
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

enum BorderState {
    INVISIBLE,
    APPEAR,
    KEEP,
    DISAPPEAR,
}

const borderDefaultValues: {
    index: number,
    pos: number, 
    scale: number, 
    state: BorderState 
} = { 
    index: 0,
    pos: 0, 
    scale: config.scale.default, 
    state: BorderState.KEEP 
};

export class NormalBorder extends GraphicsModel implements Monitorable {
    protected v: boolean;
    count: number;

    constructor(props: ModelProperty & { v: boolean }) {
        super(props);
        this.v = props.v;
        this.count = 0;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(borderDefaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        if (props.v) {
            this.props.pos = { x: props.pos, y: 0 };
        } else {
            this.props.pos = { x: 0, y: props.pos} ;
        }
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(()=> {
            this.graphics.zIndex = -1;
            this.shape();
            this.updateDestination();
            this.moveDestination();
        });
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("resize", () => this.shape());
    }

    protected shape() {
        this.graphics.clear();
        this.graphics.lineStyle(graphicsOpts.width, graphicsOpts.normal);
        this.graphics.moveTo(0, 0);
        if (this.v) {
            this.graphics.lineTo(0, this.app.renderer.height);
        } else {
            this.graphics.lineTo(this.app.renderer.width, 0);
        }
    }

    updateDisplayInfo() {
        super.updateDisplayInfo();
        if (this.v) {
            this.graphics.y = 0;
        } else {
            this.graphics.x = 0;
        }
        switch (this.props.state) {
            case BorderState.INVISIBLE:
                this.graphics.visible = false;
                break;
            case BorderState.APPEAR:
                this.graphics.visible = true;
                this.graphics.alpha = this.count / config.latency * 2;
                break;
            case BorderState.DISAPPEAR:
                this.graphics.visible = true;
                this.graphics.alpha = 1 - this.count / config.latency * 2
                break;
            case BorderState.KEEP:
                this.graphics.visible = true;
                this.graphics.alpha = 1;
                break;
        }
    }
}

enum BorderContainerState {
    KEEP,
    DESTROY,
    GENERATE
}

const borderContainerDefaultValues: { coord: Coordinates, [index: string]: any } = {
    coord:  {
        cx: config.gamePos.default.x, 
        cy: config.gamePos.default.y, 
        scale: config.scale.default,
        zoom: 0
    },
    resize: false,
    delegate: 0
};

abstract class NormalBorderContainer extends GraphicsContainer<NormalBorder> implements MonitorContainer {
    protected chunk: Chunk;
    protected coord: Coordinates;
    protected v: boolean;
    protected state: BorderContainerState;
    protected zoom: boolean;
    protected count: number;
    protected suspends: boolean

    constructor(props: ModelProperty & { v: boolean, delegate: number }) {
        super(props, NormalBorder, { v: props.v });
        this.v = props.v;
        this.props.delegate = props.delegate;

        this.chunk = getChunk(
            config.gamePos.default, 
            config.scale.default - props.delegate + 1
        );
        this.coord = {
            cx: config.gamePos.default.x, 
            cy: config.gamePos.default.y,
            scale: config.scale.default,
            zoom: 0
        };
        this.state = BorderContainerState.KEEP;
        this.zoom = false;
        this.count = 0;
        this.suspends = false;
        this.genChildren(this.chunk);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(borderContainerDefaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        this.props.delegate = this.model.delegate;
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("delegate", () => {
            this.forEachChild(c => this.removeChild(c.get("id")));
            this.changeState(BorderContainerState.KEEP);
            this.chunk = getChunk({x: this.props.coord.cx, y: this.props.coord.cy}, this.props.coord.scale - this.props.delegate + 1);
            this.genChildren(this.chunk);
        });
        this.addUpdateCallback("coord", (v: Coordinates) => {
            let nowChunk = getChunk({x: v.cx, y: v.cy}, v.scale - this.props.delegate + 1);
            let nowOffset = this.getOffset(nowChunk);

            let num = Math.pow(2, this.props.delegate);

            if (this.chunk.scale !== nowChunk.scale) {
                this.zoom = nowChunk.scale < this.chunk.scale;

                if (this.state === BorderContainerState.KEEP) {
                    this.changeState(BorderContainerState.DESTROY);
                    // 拡大(縮小)による再作成
                    this.genChildren(nowChunk);
                } else {
                    this.suspends = true;
                }
            } else {
                // 左(上)側を作成、右(下)側を削除
                for (var i = 0; i < this.getOffset(this.chunk) - nowOffset; i++) {
                    var offset = i - Math.floor(num / 2) - 1;
                    this.genChild(this.chunk, offset);

                    var offset = i + Math.floor(num / 2) + this.getOffset(this.chunk);
                    this.removeChild(this.getId(offset, this.chunk.scale));
                }
                // 右(下)側を作成、左(上)側を削除
                for (var i = 0; i < nowOffset - this.getOffset(this.chunk); i++) {
                    var offset = i + Math.floor(num / 2) + 1;
                    this.genChild(this.chunk, offset);

                    var offset = i - Math.floor(num / 2) + this.getOffset(this.chunk);
                    this.removeChild(this.getId(offset, this.chunk.scale));
                }
            }
            
            this.chunk = nowChunk;
            this.coord = v;
        });
    }

    updateDisplayInfo() {
        super.updateDisplayInfo();
        if (this.suspends && this.state === BorderContainerState.KEEP) {
            this.changeState(BorderContainerState.DESTROY);
            // 拡大(縮小)による再作成
            this.genChildren(this.chunk);
            this.suspends = false;
        }

        switch(this.state) {
            case BorderContainerState.KEEP:
                break;
            case BorderContainerState.DESTROY:
                this.count++;
                if (this.count >= config.latency / 2) {
                    this.changeState(BorderContainerState.GENERATE);
                }
                break;
            case BorderContainerState.GENERATE:
                this.count++;
                if (this.count >= config.latency / 2) {
                    this.changeState(BorderContainerState.KEEP);
                }
                break;
        }
        this.forEachChild(c => {
            c.count = this.count;
            c.updateDisplayInfo();
        });
    }

    protected getInterval(chunk: Chunk) {
        return Math.pow(2, chunk.scale);
    }

    protected getId(offset: number, scale: number) {
        return offset + "_" + scale;
    }

    protected abstract getOffset(chunk: Chunk): number;

    protected abstract isAreaIn(offset: number): boolean;

    protected genChildOpts(id: string, index: number, pos: number, scale: number) {
        var state = BorderState.KEEP;
        switch (this.state) {
            case BorderContainerState.DESTROY:
                state = BorderState.INVISIBLE;
        }

        return {
            v: this.v,
            id: id,
            index: index,
            pos: pos,
            scale: scale,
            coord: this.coord,
            state: state
        }
    }

    protected genChild(chunk: Chunk, offset: number) {
        let index = this.getOffset(chunk) + offset;
        let id = this.getId(index, chunk.scale);
        let pos = index * this.getInterval(chunk);

        if (!this.isAreaIn(pos)) {
            return;
        }

        if (this.existsChild(id)) {
            this.removeChild(id);
        }

        this.addChild(this.genChildOpts(id, index, pos, chunk.scale));
    }

    protected genChildren(chunk: Chunk) {
        let num = Math.pow(2, this.props.delegate);
        for (var offset = -Math.floor(num / 2); offset < Math.floor(num / 2) + 1; offset++) {
            this.genChild(chunk, offset);
        }
    }

    protected changeState(state: BorderContainerState) {
        switch(state) {
            case BorderContainerState.DESTROY:
                if (!this.zoom) {
                    this.forEachChild(c => {
                        if (c.get("index") % 2 !== 0) {
                            c.merge("state", BorderState.DISAPPEAR);
                        }
                    });
                }
                break;
            case BorderContainerState.GENERATE: {
                this.forEachChild(c => {
                    if (c.get("scale") !== this.chunk.scale) {
                        this.removeChild(c.get("id"));
                    }
                });

                if (this.zoom) {
                    this.forEachChild(c => {
                        if (c.get("index") % 2 !== 0) {
                            c.merge("state", BorderState.APPEAR);
                        } else {
                            c.merge("state", BorderState.KEEP);
                        }
                    });
                } else {
                    this.forEachChild(c => c.merge("state", BorderState.KEEP));
                }

                break;
            }
            case BorderContainerState.KEEP:
                this.forEachChild(c => c.merge("state", BorderState.KEEP));
                break;
        }
        this.state = state;
        this.count = 0;
    }
}

export class XBorderContainer extends NormalBorderContainer implements MonitorContainer {
    constructor(props: ModelProperty & { delegate: number }) {
        super({ ...props, v: false });
    }

    protected getOffset(chunk: Chunk) {
        return chunk.y;
    }

    protected isAreaIn(offset: number) {
        return offset > config.gamePos.min.x && offset < config.gamePos.max.x;
    }
}

export class YBorderContainer extends NormalBorderContainer implements MonitorContainer {
    constructor(props: ModelProperty & { delegate: number }) {
        super({ ...props, v: true });
    }

    protected getOffset(chunk: Chunk) {
        return chunk.x;
    }

    protected isAreaIn(offset: number) {
        return offset > config.gamePos.min.y && offset < config.gamePos.max.y;
    }
}