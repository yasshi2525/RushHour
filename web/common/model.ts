import { Residence, Company } from "./models/background";
import MonitorContainer from "./models/container";
import { Coordinates, config } from "./interfaces/gamemap";
import { Monitorable } from "./interfaces/monitor";
import { ApplicationProperty } from "./interfaces/pixi";
import { GameMap } from "../state";
import { RailNode, RailEdge } from "./models/rail";

export default class {
    protected stage: PIXI.Container; 
    protected loader: PIXI.loaders.Loader;
    renderer: PIXI.CanvasRenderer | PIXI.WebGLRenderer;
    protected payload: {[index:string]: MonitorContainer<Monitorable>} = {}
    protected changed: boolean = false;
    timestamp: number;
    coord: Coordinates;
    /**
     * サーバから全データの取得が必要か
     */
    shouldFetch: boolean;
    shouldRemoveOutsider: boolean;
    debugText: PIXI.Text;

    constructor(options: ApplicationProperty & Coordinates) {
        this.stage = options.app.stage;
        this.loader = options.app.loader;
        this.renderer = options.app.renderer;

        this.payload["residences"] = new MonitorContainer(Residence, {name: "residence", ...options});
        this.payload["companies"] = new MonitorContainer(Company, {name: "company", ...options});
        this.payload["rail_nodes"] = new MonitorContainer(RailNode, options);
        this.payload["rail_edges"] = new MonitorContainer(RailEdge, options);
        this.coord = { cx: options.cx, cy: options.cy, scale: options.scale }
        this.shouldFetch = false;
        this.shouldRemoveOutsider = false;
        this.timestamp = 0;
        this.debugText = new PIXI.Text();
        this.debugText.y = window.innerHeight - 50;
        this.stage.addChild(this.debugText)
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
        Object.keys(payload).forEach(key => {
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
        this.payload["rail_edges"].forEachChild((re: RailEdge) => 
            re.resolve(
                this.get("rail_nodes", re.get("from")),
                this.get("rail_nodes", re.get("to")),
                this.get("rail_edges", re.get("eid"))
            )
        );
    }

    setCenter(x: number, y: number) {
        if (x < config.gamePos.min.x) {
            x = config.gamePos.min.x;
        }
        if (y < config.gamePos.min.y) {
            y = config.gamePos.min.y;
        }
        if (x > config.gamePos.max.x) {
            x = config.gamePos.max.x;
        }
        if (y > config.gamePos.max.y) {
            y = config.gamePos.max.y;
        }
        if (this.coord.cx == x && this.coord.cy == y) {
            return;
        }
        this.shouldFetch = true;
        this.shouldRemoveOutsider = true;
        this.coord.cx = x;
        this.coord.cy = y;
        
        Object.keys(this.payload).forEach(key => {
            this.payload[key].mergeAll(this.coord)
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        })
    }

    setScale(v: number) {
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
            this.payload[key].mergeAll(this.coord)
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
        this.renderer.render(this.stage);
        Object.keys(this.payload).forEach(key => this.payload[key].reset());
        this.changed = false;
    }

    unmount() {
        Object.keys(this.payload).forEach(key => this.payload[key].end());
    }

    removeOutsider() {
        Object.keys(this.payload).forEach(key => this.payload[key].removeOutsider());
    }
}
