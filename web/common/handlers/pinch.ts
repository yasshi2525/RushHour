import * as React from "react";
import GameModel from "../models";
import { PointHandler } from "./point";

export class PinchHandler extends PointHandler<React.TouchEvent> {
    dist: {from: number, to: number};

    constructor(model: GameModel, dispatch: any) {
        super(model, dispatch);
        this.dist = {from: 0, to: 0};
        this.forceMove = true;
    }

    protected getClientXY(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        let pos = {x: 0, y: 0};

        for (let i = 0; i < ts.length; i++) {
            pos.x += ts.item(i).clientX / ts.length * this.model.renderer.resolution;
            pos.y += ts.item(i).clientY / ts.length * this.model.renderer.resolution;
        }

        return pos;
    }

    protected getDistance(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        let dx = (ts.item(0).clientX - ts.item(1).clientX) * this.model.renderer.resolution;
        let dy = (ts.item(0).clientY - ts.item(1).clientY) * this.model.renderer.resolution;
        return Math.sqrt(dx * dx + dy * dy);
    }

    protected shouldStart(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    protected handleStart(ev: React.TouchEvent) {
        this.dist.from = this.getDistance(ev);
        this.dist.to = this.dist.from;
    }

    protected shouldMove(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    protected handleMove(ev: React.TouchEvent) {
        this.dist.to = this.getDistance(ev);

        let d = this.dist.to - this.dist.from;

        let size = Math.max(
            this.model.renderer.width, 
            this.model.renderer.height
        );

        let ratio = d / size;
        this.scale.to = this.scale.from - ratio * this.model.renderer.resolution;
        let center = this.zoom(this.getClientXY(ev), this.scale.to - this.scale.from);
        this.server.to.x = center.x;
        this.server.to.y = center.y;
    }

    protected shouldEnd(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    protected handleEnd() {
    }
}
