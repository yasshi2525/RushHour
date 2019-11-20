import * as PIXI from "pixi.js";
import { ZIndex } from "common/interfaces/pixi";
import GroupModel from "./group";
import Anchor from "./anchor";
import Cursor from "./cursor";
import WorldBorder from "./world";
import { XBorderContainer, YBorderContainer } from "./border";
import Destroyer from "./destroy";

export default class extends GroupModel {
  init() {
    let textures = this.model.app.loader.resources;
    let props = {
      model: this.model,
      app: this.model.app,
      offset: this.model.offset,
      delegate: this.model.delegate,
      texture: PIXI.Texture.EMPTY
    };

    let anim: { [index: string]: PIXI.Texture[] } = {};

    ["anchor", "cursor", "destroy"]
      .map(key => ({ key, sheet: textures[key].spritesheet }))
      .forEach(
        entry =>
          (anim[entry.key] =
            entry.sheet !== undefined
              ? entry.sheet.animations[entry.key]
              : [PIXI.Texture.EMPTY])
      );

    let anchor = new Anchor({
      ...props,
      zIndex: ZIndex.ANCHOR,
      animation: anim["anchor"]
    });
    let destroyer = new Destroyer({
      ...props,
      zIndex: ZIndex.DESTROY,
      animation: anim["destroy"]
    });
    this.objects.cursor = new Cursor({
      ...props,
      anchor,
      destroyer,
      zIndex: ZIndex.CURSOR,
      animation: anim["cursor"]
    });
    this.objects.anchor = anchor;
    this.objects.destroyer = destroyer;
    this.containers.xborder = new XBorderContainer({
      ...props,
      zIndex: ZIndex.NORMAL_BORDER
    });
    this.containers.yborder = new YBorderContainer({
      ...props,
      zIndex: ZIndex.NORMAL_BORDER
    });
    this.objects.world = new WorldBorder({
      ...props,
      zIndex: ZIndex.WORLD_BORDER
    });

    super.init();
  }

  getCursor() {
    return this.objects.cursor as Cursor;
  }

  getAnchor() {
    return this.objects.anchor as Anchor;
  }

  updateAnchor() {
    (this.objects.anchor as Anchor).updateAnchor(false);
  }
}
