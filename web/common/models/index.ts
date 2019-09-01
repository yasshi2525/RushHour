import * as PIXI from "pixi.js";
import { MenuStatus } from "../../state";
import { Coordinates, config } from "../interfaces/gamemap";
import { GameModelProperty } from "../interfaces/pixi";
import GameMapContainer from "./map";
import Controllers from "./controller";

export default class {
    app: PIXI.Application;
    renderer: PIXI.Renderer;
    gamemap: GameMapContainer;
    controllers: Controllers;
    timestamp: number;
    coord: Coordinates;
    delegate: number;
    offset: number;
    menu: MenuStatus;
    debugText: PIXI.Text | undefined;
    debugValue: any;

    constructor(options: GameModelProperty) {
        this.app = options.app;
        this.renderer = options.app.renderer;
       
        this.coord = { cx: options.cx, cy: options.cy, scale: options.scale, zoom: options.zoom };
        this.timestamp = 0;
        this.offset = 0;
        this.delegate = this.getDelegate();

        this.menu = MenuStatus.IDLE;

        this.gamemap = new GameMapContainer(this);
        this.controllers = new Controllers(this);
    }

    init() {
        this.gamemap.init();
        this.controllers.init();
        
        this.initDebugText();
        
        this.app.ticker.add(() => this.tick());    
    }

    protected tick() {
        this.offset++;
        if (this.offset >= config.round) {
            this.offset = 0;
        }
        this.gamemap.tick();
        this.controllers.tick();
    }

    protected initDebugText() {
        this.debugText = new PIXI.Text("");
        this.debugText.style.fontSize = 14;
        this.debugText.style.fill = 0xffffff;
        this.debugText.x = 50;
        this.app.stage.addChild(this.debugText);
        setInterval(() => this.viewDebugInfo(), 250);
    }

    protected viewDebugInfo() {
        if (this.debugText === undefined) {
            return;
        }
        this.debugText.text = "FPS: " + this.app.ticker.FPS.toFixed(2)
                                + ", " + this.app.stage.children.length + " entities"
                                + ", debug=" + this.debugValue 
                                + ", type=" + this.app.renderer.type 
                                + ", menu=" + this.menu;
    }

    setCoord(x: number, y: number, scale: number, force: boolean = false) {
        this.setCenter(x, y);
        this.setScale(scale);
        this.updateCoord(force);
    }

    protected getDelegate() {
        if (this.renderer.width < 600) { // sm
            return 2
        } else if (this.renderer.width < 960) { // md
            return 3
        } else if (this.renderer.width < 1280 ) { // lg
            return 3
        } else { // xl
            return 4
        }
    }

    protected updateDelegate() {
        let old = this.delegate;
        this.delegate = this.getDelegate();
        if (this.delegate !== old) {
            this.controllers.merge("delegate", this.delegate);
        }
    }

    protected setCenter(x: number, y: number) {
        let short = Math.min(this.renderer.width, this.renderer.height);
        let long = Math.max(this.renderer.width, this.renderer.height);
        let shortRadius = Math.pow(2, this.coord.scale - 1 + Math.log2(short/long));
        let longRadius = Math.pow(2, this.coord.scale - 1);

        if (this.renderer.width < this.renderer.height) {
            // 縦長
            if (x - shortRadius < config.gamePos.min.x) {
                x = config.gamePos.min.x + shortRadius;
            }
            if (x + shortRadius > config.gamePos.max.x) {
                x = config.gamePos.max.x - shortRadius;
            }
            if (y - longRadius < config.gamePos.min.y) {
                y = config.gamePos.min.y + longRadius;
            }
            if (y + longRadius > config.gamePos.max.y) {
                y = config.gamePos.max.y - longRadius;
            }
            if (this.coord.scale > config.scale.max) { 
                y = 0;
            }
        }else {
            // 横長
            if (x - longRadius < config.gamePos.min.x) {
                x = config.gamePos.min.x + longRadius;
            }
            if (x + longRadius > config.gamePos.max.x) {
                x = config.gamePos.max.x - longRadius;
            }
            if (y - shortRadius < config.gamePos.min.y) {
                y = config.gamePos.min.y + shortRadius;
            }
            if (y + shortRadius > config.gamePos.max.y) {
                y = config.gamePos.max.y - shortRadius;
            }
            if (this.coord.scale > config.scale.max) { 
                x = 0;
            }
        }
        
        this.coord.cx = x;
        this.coord.cy = y;
    }

    protected setScale(v: number) {
        let old = this.coord.scale

        let short = Math.min(this.renderer.width, this.renderer.height);
        let long = Math.max(this.renderer.width, this.renderer.height);
        let maxScale = config.scale.max + Math.log2(long/short);

        if (v < config.scale.min) {
            v = config.scale.min;
        }
        if (v > maxScale) {
            v = maxScale;
        }
        this.coord.zoom = v < old ? 1 : v > old ? -1 : 0;

        this.coord.scale = v;
    }

    resize(width: number, height: number) {
        let oldDelegate = this.delegate;
        this.renderer.resize(width, height);
        this.updateDelegate();
        this.controllers.merge("resize", true);
        this.gamemap.merge("resize", true);
        return this.delegate !== oldDelegate;
    }

    protected updateCoord(force: boolean) {
        this.controllers.merge("coord", this.coord);
        this.gamemap.merge("coord", this.coord);
        if (force) {
            this.controllers.merge("forceMove", true);
            this.gamemap.merge("forceMove", true);
        }
    }

    setMenuState(menu: MenuStatus) {
        if (this.menu !== menu) {
            this.controllers.merge("menu", menu);
            this.gamemap.merge("menu", menu);
            this.menu = menu;
        }
    }

    unmount() {
        this.gamemap.end();
        this.controllers.end();
    }
}
