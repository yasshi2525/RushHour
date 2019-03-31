import * as PIXI from "pixi.js";
import { Monitorable } from "../interfaces/monitor";
import { LocalableProprty } from "../interfaces/pixi";
import BaseModel from "./base";
import { number } from "prop-types";

const pixiDefaultValues: {cx: number, cy: number, scale: number} = {cx: 0, cy: 0, scale: 10};

export abstract class PIXIModel extends BaseModel implements Monitorable {
    protected app: PIXI.Application;
    protected container: PIXI.Container;

    constructor(options: LocalableProprty) {
        super();
        this.app = options.app;
        this.container = new PIXI.Container();
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(pixiDefaultValues);
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

    toView(x: number, y: number): {x: number, y:number} {
        return {
        x: (x - this.props.cx)
            * Math.max(this.app.renderer.width, this.app.renderer.height)
            * Math.pow(2, -this.props.scale)
            + this.app.renderer.width / 2,
        y: (y - this.props.cy)
            * Math.max(this.app.renderer.width, this.app.renderer.height)
            * Math.pow(2, -this.props.scale)
            + this.app.renderer.height / 2
        }
    }
}

const pointDefaultValues: {x: number, y:number} = {x: 0, y: 0};

export abstract class PointModel extends PIXIModel implements Monitorable {
    protected vx: number;
    protected vy: number;

    constructor(options: LocalableProprty) {
        super(options);
        this.vx = 0;
        this.vy = 0;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(pointDefaultValues);
    }

    beforeRender() {
        super.beforeRender();
        let v = this.toView(this.props.x, this.props.y);
        this.vx = v.x;
        this.vy = v.y;
    }
}
