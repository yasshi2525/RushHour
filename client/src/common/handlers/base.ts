import { Edge } from "interfaces/gamemap";
import { fetchMap } from "actions";
import GameModel from "models";

export default abstract class<T> {
  protected isExec: boolean;
  protected client: Edge;
  protected server: Edge;
  protected scale: { from: number; to: number };
  protected model: GameModel;
  protected dispatch: any;
  protected forceMove: boolean;

  constructor(model: GameModel, dispatch: any) {
    this.model = model;
    this.isExec = false;
    this.client = { from: { x: 0, y: 0 }, to: { x: 0, y: 0 } };
    this.server = {
      from: { x: model.coord.cx, y: model.coord.cy },
      to: { x: model.coord.cx, y: model.coord.cy }
    };
    this.scale = { from: model.coord.scale, to: model.coord.scale };
    this.dispatch = dispatch;
    this.forceMove = false;
  }

  protected shouldStart(_: T) {
    return true;
  }

  onStart(ev: T) {
    if (this.shouldStart(ev)) {
      this.isExec = true;
      this.server.from = { x: this.model.coord.cx, y: this.model.coord.cy };
      this.server.to = { x: this.model.coord.cx, y: this.model.coord.cy };
      this.scale.from = this.model.coord.scale;
      this.scale.to = this.model.coord.scale;
      this.handleStart(ev);
    }
  }

  protected abstract handleStart(ev: T): void;

  onMove(ev: T) {
    if (this.shouldMove(ev)) {
      this.handleMove(ev);

      this.model.setCoord(
        this.server.to.x,
        this.server.to.y,
        this.scale.to,
        this.forceMove
      );
    }
  }

  protected abstract handleMove(ev: T): void;

  protected shouldMove(_: T) {
    return this.isExec;
  }

  protected shouldEnd(_: T) {
    return true;
  }

  protected shouldFetch(_: T) {
    return true;
  }

  onEnd(ev: T) {
    if (this.shouldEnd(ev)) {
      this.isExec = false;
      this.handleEnd(ev);
      if (this.shouldFetch(ev)) {
        this.dispatch(fetchMap.request({ model: this.model }));
      }
    }
  }

  protected abstract handleEnd(ev: T): void;
}
