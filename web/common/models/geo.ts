import * as PIXI from "pixi.js";
import { Monitorable } from "../interfaces/monitor";
import BaseModel from "./base";

export abstract class PIXIModel extends BaseModel implements Monitorable {
    protected app: PIXI.Application;
    protected container: PIXI.Container;

    constructor(options: {app: PIXI.Application}) {
        super();
        this.app = options.app;
        this.container = new PIXI.Container();
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.app.stage.addChild(this.container);
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            this.app.stage.removeChild(this.container);
        })
    }
}

const defaultValues: {x: number, y:number} = {x: 0, y: 0};

export abstract class PointModel extends PIXIModel implements Monitorable {

    constructor(options: {app: PIXI.Application}) {
        super(options);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}
