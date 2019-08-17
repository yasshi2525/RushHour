import * as PIXI from "pixi.js";
import { ModelProperty } from "../interfaces/pixi";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { PointModel, PointContainer } from "./point"



export abstract class GraphicsModel extends PointModel implements Monitorable {
    protected graphics: PIXI.Graphics;

    constructor(options: ModelProperty) {
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
        if (this.current === undefined) {
            this.graphics.visible = false;
        } else {
            this.graphics.visible = true;
            this.graphics.x = this.current.x;
            this.graphics.y = this.current.y;
        }
    }    
}

export abstract class GraphicsContainer<T extends GraphicsModel> extends PointContainer<T> implements MonitorContrainer {

}
