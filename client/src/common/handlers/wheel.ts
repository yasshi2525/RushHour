import * as React from "react";
import { PointHandler } from "./point";

const sensitivity = 0.1;

export class WheelHandler extends PointHandler<React.WheelEvent> {
  protected getClientXY(ev: React.WheelEvent) {
    return {
      x: ev.clientX * this.model.renderer.resolution,
      y: ev.clientY * this.model.renderer.resolution
    };
  }

  protected shouldStart() {
    return true;
  }

  protected handleStart() {}

  protected handleMove(ev: React.WheelEvent): void {
    let ds = ev.deltaY > 0 ? sensitivity : -sensitivity;
    this.scale.to = Math.round((this.scale.from + ds) * 10) / 10;
    let center = this.zoom(
      this.getClientXY(ev),
      this.scale.to - this.scale.from
    );
    this.server.to.x = center.x;
    this.server.to.y = center.y;
  }

  protected shouldEnd(): boolean {
    return true;
  }

  protected handleEnd() {}

  protected shouldFetch() {
    return Math.floor(this.scale.from) != Math.floor(this.scale.to);
  }
}
