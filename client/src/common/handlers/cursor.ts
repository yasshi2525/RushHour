import * as React from "react";
import { MenuStatus } from "../../state";
import { depart, extend, connect, fetchMap, destroy } from "../../actions";
import GameModel from "../models";
import Anchor from "../models/anchor";
import Cursor from "../models/cursor";
import { RailNode } from "../models/rail";
import { Point } from "../interfaces/gamemap";
import { getZoomPos } from "./point";

export abstract class CursorHandler<T> {
    protected model: GameModel;
    protected cursor: Cursor;
    protected anchor: Anchor;
    protected dispatch: any;

    constructor(model: GameModel, dispatch: any) {
        this.model = model;
        this.cursor = model.controllers.getCursor();
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
        this.dispatch(fetchMap.request({ model: this.model }));
    }

    protected handle(client: Point) {
        let server = this.cursor.get("pos");
        if (server === undefined) {
            return;
        }
        switch(this.model.menu) {
            case MenuStatus.SEEK_DEPARTURE:
                if (this.cursor.selected === undefined) {
                    this.dispatch(depart.request({
                        model: this.model,
                        x: server.x, y: server.y,
                        scale: Math.floor(this.model.coord.scale - this.model.delegate + 1)
                    }));
                } else {
                    if (this.cursor.selected.get("mul") === 1) {
                        this.model.setMenuState(MenuStatus.EXTEND_RAIL);
                        this.anchor.merge("anchor", this.cursor.genAnchorStatus());
                    } else {
                        this.requestZoom(client);
                    }
                }
                break;
            case MenuStatus.EXTEND_RAIL:
                if (this.cursor.selected === undefined) {
                    if (this.anchor.object !== undefined && this.cursor.get("activation")) {
                        this.dispatch(extend.request({
                            model: this.model,
                            x: server.x, y: server.y,
                            rnid: this.anchor.object.get("cid"),
                            scale: Math.floor(this.model.coord.scale - this.model.delegate + 1)
                        }));
                    }
                } else {
                    if (this.cursor.selected.get("mul") === 1) {
                        if (this.anchor.object !== undefined && this.cursor.get("activation")) {
                            this.dispatch(connect.request({
                                model: this.model,
                                from: this.anchor.object.get("cid"),
                                to: this.cursor.selected.get("cid"),
                                scale: Math.floor(this.model.coord.scale - this.model.delegate + 1)
                            }));
                        }
                    } else {
                        this.requestZoom(client);
                    }
                }
                break;
            case MenuStatus.DESTROY:
                if (this.cursor.destroyer.selected !== undefined) {
                    if (this.cursor.destroyer.selected.get("mul") === 1) {
                        let type: string = "";

                        if (this.cursor.destroyer.selected instanceof RailNode) {
                            type = "rail_nodes"
                        }

                        this.dispatch(destroy.request({
                            model: this.model,
                            resource: type,
                            id: this.cursor.destroyer.selected.get("id"),
                            cid: this.cursor.destroyer.selected.get("cid"),
                            scale: Math.floor(this.model.coord.scale - this.model.delegate + 1)
                        }))
                    } else {
                        this.requestZoom(client);
                    }
                }
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
        this.cursor.merge("client", this.getClientXY(ev));
        this.moveCnt++;
    }

    onOut(_ev: React.MouseEvent) {
        this.cursor.merge("client", undefined);
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
            this.cursor.merge("client", this.pos);
            this.handle(this.pos);
            this.cursor.merge("client", undefined);
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