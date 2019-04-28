import * as styles from "./style.css"
import * as React from "react";
import GameModel from "../model";
import { Edge } from "../interfaces/gamemap";

export class DragHandler {
    isDrag: boolean;
    client: Edge;
    server: Edge;
    model: GameModel;

    constructor(model: GameModel) {
        this.model = model;
        this.isDrag = false;
        this.client = {from: {x: 0, y: 0}, to: {x: 0, y: 0}};
        this.server = {from: {x: 0, y: 0}, to: {x: 0, y: 0}};;
    }

    onDragStart(ev: React.MouseEvent) {
        this.isDrag = true;
        this.client.from = {x: ev.clientX, y: ev.clientY};
        this.server.from = {x: this.model.cx, y: this.model.cy};
        ev.currentTarget.classList.add(styles.dragging);
    }

    onDragMove(ev: React.MouseEvent) {
        if (this.isDrag) {
            this.client.to = {x: ev.clientX, y: ev.clientY};

            let dx = this.client.to.x - this.client.from.x;
            let dy = this.client.to.y - this.client.from.y;
            let size = Math.max(ev.currentTarget.clientWidth, ev.currentTarget.clientHeight);
            let zoom = Math.pow(2, this.model.scale);

            this.server.to = {
                x: this.server.from.x - dx / size * zoom,
                y: this.server.from.y - dy / size * zoom
            };

            this.model.setCenter(this.server.to.x, this.server.to.y);
            if (this.model.isChanged()) {
                this.model.render();
            }
        } 
    }

    onDragEnd(ev: React.MouseEvent) {
        this.isDrag = false;
        ev.currentTarget.classList.remove(styles.dragging);
    }
}