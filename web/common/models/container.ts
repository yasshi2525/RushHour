import BaseModel from "./base";
import { MonitorContrainer, Monitorable } from "../interfaces/monitor";


export default class <T extends Monitorable> extends BaseModel implements MonitorContrainer {

    childOptions: {[index:string]: {}};
    Child: { new (props: {[index:string]: {}}): T };

    children: {[index: string]: T};

    constructor(
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: {[index:string]: {}} = {}) {
        super();

        this.Child = newInstance;
        this.childOptions = newInstanceOptions;
        this.children = {};
    }

    existsChild(id: string) {
        return this.children[id] !== undefined;
    }

    getChild(id: string) {
        return this.children[id];
    }

    addChild(props: {id: string}) {
        let child = new this.Child(this.childOptions);
        child.setupDefaultValues();
        child.setupUpdateCallback();
        child.setupBeforeCallback();
        child.setupAfterCallback();
        child.setInitialValues(props);
        child.begin();

        this.children[props.id] = child;
        this.change();
    }
    
    updateChild(props: {id: string}) {
        let target = this.children[props.id];

        target.mergeAll(props);

        if (target.isChanged()) {
            this.change();
        }
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
            this.updateChild(props);
        } else {
            this.addChild(props);
        }
    }

    mergeChildren(payload: {id: string, [propName: string]: any}[]) {
        payload.forEach((props => this.mergeChild(props)));

        // payloadに含まれない child を削除する
        let aliveIds = payload.map(props => props.id);
        Object.keys(this.children)
            .filter(myId => !aliveIds.find(id => myId == id))
            .forEach(id => this.removeChild(id));
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