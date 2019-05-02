import * as React from "react";
import PointHandler from "./point";

const sensitivity = 0.001;

export class WheelHandler extends PointHandler<React.WheelEvent> {

    protected getClientXY(ev: React.WheelEvent) {
        return {
            x: ev.clientX, 
            y: ev.clientY
        };
    }

    protected shouldStart() {
        return true;
    }    
    
    protected handleStart() {
    }

    protected handleMove(ev: React.WheelEvent): void {
        // scaleの変更
        this.model.setScale(this.scale.from + ev.deltaY * sensitivity)
        this.scale.to = this.model.coord.scale;

        // 画面中心座標の変更
        let center = this.zoom(this.getClientXY(ev), this.scale.to - this.scale.from);
        this.server.to.x = center.x;
        this.server.to.y = center.y;
    }

    protected shouldEnd(): boolean {
        return true;
    }

    protected handleEnd() {
    }
} 