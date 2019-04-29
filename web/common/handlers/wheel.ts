import * as React from "react";
import BaseHandler from "./base";

const delta = 0.4;

export class WheelHandler extends BaseHandler<React.WheelEvent> {

    protected shouldStart() {
        return true;
    }    
    
    protected handleStart() {
    }

    protected handleMove(ev: React.WheelEvent): void {
        if (ev.deltaY > 0) {
            this.scale.to = this.scale.from + delta;
        }
        if (ev.deltaY < 0) {
            this.scale.to = this.scale.from - delta;
        }
    }

    protected shouldEnd(): boolean {
        return true;
    }

    protected handleEnd() {
    }
} 