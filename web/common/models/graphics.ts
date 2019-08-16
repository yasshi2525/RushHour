import * as PIXI from "pixi.js";
import { ApplicationProperty } from "../interfaces/pixi";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { PointModel, PointContainer } from "./point"



export abstract class GraphicsModel extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;

    constructor(options: ApplicationProperty) {
        super(options);
        this.graphics = new PIXI.Graphics();
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.container.addChild(this.graphics);
        });
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.container.removeChild(this.graphics));
    }

    beforeRender() {
        super.beforeRender();
        this.graphics.x = this.current.x;
        this.graphics.y = this.current.y;
    }    
}

export abstract class GraphicsContainer<T extends GraphicsModel> extends PointContainer<T> implements MonitorContrainer {

}
