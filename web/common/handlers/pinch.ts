import * as React from "react";
import GameModel from "../model";
import BaseHandler from "./base";

export class PinchHandler extends BaseHandler<React.TouchEvent> {
    dist: {from: number, to: number};

    constructor(model: GameModel) {
        super(model);
        this.dist = {from: 0, to: 0};
    }

    getDistance(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        let dx = ts.item(0).clientX - ts.item(1).clientX;
        let dy = ts.item(0).clientY - ts.item(1).clientY;
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
        this.dist.to = this.getDistance(ev);

        let ratio = this.dist.to / this.dist.from;
        this.scale.to = this.scale.from + (ratio - 1);
    }

    shouldEnd(ev: React.TouchEvent) {
        return ev.targetTouches.length > 1;
    }

    handleEnd() {
    }
}
