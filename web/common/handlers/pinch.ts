import * as React from "react";
import GameModel from "../model";
import PointHandler from "./point";

export class PinchHandler extends PointHandler<React.TouchEvent> {
    dist: {from: number, to: number};

    constructor(model: GameModel, dispatch: any) {
        super(model, dispatch);
        this.dist = {from: 0, to: 0};
    }

    /**
     * 2つのタッチポイントの重心を取得します。
     * @param ev イベント
     */
    protected getClientXY(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        return {
            x: (ts.item(0).clientX + ts.item(1).clientX) * this.model.renderer.resolution / 2, 
            y: (ts.item(0).clientY + ts.item(1).clientY) * this.model.renderer.resolution / 2
        };
    }

    protected getDistance(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        let dx = (ts.item(0).clientX - ts.item(1).clientX) * this.model.renderer.resolution;
        let dy = (ts.item(0).clientY - ts.item(1).clientY) * this.model.renderer.resolution;
        return Math.sqrt(dx * dx + dy * dy);
    }

    shouldStart(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    handleStart(ev: React.TouchEvent) {
        this.dist.from = this.getDistance(ev);
        this.dist.to = this.dist.from;
    }

    handleMove(ev: React.TouchEvent) {
        // scaleの変更
        this.dist.to = this.getDistance(ev);

        let ratio = this.dist.from / this.dist.to;
        this.model.setScale(this.scale.from + (ratio - 1))
        this.scale.to = this.model.coord.scale;

        // 画面中心座標の変更
        let center = this.zoom(this.getClientXY(ev), this.scale.to - this.scale.from);
        this.server.to.x = center.x;
        this.server.to.y = center.y;
    }

    shouldEnd(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    handleEnd() {
    }
}
