import * as PIXI from "pixi.js";
import { ZIndex } from "interfaces/pixi";
import { ResolveError } from "interfaces/gamemap";
import { Monitorable } from "interfaces/monitor";
import { Entity } from "interfaces";
import { FetchMapResponse, FetchMapResponseKeys } from "interfaces/endpoint";
import GroupModel from "./group";
import { ResidenceContainer, CompanyContainer } from "./background";
import { RailNodeContainer, RailEdgeContainer } from "./rail";
import { PlayerContainer } from "./player";

export default class extends GroupModel {
  init() {
    let textures = this.model.app.loader.resources;
    let base = {
      model: this.model,
      app: this.model.app
    };
    this.containers["players"] = new PlayerContainer(base);
    // this.containers["stations"] = new StationContainer({ ...base, zIndex: ZIndex.STATION, texture: textures["station"].texture});

    let anim: { [index: string]: PIXI.Texture[] } = {};

    ["residence", "company", "rail_node", "rail_edge"]
      .map(key => ({ key, sheet: textures[key].spritesheet }))
      .forEach(
        entry =>
          (anim[entry.key] =
            entry.sheet !== undefined
              ? entry.sheet.animations[entry.key]
              : [PIXI.Texture.EMPTY])
      );

    this.containers["residences"] = new ResidenceContainer({
      ...base,
      zIndex: ZIndex.RESIDENCE,
      animation: anim["residence"]
    });
    this.containers["companies"] = new CompanyContainer({
      ...base,
      zIndex: ZIndex.COMPANY,
      animation: anim["company"]
    });
    this.containers["rail_nodes"] = new RailNodeContainer({
      ...base,
      zIndex: ZIndex.RAIL_NODE,
      animation: anim["rail_node"]
    });
    this.containers["rail_edges"] = new RailEdgeContainer({
      ...base,
      zIndex: ZIndex.RAIL_EDGE,
      animation: anim["rail_edge"]
    });

    super.init();
  }

  mergeChild(key: string, props: { id: number }): undefined | Monitorable {
    if (this.containers[key] === undefined) {
      return undefined;
    }
    return this.containers[key].mergeChild(props);
  }

  mergeChildren(
    key: string,
    props: Entity[],
    opts: { [index: string]: any } = {}
  ) {
    if (this.containers[key] !== undefined) {
      this.containers[key].mergeChildren(props, opts);
      if (this.containers[key].isChanged()) {
        this.changed = true;
      }
    }
  }

  mergeAll(payload: FetchMapResponse) {
    FetchMapResponseKeys.forEach(key => {
      const list = Object.values(payload[key]);
      this.mergeChildren(key, list, { coord: this.model.coord });
    });
    let error = this.resolve();
    this.model.controllers.updateAnchor();
    return error;
  }

  removeChild(key: string, id: number) {
    this.containers[key].removeChild(id);
  }

  resolve() {
    let error: ResolveError = {};
    this.forEach(v => v.resolve(error));
    return error;
  }
}
