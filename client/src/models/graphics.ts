import * as PIXI from "pixi.js";
import { PIXIProperty } from "interfaces/pixi";
import { Monitorable, MonitorContainer } from "interfaces/monitor";
import { PointModel, PointContainer } from "./point";

export abstract class GraphicsModel extends PointModel implements Monitorable {
  protected graphics: PIXI.Graphics;

  constructor(options: PIXIProperty) {
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

  updateDisplayInfo() {
    super.updateDisplayInfo();
    this.setDisplayPosition();
  }

  protected getPIXIObject() {
    return this.graphics;
  }
}

export abstract class GraphicsContainer<
  T extends GraphicsModel,
  C extends PIXIProperty
> extends PointContainer<T, C> implements MonitorContainer {}
