import { GraphicsModel, GraphicsContainer } from "./graphics";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { ApplicationProperty } from "../interfaces/pixi";
import { config, Coordinates, Point } from "../interfaces/gamemap";

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

const borderDefaultValues: { offset: number } = {offset: 0};

export class NormalBorder extends GraphicsModel implements Monitorable {
    protected v: boolean;

    constructor(props: ApplicationProperty & {v: boolean}) {
        super(props);
        this.v = props.v;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(borderDefaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        if (props.v) {
            this.props.x = props.offset;
        } else {
            this.props.y = props.offset;
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

    beforeRender() {
        super.beforeRender();
        if (this.v) {
            this.graphics.y = 0;
        } else {
            this.graphics.x = 0;
        }
    }
}

interface Chunk {
    chx: number,
    chy: number,
    scale: number
}

const borderContainerDefaultValues: { coord: Coordinates, [index: string]: any } = {
    coord:  {
        cx: config.gamePos.default.x, 
        cy: config.gamePos.default.y, 
        scale: config.scale.default,
        zoom: 0
    },
    resize: false,
    forceMove: false,
    outMap: false,
    visible: true
};

abstract class NormalBorderContainer extends GraphicsContainer<NormalBorder> implements MonitorContrainer {
    protected chunk: Chunk;
    protected v: boolean;

    constructor(props: ApplicationProperty & { v: boolean }) {
        super(props, NormalBorder, {v: props.v});
        this.v = props.v;

        this.chunk = this.getChunk(config.gamePos.default, config.scale.default);
        this.genChildren(this.chunk);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(borderContainerDefaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("coord", (v: Coordinates) => {
            let nowChunk = this.getChunk({x: v.cx, y: v.cy }, v.scale);
            let nowOffset = this.getOffset(nowChunk);

            let num = Math.pow(2, config.scale.delegate);

            if (this.chunk.scale !== nowChunk.scale) {
                // 拡大(縮小)による再作成
                this.forEachChild(c => this.removeChild(c.get("id")));
                this.genChildren(nowChunk);
            } else {
                // 左(上)側を作成、右(下)側を削除
                for (var i = 0; i < this.getOffset(this.chunk) - nowOffset; i++) {
                    var offset = i - Math.floor(num / 2) - 1;
                    this.addChild(this.genChildOpts(this.chunk, offset));

                    var offset = i + Math.floor(num / 2) + this.getOffset(this.chunk);
                    this.removeChild(this.getId(offset, this.chunk.scale));
                }
                // 右(下)側を作成、左(上)側を削除
                for (var i = 0; i < nowOffset - this.getOffset(this.chunk); i++) {
                    var offset = i + Math.floor(num / 2) + 1;
                    this.addChild(this.genChildOpts(this.chunk, offset));

                    var offset = i - Math.floor(num / 2) + this.getOffset(this.chunk);
                    this.removeChild(this.getId(offset, this.chunk.scale));
                }
            }
            this.chunk = nowChunk;
        });
    }

    protected getInterval(chunk: Chunk) {
        return Math.pow(2, chunk.scale - config.scale.delegate + 1);
    }

    protected getChunk(pos: Point, scale: number): Chunk {
        scale = Math.floor(scale);
        let interval = Math.pow(2, scale - config.scale.delegate + 1);
        return {
            chx: Math.floor(pos.x / interval),
            chy: Math.floor(pos.y / interval),
            scale: scale
        };
    }

    protected getId(offset: number, scale: number) {
        return offset + "_" + scale;
    }

    protected abstract getOffset(chunk: Chunk): number;

    protected abstract isAreaIn(offset: number): boolean;

    protected genChildOpts(chunk: Chunk, offset: number) {
        let coord: Coordinates = (this.props.coord !== undefined) ? this.props.coord : {
            cx: config.gamePos.default.x, 
            cy: config.gamePos.default.y,
            scale: config.scale.default,
            zoom: 0
        }
        return {
            id: this.getId(this.getOffset(chunk) + offset, chunk.scale),
            offset: (this.getOffset(chunk) + offset) * this.getInterval(chunk),
            v: this.v,
            coord: coord
        }
    }

    protected genChildren(chunk: Chunk) {
        let num = Math.pow(2, config.scale.delegate);
        for (var offset = -Math.floor(num / 2); offset < Math.floor(num / 2) + 1; offset++) {
            let opt = this.genChildOpts(chunk, offset);
            if (this.isAreaIn(opt.offset)) {
                this.addChild(opt);
            }
        }
    }
}

export class XBorderContainer extends NormalBorderContainer implements MonitorContrainer {
    constructor(props: ApplicationProperty) {
        super({ ...props, v: false });
    }

    protected getOffset(chunk: Chunk) {
        return chunk.chy;
    }

    protected isAreaIn(offset: number) {
        return offset > config.gamePos.min.x && offset < config.gamePos.max.x;
    }
}

export class YBorderContainer extends NormalBorderContainer implements MonitorContrainer {
    constructor(props: ApplicationProperty) {
        super({ ...props, v: true });
    }

    protected getOffset(chunk: Chunk) {
        return chunk.chx;
    }

    protected isAreaIn(offset: number) {
        return offset > config.gamePos.min.y && offset < config.gamePos.max.y;
    }
}