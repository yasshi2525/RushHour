import { Monitorable, MonitorContainer } from "interfaces/monitor";
import { SpriteProperty } from "interfaces/pixi";
import { config, ResolveError } from "interfaces/gamemap";
import { PIXIProperty } from "interfaces/pixi";
import { PointModel, PointContainer } from "./point";

const defaultValues: {
  pid: number;
  cid: number;
  mul: number;
} = {
  pid: 0,
  cid: 0,
  mul: 1
};

export abstract class ZoomablePointModel extends PointModel
  implements Monitorable {
  parentModel: PointModel | undefined;

  constructor(options: SpriteProperty) {
    super(options);
    this.parentModel = undefined;
  }

  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  setupAfterCallback() {
    super.setupAfterCallback();
    this.addAfterCallback(() => {
      if (this.parentModel !== undefined) {
        this.parentModel.merge("visible", true);
      }
    });
  }

  getResourceName(): string | undefined {
    return undefined;
  }

  resolve(error: ResolveError) {
    let resourceName = this.getResourceName();
    if (resourceName !== undefined) {
      let parent = this.model.gamemap.get(resourceName, this.props.pid) as
        | PointModel
        | undefined;
      if (parent !== undefined) {
        this.parentModel = parent;
        // 拡大時、派生元の座標から移動を開始する
        if (this.props.coord.zoom == 1) {
          this.current = Object.assign({}, parent.current);
          this.latency = config.latency;
        }
        // 縮小時、集約先の座標に向かって移動する
        if (this.props.coord.zoom == -1) {
          this.merge("pos", parent.get("pos"));
        }
        parent.merge("visible", false);
      }
    }
    return error;
  }
}

export abstract class ZoomablePointModelContainer<
  T extends ZoomablePointModel,
  C extends PIXIProperty
> extends PointContainer<T, C> implements MonitorContainer {}
