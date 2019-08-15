import * as React from "react";
import { MenuStatus } from "../../state";
import { depart } from "../../actions";
import GameModel from "../model";
import CursorModel from "../models/cursor";
import { Point } from "../interfaces/gamemap";

const offset = 2;

abstract class CursorHandler<T> {
    protected model: GameModel;
    protected cursor: CursorModel;
    protected dispatch: any;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.cursor = model.cursor;
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
            x: (client.x + offset - w / 2) / size,
            y: (client.y + offset - h / 2) / size
        }

        let zoom = Math.pow(2, this.model.coord.scale);
        return {
            x: this.model.coord.cx + d.x * zoom,
            y: this.model.coord.cy + d.y * zoom
        }
    }

    onMove(ev: T) {
        let client = this.getClientXY(ev);
        this.cursor.merge("x", client.x);
        this.cursor.merge("y", client.y);
        if (this.cursor.isChanged()) {
            this.cursor.beforeRender();
        }
    }

    onMouseOut() {
        this.cursor.merge("x", -1);
        this.cursor.merge("y", -1);
    }
    
    onClick(ev: T) {
        let server = this.toServerXY(this.getClientXY(ev));
        switch(this.cursor.get("menu")) {
            case MenuStatus.SEEK_DEPARTURE:
                this.dispatch(depart.request({
                    oid: 1, // TODO
                    x: server.x,
                    y: server.y
                }))
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