import { Edge } from "../interfaces/gamemap";
import { fetchMap } from "../../actions";
import GameModel from "../model";

export default abstract class <T> {
    protected isExec: boolean;
    protected client: Edge;
    protected server: Edge;
    protected scale: {from: number, to: number};
    protected model: GameModel;
    protected dispatch: any;
    protected forceMove: boolean;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.isExec = false;
        this.client = {from: {x: 0, y: 0}, to: {x: 0, y: 0}};
        this.server = {from: {x: model.coord.cx, y: model.coord.cy}, to: {x: model.coord.cx, y: model.coord.cy}};
        this.scale = {from: model.coord.scale, to: model.coord.scale};
        this.dispatch = dispatch;
        this.forceMove = false;
    }

    protected shouldStart(_: T) {
        return true;
    }

    onStart(ev: T) {
        if (this.shouldStart(ev)) {
            this.isExec = true;
            this.server.from = {x: this.model.coord.cx, y: this.model.coord.cy};
            this.server.to = {x: this.model.coord.cx, y: this.model.coord.cy};
            this.scale.from = this.model.coord.scale;
            this.scale.to = this.model.coord.scale;
            this.handleStart(ev);
        }
    }

    protected abstract handleStart(ev: T): void

    onMove(ev: T) {
        if (this.shouldMove(ev)) {
            this.handleMove(ev);

            this.model.setScale(this.scale.to, this.forceMove);
            this.model.setCenter(this.server.to.x, this.server.to.y, this.forceMove);
        }
        if (this.shouldEnd(ev)) {

        }
    }

    protected abstract handleMove(ev: T): void

    protected shouldMove(_: T) {
        return this.isExec;
    }

    protected shouldEnd(_: T) {
        return true;
    }

    onEnd(ev: T) {
        if (this.shouldEnd(ev)) {
            this.isExec = false;
            this.handleEnd(ev);
            if (this.model.shouldRemoveOutsider) {
                this.model.removeOutsider();
                this.model.shouldRemoveOutsider = false;
            }
            if (this.model.shouldFetch) {
                this.dispatch(fetchMap.request({
                    cx: this.model.coord.cx, 
                    cy: this.model.coord.cy, 
                    scale: this.model.coord.scale + 1
                }));
                this.model.shouldFetch = false;
            }
        }
    }

    protected abstract handleEnd(ev: T): void;
}