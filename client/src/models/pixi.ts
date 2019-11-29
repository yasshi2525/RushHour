import * as PIXI from "pixi.js";
import { config, Coordinates, Point } from "interfaces/gamemap";
import { Monitorable, MonitorContainer } from "interfaces/monitor";
import { PIXIProperty } from "interfaces/pixi";
import BaseContainer from "./container";
import BaseModel from "./base";

const defaultValues: { coord: Coordinates; [index: string]: any } = {
  coord: {
    cx: config.gamePos.default.x,
    cy: config.gamePos.default.y,
    scale: config.scale.default,
    zoom: 0
  },
  resize: false,
  forceMove: false,
  outMap: false,
  visible: true
};

export abstract class PIXIModel extends BaseModel implements Monitorable {
  protected app: PIXI.Application;
  protected parent: PIXI.Container;
  protected container: PIXI.Container;
  /**
   * smoothMove後、描画する座標(クライアント座標系)
   */
  destination: Point | undefined;
  /**
   * 描画する座標(クライアント座標系)
   */
  current: Point | undefined;
  /**
   * (x, y)が変化したとき、destination に移動するまでの残りフレーム数。
   */
  protected latency: number;

  protected smoothMoveFn: () => void;

  protected zIndex: number;

  constructor(options: PIXIProperty) {
    super(options);
    this.app = options.app;
    this.parent = options.container;
    this.container = new PIXI.Container();
    this.destination = { x: 0, y: 0 };
    this.current = { x: 0, y: 0 };
    this.latency = 0;
    this.zIndex = options.zIndex;
    this.smoothMoveFn = () => this.smoothMove();
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  setInitialValues(props: { [index: string]: any }) {
    super.setInitialValues(props);
    this.updateDestination();
    this.moveDestination();
  }

  setupBeforeCallback() {
    super.setupBeforeCallback();
    this.addBeforeCallback(() => {
      this.app.stage.addChild(this.container);
      this.app.ticker.add(this.smoothMoveFn);
      this.container.visible = this.props.visible;
      this.container.zIndex = this.zIndex;
    });
  }

  setupUpdateCallback() {
    super.setupUpdateCallback();
    this.addUpdateCallback("coord", () => {
      if (!this.props.deamon) {
        this.updateDestination();
      }
    });
    this.addUpdateCallback("forceMove", (v: boolean) => {
      if (v) {
        this.moveDestination();
      }
    });
    this.addUpdateCallback("resize", (v: boolean) => {
      if (v) {
        this.updateDestination();
        this.props.resize = false;
      }
    });
    this.addUpdateCallback("visible", v => {
      this.container.visible = v;
    });
  }

  setupAfterCallback() {
    super.setupAfterCallback();
    this.addAfterCallback(() => {
      this.app.ticker.remove(this.smoothMoveFn);
      this.app.stage.removeChild(this.container);
    });
  }

  toView(pos: Point | undefined): Point | undefined {
    if (pos === undefined) {
      return undefined;
    }
    let center = {
      x: this.app.renderer.width / this.app.renderer.resolution / 2,
      y: this.app.renderer.height / this.app.renderer.resolution / 2
    };
    let size = Math.max(
      this.app.renderer.width / this.app.renderer.resolution,
      this.app.renderer.height / this.app.renderer.resolution
    );
    let zoom = Math.pow(2, -this.props.coord.scale);

    return {
      x: (pos.x - this.props.coord.cx) * size * zoom + center.x,
      y: (pos.y - this.props.coord.cy) * size * zoom + center.y
    };
  }

  toServer(client: Point | undefined, offset: number = 0) {
    if (client === undefined) {
      return undefined;
    }
    let w = this.model.renderer.width;
    let h = this.model.renderer.height;
    let size = Math.max(this.model.renderer.width, this.model.renderer.height);

    let d = {
      x: (client.x + offset - w / 2) / size,
      y: (client.y + offset - h / 2) / size
    };

    let zoom = Math.pow(2, this.model.coord.scale);
    return {
      x: this.model.coord.cx + d.x * zoom,
      y: this.model.coord.cy + d.y * zoom
    };
  }

  shouldEnd() {
    return this.props.outMap && this.current == this.destination;
  }

  protected calcDestination() {
    return this.toView(this.props.pos);
  }

  updateDestination() {
    this.destination = this.calcDestination();
    this.latency = config.latency;
  }

  moveDestination() {
    this.current = this.destination;
    this.latency = 0;
    this.props.forceMove = false;
  }

  protected smoothMove() {
    if (this.latency > 0) {
      let ratio = this.latency / config.latency;
      if (ratio < 0.5) {
        ratio = 1.0 - ratio;
      }
      this.mapRatioToVariable(ratio);
      this.latency--;
    } else {
      this.moveDestination();
    }
    this.updateDisplayInfo();
  }

  protected mapRatioToVariable(ratio: number) {
    if (this.current !== undefined && this.destination !== undefined) {
      this.current.x =
        this.current.x * ratio + this.destination.x * (1 - ratio);
      this.current.y =
        this.current.y * ratio + this.destination.y * (1 - ratio);
    }
  }

  protected abstract getPIXIObject(): PIXI.DisplayObject;
}

export abstract class PIXIContainer<T extends PIXIModel, C extends PIXIProperty>
  extends BaseContainer<T, C>
  implements MonitorContainer {
  protected app: PIXI.Application;
  protected zIndex: number;
  constructor(options: PIXIProperty, newInstance: { new (props: C): T }) {
    super(options, newInstance);
    this.app = options.app;
    this.zIndex = options.zIndex;
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  protected getBasicChildOptions(): PIXIProperty {
    return {
      ...super.getBasicChildOptions(),
      app: this.app,
      zIndex: this.zIndex
    };
  }
}
