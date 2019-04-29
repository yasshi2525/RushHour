import * as styles from "./style.css"
import * as React from "react";
import BaseHandler from "./base";
import { Point } from "../interfaces/gamemap";

abstract class DragHandler<T extends React.SyntheticEvent> extends BaseHandler<T> {

    protected abstract getClientXY(ev: T): Point;

    handleStart(ev: T) {
        this.client.from = this.getClientXY(ev);
        this.client.to = this.getClientXY(ev);
        ev.currentTarget.classList.add(styles.dragging);
    }

    handleMove(ev: T) {
        this.client.to = this.getClientXY(ev);

        let dx = this.client.to.x - this.client.from.x;
        let dy = this.client.to.y - this.client.from.y;
        let size = Math.max(ev.currentTarget.clientWidth, ev.currentTarget.clientHeight);
        let zoom = Math.pow(2, this.model.scale);

        this.server.to = {
            x: this.server.from.x - dx / size * zoom,
            y: this.server.from.y - dy / size * zoom
        };
    }

    handleEnd(ev: T) {
        ev.currentTarget.classList.remove(styles.dragging);
    }
}

export class MouseDragHandler extends DragHandler<React.MouseEvent> {

    protected getClientXY(ev: React.MouseEvent) {
        return {x: ev.clientX, y: ev.clientY};
    }

    shouldStart(): boolean {
        return true;
    }

    shouldEnd(): boolean {
        return true;
    }
}

export class TouchDragHandler extends DragHandler<React.TouchEvent> {

    protected getClientXY(ev: React.TouchEvent) {
        let touch = ev.targetTouches.item(0);
        return {x: touch.clientX, y: touch.clientY};
    }

    shouldStart(ev: React.TouchEvent) {
        return ev.targetTouches.length == 1;
    }

    shouldEnd(ev: React.TouchEvent) {
        return ev.targetTouches.length == 1;
    }
}
