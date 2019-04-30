import * as React from "react";
import BaseHandler from "./base";

const sensitivity = 0.001;

export class WheelHandler extends BaseHandler<React.WheelEvent> {

    protected shouldStart() {
        return true;
    }    
    
    protected handleStart() {
    }

    protected handleMove(ev: React.WheelEvent): void {
        // scale変更
        this.model.setScale(this.scale.from + ev.deltaY * sensitivity)
        this.scale.to = this.model.scale;

        // 画面中心からのカーソル相対位置取得
        let size = Math.max(
            this.model.renderer.width, 
            this.model.renderer.height
        );
        let zoom = Math.pow(2, this.scale.from);

        let dC = {
            x: ev.clientX - this.model.renderer.width / 2,
            y: ev.clientY - this.model.renderer.height / 2
        };
        
        let dS = {
            x: dC.x / size * zoom,
            y: dC.y / size * zoom,
        };

        // 中心座標の変更
        let fromLen = Math.sqrt(dS.x * dS.x + dS.y * dS.y);
        if (fromLen > 0) {
            let toLen = fromLen * Math.pow(2, this.scale.to - this.scale.from);

            let theta = Math.atan2(dS.y, dS.x);
            this.server.to = {
                x: this.server.from.x - (toLen - fromLen) * Math.cos(theta),
                y: this.server.from.y - (toLen - fromLen) * Math.sin(theta)
            };  
        }
    }

    protected shouldEnd(): boolean {
        return true;
    }

    protected handleEnd() {
    }
} 