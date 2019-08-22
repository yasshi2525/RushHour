import GameModel from ".";
import BaseModel from "./base";
import { MonitorContainer, Monitorable } from "../interfaces/monitor";
import { Chunk } from "../interfaces/gamemap";
import { PointModel } from "./point";

const defaultValues: {offset: number} = {offset: 0};

export default abstract class <T extends Monitorable> extends BaseModel implements MonitorContainer {

    childOptions: {[index:string]: {}};
    Child: { new (props: {[index:string]: {}}): T };

    children: {[index: string]: T};

    constructor(
        model: GameModel,
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: { [index:string]: {}}) {
        super({ model: model});
        this.Child = newInstance;
        this.childOptions = { ...newInstanceOptions, model };
        this.children = {};
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("offset", v => this.childOptions.offset = v);
    }

    existsChild(id: string) {
        return this.children[id] !== undefined;
    }

    getChild(id: string) {
        return this.children[id];
    }

    getChildOnChunk(chunk: Chunk, oid: number): PointModel | undefined {
        let result = Object.keys(this.children).map(id => this.children[id])
            .find(c => c.get("oid") === oid && !c.get("outMap") && c.standOnChunk(chunk));
        return (result instanceof PointModel) ? result : undefined;
    }

    addChild(props: {id: string}) {
        let child = new this.Child(this.childOptions);
        child.setupDefaultValues();
        child.setupUpdateCallback();
        child.setupBeforeCallback();
        child.setupAfterCallback();
        child.setInitialValues({ ...props, coord: this.model.coord });
        child.begin();

        this.children[props.id] = child;
        this.change();

        return child;
    }
    
    updateChild(props: {id: string}) {
        let target = this.children[props.id];

        target.mergeAll(props);

        if (target.isChanged()) {
            this.change();
        }
        return target;
    }

    endChildren() {
        this.forEachChild(ch => {
            if (ch.shouldEnd()) {
                this.removeChild(ch.get("id"))
            }
        })
    }

    removeChild(id: string) {
        if (this.existsChild(id)) {
            this.children[id].end();
            delete this.children[id];
            this.change();
        }
    }

    mergeChild(props: {id: string, [propName: string]: any}) {
        if (this.existsChild(props.id)) {
            return this.updateChild(props);
        } else {
            return this.addChild(props);
        }
    }

    mergeChildren(payload: {id: string, [propName: string]: any}[], opts: {[index: string]: any}) {
        if (payload === undefined) {
            return
        }
        payload.forEach(props => {
            Object.assign(props, opts, {outMap: false})
            this.mergeChild(props);
        });

        // payloadに含まれない child に outMap をつける
        let aliveIds = payload.map(props => props.id);
        Object.keys(this.children)
            .filter(myId => !aliveIds.find(id => myId == id))
            .forEach(id => this.getChild(id).merge("outMap", true));
    }

    /**
     * すべてのchildのkeyにvalueを設定します
     * @param key 設定するプロパティ名
     * @param value プロパティに設定する値
     */
    merge(key: string, value: any) {
        super.merge(key, value);
        this.forEachChild(c => {
            c.merge(key, value);
            if (c.isChanged()) {
                this.change();
            }
        });
    }

    reset() {
        this.forEachChild(c => c.reset());
        super.reset();
    }

    end() {
        this.forEachChild(c => c.end());
        super.end();
    }

    forEachChild(func: (c: T) => any) {
        Object.keys(this.children).forEach(id => func(this.children[id]));
    }
}