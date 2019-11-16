import { Point, getChunkByPos } from "../interfaces/gamemap";
import { startMonitor, MonitorContainer, Monitorable } from "../interfaces/monitor";
import GameModel from ".";


export default abstract class {
    protected model: GameModel;
    protected containers: {[index: string]: MonitorContainer } = {};
    protected objects: {[index: string]: Monitorable } = {};
    protected changed: boolean = false;

    constructor(game: GameModel) {
        this.model = game;
    }

    init() {
        this.forEach(v => startMonitor(v, {}));
    }

    tick() {
        this.forEachContainer(v => v.merge("offset", this.model.offset));
        this.forEachContainer(v => v.endChildren());
    }

    /**
     * 指定した id に対応するリソースを取得します
     * @param key リソース型
     * @param id id
     */
    get(key: string, id: string | undefined) {
        let container = this.containers[key];
        if (container !== undefined) {
            return container.getChild(id);
        }
        return undefined;
    }

    getOnChunk(key: string, pos: Point | undefined, oid: number): Monitorable | undefined {
        if (this.containers[key] === undefined || pos === undefined) {
            return undefined;
        }
        return this.containers[key].getChildOnChunk(getChunkByPos(pos, this.model.coord.scale - this.model.delegate + 1), oid)
    }

    merge(key: string, value: any) {
        this.forEach(v => {
            v.merge(key, value);
            if (v.isChanged()) {
                this.changed = true;
            }
        });
    }

    updateDisplayInfo() {
        this.forEachContainer(v => {
            v.forEachChild(c => c.updateDisplayInfo());
        });
        this.forEach(v => v.reset());
        this.changed = false;
    }

    isChanged() {
        return this.changed;
    }

    end() {
        this.forEach(v => v.end());
    }

    protected forEach(fn: (v: Monitorable) => any) {
        this.forEachContainer(fn);
        this.forEachObject(fn);
    }

    protected forEachContainer(fn: (v: MonitorContainer) => any) {
        Object.keys(this.containers).forEach(key => fn(this.containers[key]));
    }

    protected forEachObject(fn: (v: Monitorable) => any) {
        Object.keys(this.objects).forEach(key => fn(this.objects[key]));
    }

}