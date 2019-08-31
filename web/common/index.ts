import * as PIXI from "pixi.js";
import GameModel from "./models";
import { config } from "./interfaces/gamemap";

export default class {
    app: PIXI.Application;
    protected images = ["residence", "company", "station", "train"];
    sheets = ["cursor", "anchor", "rail_node", "rail_edge"];
    model: GameModel;

    constructor() {
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
            zoom: 0
        });

        this.images.forEach(key => this.app.loader.add(key, `public/img/${key}.png`));
        
        this.sheets.forEach(key => {
            this.app.loader.add(key, `public/spritesheet/${key}@${this.model.renderer.resolution}x.json`);
        });
    }

    initModel() {
        this.model.init();
    }
}