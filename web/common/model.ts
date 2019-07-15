import * as PIXI from "pixi.js";
import { ResidenceContainer, CompanyContainer } from "./models/background";
import MonitorContainer from "./models/container";
import { Coordinates, config } from "./interfaces/gamemap";
import { Monitorable } from "./interfaces/monitor";
import { GameModelProperty, ResourceAttachable } from "./interfaces/pixi";
import { GameMap } from "../state";
import { RailEdge, RailNodeContainer, RailEdgeContainer } from "./models/rail";

const forceMove = { forceMove: true };

export default class implements ResourceAttachable {
    protected app: PIXI.Application;
    renderer: PIXI.Renderer;
    protected payload: {[index:string]: MonitorContainer<Monitorable>} = {}
    protected changed: boolean = false;
    timestamp: number;
    textures: {[index: string]: PIXI.Texture};
    coord: Coordinates;
    offset: number;
    /**
     * サーバから全データの取得が必要か
     */
    shouldFetch: boolean;
    shouldRemoveOutsider: boolean;
    debugText: PIXI.Text;
    debugValue: any;

    constructor(options: GameModelProperty) {
        this.app = options.app;
        this.renderer = options.app.renderer;
        this.textures = {};
       
        this.coord = { cx: options.cx, cy: options.cy, scale: options.scale }
        this.shouldFetch = false;
        this.shouldRemoveOutsider = false;
        this.timestamp = 0;
        this.offset = 0;

        this.app.ticker.add(() => {
            this.offset++;
            if (this.offset >= config.round) {
                this.offset = 0;
            }
            Object.keys(this.payload).forEach(key => this.payload[key].merge("offset", this.offset));
        });

        this.debugText = new PIXI.Text("");
        this.debugText.style.fontSize = 14;
        this.debugText.style.fill = 0xffffff;
        this.app.stage.addChild(this.debugText);
        setInterval(() => this.viewDebugInfo(), 250);
    }

    attach(textures: {[index: string]: PIXI.Texture}) {
        this.payload["residences"] = new ResidenceContainer({ app: this.app, texture: textures.residence});
        this.payload["companies"] = new CompanyContainer({ app: this.app, texture: textures.company});
        this.payload["rail_nodes"] = new RailNodeContainer({ app: this.app});
        this.payload["rail_edges"] = new RailEdgeContainer({ app: this.app});

        Object.keys(this.payload).forEach(key => {
            this.payload[key].setupDefaultValues();
            this.payload[key].setupUpdateCallback();
            this.payload[key].setupBeforeCallback();
            this.payload[key].setupAfterCallback();
            this.payload[key].begin();
        });
    }

    protected viewDebugInfo() {
        this.debugText.text = "FPS: " + this.app.ticker.FPS.toFixed(2)
                                + ", " + this.app.stage.children.length + " entities"
                                + ", debug=" + this.debugValue;
    }

    /**
     * 指定した id に対応するリソースを取得します
     * @param key リソース型
     * @param id id
     */
    get(key: string, id: string) {
        let container = this.payload[key];
        if (container !== undefined) {
            return container.getChild(id);
        }
        return undefined;
    }

    mergeAll(payload: GameMap) {
        config.zIndices.forEach(key => {
            if (this.payload[key] !== undefined) {
                this.payload[key].mergeChildren(payload[key], this.coord);
                if (this.payload[key].isChanged()) {
                    this.changed = true;
                }
            }
        });
        this.resolve();
        this.shouldFetch = false;
    }

    resolve() {
        if (this.payload["rail_edges"] !== undefined) {
            this.payload["rail_edges"].forEachChild((re: RailEdge) => 
                re.resolve(
                    this.get("rail_nodes", re.get("from")),
                    this.get("rail_nodes", re.get("to")),
                    this.get("rail_edges", re.get("eid"))
                )
            );
        }
    }

    setCenter(x: number, y: number, force: boolean = false) {
        let radius = Math.pow(2, this.coord.scale - 1);
        if (x - radius < config.gamePos.min.x) {
            x = config.gamePos.min.x + radius;
        }
        if (y - radius < config.gamePos.min.y) {
            y = config.gamePos.min.y + radius;
        }
        if (x + radius > config.gamePos.max.x) {
            x = config.gamePos.max.x - radius;
        }
        if (y + radius > config.gamePos.max.y) {
            y = config.gamePos.max.y - radius;
        }
        if (this.coord.cx == x && this.coord.cy == y) {
            return;
        }
        this.shouldFetch = true;
        this.shouldRemoveOutsider = true;
        this.coord.cx = x;
        this.coord.cy = y;
        
        Object.keys(this.payload).forEach(key => {
            this.payload[key].mergeAll(this.coord);
            if (force) {
                this.payload[key].mergeAll(forceMove);
            }
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        });
    }

    setScale(v: number, force: boolean = false) {
        if (v < config.scale.min) {
            v = config.scale.min;
        }
        if (v > config.scale.max) {
            v = config.scale.max;
        }
        if (this.coord.scale == v) {
            return;
        } else if (v > this.coord.scale) { // 縮小
            this.shouldFetch = true;
        } else { // 拡大
            this.shouldRemoveOutsider = true;
        }
        this.coord.scale = v;

        Object.keys(this.payload).forEach(key => {
            this.payload[key].mergeAll(this.coord);
            if (force) {
                this.payload[key].mergeAll(forceMove);
            }
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        });
    }

    isChanged() {
        return this.changed;
    }

    render() {
        Object.keys(this.payload).forEach(key => 
            this.payload[key].forEachChild((c) => c.beforeRender())
        );
        Object.keys(this.payload).forEach(key => this.payload[key].reset());
        this.changed = false;
    }

    unmount() {
        Object.keys(this.payload).reverse().forEach(key => this.payload[key].end());

        Object.keys(this.payload).reverse().forEach(key => {
            this.payload[key].end();
        });
    }

    removeOutsider() {
        Object.keys(this.payload).forEach(key => this.payload[key].removeOutsider());
    }
}
