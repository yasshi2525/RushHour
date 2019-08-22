import * as React from "react";
import { MenuStatus } from "../../state";
import { depart, fetchMap } from "../../actions";
import GameModel from "../models";
import { Cursor, Anchor } from "../models/cursor";
import { Point } from "../interfaces/gamemap";
import { getZoomPos } from "./point";

const offset = 2;
export abstract class CursorHandler<T> {
    protected model: GameModel;
    protected view: Cursor;
    protected anchor: Anchor;
    protected dispatch: any;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.view = model.controllers.getCursor();
        this.anchor = model.controllers.getAnchor();
        this.dispatch = dispatch;
    }

    /**
     * 選択した点に複数の線路ノードが存在したため、拡大して選ばせる
     */
    protected requestZoom(client: Point) {
        let center = getZoomPos(
            this.model, {x: this.model.coord.cx, y: this.model.coord.cy}, 
            this.model.coord.scale, client, -1)
        this.model.setCoord(center.x, center.y, this.model.coord.scale - 1);
        this.dispatch(fetchMap.request({
            model: this.model,
            dispatch: this.dispatch
        }));
    }

    protected handle(client: Point) {
        let server = this.toServerXY(client, offset);
        if (server === undefined) {
            return;
        }
        switch(this.model.menu) {
            case MenuStatus.SEEK_DEPARTURE:
                if (this.view.selected === undefined) {
                    this.dispatch(depart.request({
                        model: this.model,
                        dispatch: this.dispatch,
                        oid: 1, // TODO
                        x: server.x, y: server.y,
                        scale: Math.floor(this.model.coord.scale - this.model.delegate + 1)
                    }));
                } else {
                    if (this.view.selected.get("mul") === 1) {
                        this.anchor.merge("anchor", this.view.genAnchorStatus())
                        this.model.setMenuState(MenuStatus.EXTEND_RAIL);
                    } else {
                        this.requestZoom(client);
                    }
                }
                break;
            case MenuStatus.EXTEND_RAIL:
                if (this.view.selected === undefined) {
                    console.log("TODO: send extend request");
                } else {
                    if (this.view.selected.get("mul") === 1) {
                        console.log("TODO: connect");
                    } else {
                        this.requestZoom(client);
                    }
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