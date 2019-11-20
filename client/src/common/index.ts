import * as PIXI from "pixi.js";
import GameModel from "./models";
import { config } from "./interfaces/gamemap";

export default class {
  app: PIXI.Application;
  sheets = [
    "cursor",
    "anchor",
    "residence",
    "company",
    "rail_node",
    "rail_edge",
    "destroy"
  ];
  model: GameModel;

  constructor(myid: number) {
    this.app = new PIXI.Application({
      width: window.innerWidth,
      height: window.innerHeight,
      backgroundColor: config.background,
      autoStart: true,
      antialias: true,
      resolution: window.devicePixelRatio,
      autoDensity: true
    });
    this.app.stage.sortableChildren = true;
    this.model = new GameModel({
      app: this.app,
      cx: config.gamePos.default.x,
      cy: config.gamePos.default.y,
      scale: config.scale.default,
      zoom: 0,
      myid
    });
    this.sheets.forEach(key => {
      this.app.loader.add(
        key,
        `assets/bundle/spritesheet/${key}@${Math.floor(
          this.model.renderer.resolution
        )}x.json`
      );
    });
  }

  init() {
    this.model.init();
  }
}
