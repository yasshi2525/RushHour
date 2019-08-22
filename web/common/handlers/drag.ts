import * as styles from "./style.css";
import * as React from "react";
import { PointHandler } from "./point";
import GameModel from "../models";

abstract class DragHandler<T extends React.SyntheticEvent> extends PointHandler<T> {
    constructor(model: GameModel, dispatch: any) {
        super(model, dispatch);
        this.forceMove = true;
    }

    protected handleStart(ev: T) {
        this.client.from = this.getClientXY(ev);
        this.client.to = this.getClientXY(ev);
        ev.currentTarget.classList.add(styles.dragging);
    }

    protected handleMove(ev: T) {
        this.client.to = this.getClientXY(ev);

        let dx = (this.client.to.x - this.client.from.x);
        let dy = (this.client.to.y - this.client.from.y);
        let size = Math.max(
            this.model.renderer.width, 
            this.model.renderer.height
        );
        let zoom = Math.pow(2, this.model.coord.scale);

        this.server.to = {
            x: this.server.from.x - dx / size * zoom,
            y: this.server.from.y - dy / size * zoom
        };
    }

    protected handleEnd(ev: T) {
        ev.currentTarget.classList.remove(styles.dragging);
    }

    protected shouldFetch(ev: T) {
        return ev.type != "mouseout"
    }
}

export class MouseDragHandler extends DragHandler<React.MouseEvent> {

    protected getClientXY(ev: React.MouseEvent) {
        return {
            x: ev.clientX * this.model.renderer.resolution, 
            y: ev.clientY * this.model.renderer.resolution
        };
    }
}

export class TouchDragHandler extends DragHandler<React.TouchEvent> {
    constructor(model: GameModel, dispatch: any) {
        super(model, dispatch);
    }

    protected shouldStart(ev: React.TouchEvent) {
        return ev.touches.length == 1;
    }

    protected getClientXY(ev: React.TouchEvent) {
        let ts = ev.targetTouches;
        let pos = {x: 0, y: 0};

        for (let i = 0; i < ts.length; i++) {
            pos.x += ts.item(i).clientX / ts.length * this.model.renderer.resolution;
            pos.y += ts.item(i).clientY / ts.length * this.model.renderer.resolution;
        }

        return pos;
    }

    protected shouldEnd(ev: React.TouchEvent) {
        return ev.touches.length == 1;
    }
}
