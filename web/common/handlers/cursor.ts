import * as React from "react";
import { MenuStatus } from "../../state";
import GameModel from "../model";
import { Point } from "../interfaces/gamemap";

abstract class CursorHandler<T> {
    protected model: GameModel;
    protected dispatch: any;
    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.dispatch = dispatch;
    }

    protected toServerXY(client: Point) {
        let w = this.model.renderer.width;
        let h = this.model.renderer.height;
        let size = Math.max(
            this.model.renderer.width, 
            this.model.renderer.height
        );

        let d = {
            x: (client.x - w / 2) / size,
            y: (client.y - h / 2) / size
        }

        let zoom = Math.pow(2, this.model.coord.scale);
        return {
            x: this.model.coord.cx + d.x * zoom,
            y: this.model.coord.cy + d.y * zoom
        }
    }
    
    onClick(menu: MenuStatus, _ev: T) {
        //let loc = this.toServerXY(this.getClientXY(ev));
        switch(menu) {
            case MenuStatus.SEEK_DEPARTURE:
                break;
        }
    }

    protected abstract getClientXY(ev: T): Point;
}

export class ClickCursor extends CursorHandler<React.MouseEvent> {
    protected getClientXY(ev: React.MouseEvent) {
        return {x: ev.clientX, y: ev.clientY};
    }
}