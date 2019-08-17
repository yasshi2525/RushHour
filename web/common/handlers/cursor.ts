import * as React from "react";
import { MenuStatus } from "../../state";
import { depart } from "../../actions";
import GameModel from "../model";
import CursorModel from "../models/cursor";
import { Point } from "../interfaces/gamemap";

const offset = 2;
export abstract class CursorHandler<T> {
    protected model: GameModel;
    protected view: CursorModel;
    protected dispatch: any;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.view = model.cursor;
        this.dispatch = dispatch;
    }

    protected handle(client: Point) {
        let server = this.toServerXY(client, offset);
        if (server === undefined) {
            return;
        }
        switch(this.view.get("menu")) {
            case MenuStatus.SEEK_DEPARTURE:
                if (this.view.selected === undefined) {
                    this.dispatch(depart.request({
                        oid: 1, // TODO
                        x: server.x,
                        y: server.y
                    }));
                }
                break;
        }
    }

    protected toServerXY(client: Point | undefined, offset: number = 0) {
        if (client === undefined) {
            return undefined;
        }
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

    protected abstract getClientXY(ev: T): Point | undefined;
}

const sensitivity = 1;

export class ClickCursor extends CursorHandler<React.MouseEvent> {
    protected moveCnt = 0;

    onStart(_ev: React.MouseEvent) {
        this.moveCnt = 0;
    }

    onMove(ev: React.MouseEvent) {
        this.view.merge("pos", this.toServerXY(this.getClientXY(ev)));
        this.moveCnt++;
    }

    onOut(_ev: React.MouseEvent) {
        this.view.merge("pos", undefined);
        this.moveCnt = 0;
    }

    onEnd(ev: React.MouseEvent) {
        if (this.moveCnt <= sensitivity) {
            this.handle(this.getClientXY(ev));
        }
        this.moveCnt = 0;
    }

    protected getClientXY(ev: React.MouseEvent) {
        let result = {
            x: ev.clientX * this.model.renderer.resolution, 
            y: ev.clientY * this.model.renderer.resolution
        };
        return result;
    }
}


export class TapCursor extends CursorHandler<React.TouchEvent> {
    protected pos: Point | undefined;
    protected moveCnt = 0;

    onStart(ev: React.TouchEvent) {
        this.pos = this.getClientXY(ev);
        this.moveCnt = 0;
    }

    onMove(_ev: React.TouchEvent) {
        this.moveCnt++;
    }

    onEnd(_ev: React.TouchEvent) {
        if (this.pos !== undefined && this.moveCnt <= sensitivity) {
            this.view.merge("pos", this.toServerXY(this.pos));
            this.handle(this.pos);
            this.view.merge("pos", undefined);
        }
        this.pos = undefined;
        this.moveCnt = 0;
    }

    protected getClientXY(ev: React.TouchEvent) {
        let ts = ev.touches;

        if (ts.length === 1) {
            let result = { 
                x: ts.item(0).clientX * this.model.renderer.resolution, 
                y: ts.item(0).clientY * this.model.renderer.resolution
            }
            return result;
        }
        return undefined;
    }
}