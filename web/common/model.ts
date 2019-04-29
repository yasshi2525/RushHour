import { Residence, Company } from "./models/background";
import MonitorContainer from "./models/container";
import { Monitorable } from "./interfaces/monitor";
import { LocalableProprty } from "./interfaces/pixi";
import { GameMap } from "../state";
import { RailNode, RailEdge } from "./models/rail";

export default class {
    protected stage: PIXI.Container; 
    protected loader: PIXI.loaders.Loader;
    protected renderer: PIXI.CanvasRenderer | PIXI.WebGLRenderer;
    protected payload: {[index:string]: MonitorContainer<Monitorable>} = {}
    protected changed: boolean = false;
    /**
     * 中心x座標(サーバにおけるマップ座標系)
     */
    cx: number;
    /**
     * 中心y座標(サーバにおけるマップ座標系)
     */
    cy: number;
    /**
     * 拡大率(クライエントウィンドウの幅が2^scaleに対応する)
     */
    scale: number;

    constructor(options: LocalableProprty) {
        this.stage = options.app.stage;
        this.loader = options.app.loader;
        this.renderer = options.app.renderer;

        this.payload["residences"] = new MonitorContainer(Residence, {name: "residence", ...options});
        this.payload["companies"] = new MonitorContainer(Company, {name: "company", ...options});
        this.payload["rail_nodes"] = new MonitorContainer(RailNode, options);
        this.payload["rail_edges"] = new MonitorContainer(RailEdge, options);
        this.cx = 0;
        this.cy = 0;
        this.scale = 10;
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
                this.payload[key].mergeChildren(payload[key]);
                if (this.payload[key].isChanged()) {
                    this.changed = true;
                }
            }
        });
        this.resolve();
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
        this.cx = x;
        this.cy = y;
        Object.keys(this.payload).forEach(key => {
            this.payload[key].childOptions.cx = x;
            this.payload[key].childOptions.cy = y;
            this.payload[key].mergeAll({cx: x, cy: y})
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        })
    }

    setScale(v: number) {
        Object.keys(this.payload).forEach(key => {
            this.payload[key].mergeAll({scale: v})
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        })
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
}
