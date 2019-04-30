import { Edge } from "../interfaces/gamemap";
import GameModel from "../model";

export default abstract class <T> {
    protected isExec: boolean;
    protected client: Edge;
    protected server: Edge;
    protected scale: {from: number, to: number};
    protected model: GameModel;

    constructor(model: GameModel) {
        this.model = model;
        this.isExec = false;
        this.client = {from: {x: 0, y: 0}, to: {x: 0, y: 0}};
        this.server = {from: {x: model.cx, y: model.cy}, to: {x: model.cx, y: model.cy}};
        this.scale = {from: model.scale, to: model.scale};
    }

    protected abstract shouldStart(ev: T): boolean

    onStart(ev: T) {
        if (this.shouldStart(ev)) {
            this.isExec = true;
            this.server.from = {x: this.model.cx, y: this.model.cy};
            this.server.to = {x: this.model.cx, y: this.model.cy};
            this.scale.from = this.model.scale;
            this.scale.to = this.model.scale;
            this.handleStart(ev);
        }
    }

    protected abstract handleStart(ev: T): void

    onMove(ev: T) {
        if (this.isExec) {
            this.handleMove(ev);

            this.model.setScale(this.scale.to);
            this.model.setCenter(this.server.to.x, this.server.to.y);

            if (this.model.isChanged()) {
                this.model.render();
            }
        } 
    }

    protected abstract handleMove(ev: T): void

    protected abstract shouldEnd(ev: T): boolean

    onEnd(ev: T) {
        if (this.shouldEnd(ev)) {
            this.isExec = false;
            this.handleEnd(ev);
        }
    }

    protected abstract handleEnd(ev: T): void;
}