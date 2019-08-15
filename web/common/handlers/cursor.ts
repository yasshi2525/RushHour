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

    protected handle(server: Point) {
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

    protected abstract getClientXY(ev: T): Point;
}

export class ClickCursor extends CursorHandler<React.MouseEvent> {
    onMove(ev: React.MouseEvent) {
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
    
    onClick(ev: React.MouseEvent) {
        this.handle(this.toServerXY(this.getClientXY(ev)));
    }

    protected getClientXY(ev: React.MouseEvent) {
        return {x: ev.clientX, y: ev.clientY};
    }
}

const sensitivity = 1;

export class TapCursor extends CursorHandler<React.TouchEvent> {
    protected pos: Point = {x: -1, y: -1};
    protected moveCnt = 0;

    onStart(ev: React.TouchEvent) {
        this.pos = this.getClientXY(ev);
        this.moveCnt = 0;
    }

    onMove(_ev: React.TouchEvent) {
        this.moveCnt++;
    }

    onEnd(_ev: React.TouchEvent) {
        if (this.pos.x !== -1 && this.pos.y !== -1 && this.moveCnt <= sensitivity) {
            this.cursor.merge("x", this.pos.x);
            this.cursor.merge("y", this.pos.y);
            this.cursor.beforeRender();

            this.handle(this.toServerXY(this.pos));

            this.cursor.merge("x", -1);
            this.cursor.merge("y", -1);
            this.cursor.beforeRender();
        }
        this.pos = {x: -1, y: -1};
        this.moveCnt = 0;
    }

    protected getClientXY(ev: React.TouchEvent) {
        let ts = ev.touches;
        return (ts.length === 1) ? { 
            x: ts.item(0).clientX * this.model.renderer.resolution, 
            y: ts.item(0).clientY * this.model.renderer.resolution}
            : {x: -1, y: -1}
    }
}