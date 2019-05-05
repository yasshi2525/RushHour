import * as PIXI from "pixi.js";
import { config, Coordinates, Point } from "../interfaces/gamemap";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { ApplicationProperty, ContainerProperty } from "../interfaces/pixi";
import BaseContainer from "./container";
import BaseModel from "./base";

const defaultValues: Coordinates & {[index: string]: any} = {
    cx: config.gamePos.default.x, 
    cy: config.gamePos.default.y, 
    scale: config.scale.default
};

export abstract class PIXIModel extends BaseModel implements Monitorable {
    protected app: PIXI.Application;
    protected parent: PIXI.Container;
    protected container: PIXI.Container;

    constructor(options: ContainerProperty) {
        super();
        this.app = options.app;
        this.parent = options.container;
        this.container = new PIXI.Container();
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.parent.addChild(this.container);
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            this.parent.removeChild(this.container);
        })
    }

    toView(x: number, y: number): Point {
        let center = {
            x: this.app.renderer.width / this.app.renderer.resolution / 2,
            y: this.app.renderer.height / this.app.renderer.resolution / 2,
        }
        let size = Math.max(this.app.renderer.width / this.app.renderer.resolution, this.app.renderer.height / this.app.renderer.resolution)
        let zoom = Math.pow(2, -this.props.scale)

        return {
            x: (x - this.props.cx) * size * zoom + center.x,
            y: (y - this.props.cy) * size * zoom + center.y
        }
    }

    /**
     * scale + 1 の範囲をキャッシュ保持領域としたとき、それを外れたかどうか判定する
     * @param x サーバ座標系x座標
     * @param y サーバ座標系y座標
     */
    protected isOut(x: number, y: number) {
        let zoom = Math.pow(2, this.props.scale);
        return Math.abs(x - this.props.cx) > zoom || Math.abs(y - this.props.cy) > zoom;
    }

    shouldEnd() {
        return this.isOut(this.props.x, this.props.y);
    }
}

const animationOpts = {
    round: 5000
};

export abstract class PIXIContainer<T extends PIXIModel> extends BaseContainer<T> implements MonitorContrainer {
    protected app: PIXI.Application;
    protected container: PIXI.Container;
    /**
     * アニメーション時の進行率(0-1) cosカーブ
     */
    protected offset: number;

    /**
     * インスタンス作成からの累計時間
     */
    protected tick: number;

    protected counterFn: () => void;

    constructor(
        options: ApplicationProperty,
        newInstance: { new (props: {[index:string]: {}}): T }, 
        newInstanceOptions: {[index:string]: {}}) {
        super(newInstance, newInstanceOptions);
        this.app = options.app;
        this.container = new PIXI.Container();
        this.childOptions.app = this.app;
        this.childOptions.container = this.container;
        this.tick = 0;
        this.offset = 0;
        this.counterFn = () => this.counter();
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.app.stage.addChild(this.container);
            this.app.ticker.add(this.counterFn)
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            this.app.ticker.remove(this.counterFn);
            this.app.stage.removeChild(this.container);
        })
    }

    protected counter() {
        this.tick += this.app.ticker.elapsedMS;
        let ratio = (this.tick % animationOpts.round) / animationOpts.round;
        this.offset = Math.cos(ratio * Math.PI * 2) / 2 + 0.5;
    }
}