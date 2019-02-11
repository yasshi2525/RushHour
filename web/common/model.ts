import { Residence, Company } from "./models/background";
import MonitorContainer from "./models/container";
import { Monitorable } from "./interfaces/monitor";
import { GameMap } from "../state";
import { RailNode, RailEdge } from "./models/rail";

export default class {
    stage: PIXI.Container; 
    loader: PIXI.loaders.Loader;
    renderer: PIXI.CanvasRenderer | PIXI.WebGLRenderer;
    payload: {[index:string]: MonitorContainer<Monitorable>} = {}
    changed: boolean = false;

    constructor(options: {app: PIXI.Application}) {
        this.stage = options.app.stage;
        this.loader = options.app.loader;
        this.renderer = options.app.renderer;

        this.payload["residences"] = new MonitorContainer(Residence, {name: "residence", app: options.app});
        this.payload["companies"] = new MonitorContainer(Company, {name: "company", app: options.app});
        this.payload["rail_nodes"] = new MonitorContainer(RailNode, {app: options.app})
        this.payload["rail_edes"] = new MonitorContainer(RailEdge, {app: options.app})
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
    }

    isChanged() {
        return this.changed;
    }

    render() {
        this.renderer.render(this.stage);
        Object.keys(this.payload).forEach(key => this.payload[key].reset());
        this.changed = false;
    }

    unmount() {
        Object.keys(this.payload).forEach(key => this.payload[key].end());
    }
}
