import { Residence, Company } from "./models/background";
import MonitorContainer from "./models/container";
import { Monitorable } from "./interfaces/monitor";
import { GameMap } from "../state";

const spriteOption = (name: string, app: PIXI.Application) => ({
    name: name,
    container: app.stage,
    loader: app.loader
})

export class GameModel {
    stage: PIXI.Container; 
    loader: PIXI.loaders.Loader;
    renderer: PIXI.CanvasRenderer | PIXI.WebGLRenderer;
    payload: {[index:string]: MonitorContainer<Monitorable>} = {}

    constructor(options: {app: PIXI.Application}) {
        this.stage = options.app.stage;
        this.loader = options.app.loader;
        this.renderer = options.app.renderer;

        this.payload["residences"] = new MonitorContainer(Residence, spriteOption("residence", options.app));
        this.payload["companies"] = new MonitorContainer(Company, spriteOption("company", options.app));
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
            }
        });
    }

    isChanged() {
        return Object.keys(this.payload).find(key => this.payload[key].isChanged());
    }

    render() {
        this.renderer.render(this.stage);
        Object.keys(this.payload).forEach(key => this.payload[key].reset());
    }
}
