import * as PIXI from "pixi.js";
import { ResidenceContainer, CompanyContainer } from "./models/background";
import MonitorContainer from "./models/container";
import { Coordinates, config, getChunk, Point } from "./interfaces/gamemap";
import { Monitorable } from "./interfaces/monitor";
import { GameModelProperty, ResourceAttachable } from "./interfaces/pixi";
import { GameMap, MenuStatus } from "../state";
import { RailEdge, RailNodeContainer, RailEdgeContainer, RailNode } from "./models/rail";
import { StationContainer } from "./models/station";
import { Cursor, Anchor } from "./models/cursor";
import { WorldBorder, XBorderContainer, YBorderContainer } from "./models/border";
import { PointModel } from "./models/point";

const forceMove = { forceMove: true };
const resize = { resize: true };

export default class implements ResourceAttachable {
    protected app: PIXI.Application;
    renderer: PIXI.Renderer;
    protected xborder: XBorderContainer;
    protected yborder: YBorderContainer;
    protected world: WorldBorder;
    protected payload: {[index:string]: MonitorContainer<Monitorable>} = {};
    protected changed: boolean = false;
    cursor: Cursor;
    anchor: Anchor;
    timestamp: number;
    textures: {[index: string]: PIXI.Texture};
    coord: Coordinates;
    delegate: number;
    offset: number;
    menu: MenuStatus;
    debugText: PIXI.Text;
    debugValue: any;

    constructor(options: GameModelProperty) {
        this.app = options.app;
        this.renderer = options.app.renderer;
        this.textures = {};
       
        this.coord = { cx: options.cx, cy: options.cy, scale: options.scale, zoom: options.zoom };
        this.timestamp = 0;
        this.offset = 0;
        this.delegate = this.getDelegate();

        this.menu = MenuStatus.IDLE;
        this.anchor = new Anchor({ model: this, app: this.app, offset: this.offset });
        this.cursor = new Cursor({ model: this, app: this.app, offset: this.offset, anchor: this.anchor });

        [this.cursor, this.anchor].forEach((v: Monitorable) => {
            v.setupDefaultValues();
            v.setupUpdateCallback();
            v.setupBeforeCallback();
            v.setupAfterCallback();
            v.setInitialValues({});
            v.begin();
        });

        this.xborder = new XBorderContainer({ model: this, app: this.app, delegate: this.delegate });
        this.yborder = new YBorderContainer({ model: this, app: this.app, delegate: this.delegate });
        this.world = new WorldBorder({ model: this, app: this.app });

        [this.xborder, this.yborder, this.world].forEach((v: Monitorable) => {         
            v.setupDefaultValues();
            v.setupUpdateCallback();
            v.setupBeforeCallback();
            v.setupAfterCallback();
            v.setInitialValues({});
            v.begin();
        })

        this.app.ticker.add(() => {
            this.offset++;
            if (this.offset >= config.round) {
                this.offset = 0;
            }
            [this.cursor, this.anchor, this.xborder, this.yborder, this.world].forEach((v: Monitorable) => {
                v.merge("offset", this.offset);
                v.beforeRender();
            })
            Object.keys(this.payload).forEach(key => {
                this.payload[key].merge("offset", this.offset);
                this.payload[key].endChildren();
            });
        });

        this.debugText = new PIXI.Text("");
        this.debugText.style.fontSize = 14;
        this.debugText.style.fill = 0xffffff;
        this.debugText.x = 50;
        this.app.stage.addChild(this.debugText);
        setInterval(() => this.viewDebugInfo(), 250);
    }

    attach(textures: {[index: string]: PIXI.Texture}) {
        this.payload["residences"] = new ResidenceContainer({ model: this, app: this.app, texture: textures.residence});
        this.payload["companies"] = new CompanyContainer({ model: this, app: this.app, texture: textures.company});
        this.payload["stations"] = new StationContainer({ model: this, app: this.app, texture: textures.station});
        this.payload["rail_nodes"] = new RailNodeContainer({ model: this,app: this.app});
        this.payload["rail_edges"] = new RailEdgeContainer({ model: this, app: this.app});

        Object.keys(this.payload).forEach(key => {
            this.payload[key].setupDefaultValues();
            this.payload[key].setupUpdateCallback();
            this.payload[key].setupBeforeCallback();
            this.payload[key].setupAfterCallback();
            this.payload[key].begin();
        });
    }

    protected viewDebugInfo() {
        this.debugText.text = "FPS: " + this.app.ticker.FPS.toFixed(2)
                                + ", " + this.app.stage.children.length + " entities"
                                + ", debug=" + this.debugValue + ", type=" + this.app.renderer.type;
    }

    /**
     * 指定した id に対応するリソースを取得します
     * @param key リソース型
     * @param id id
     */
    get(key: string, id: string) {
        let container = this.payload[key];
        if (container !== undefined) {
            return container.getChild(id);
        }
        return undefined;
    }

    getOnChunk(key: string, pos: Point | undefined, oid: number): PointModel | undefined {
        if (this.payload[key] === undefined || pos === undefined) {
            return undefined;
        }
        return this.payload[key].getChildOnChunk(getChunk(pos, this.coord.scale - this.delegate + 1), oid)
    }

    merge(key: string, props: {id: string}) {
        return this.payload[key].mergeChild(props);
    }

    mergeAll(payload: GameMap) {
        config.zIndices.forEach(key => {
            if (this.payload[key] !== undefined) {
                this.payload[key].mergeChildren(payload[key], {coord: this.coord});
                if (this.payload[key].isChanged()) {
                    this.changed = true;
                }
            }
        });
        this.resolve();
        this.anchor.updateAnchor();
    }

    resolve() {
        if (this.payload["rail_nodes"] !== undefined) {
            this.payload["rail_nodes"].forEachChild((rn : RailNode) => {
                rn.resolve(this.get("rail_nodes", rn.get("pid")))
            });
        }
        if (this.payload["rail_edges"] !== undefined) {
            this.payload["rail_edges"].forEachChild((re: RailEdge) => 
                re.resolve(
                    this.get("rail_nodes", re.get("from")),
                    this.get("rail_nodes", re.get("to")),
                    this.get("rail_edges", re.get("eid"))
                )
            );
        }
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
            [this.xborder, this.yborder].forEach((v: Monitorable) => v.merge("delegate", this.delegate));
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
        [this.xborder, this.yborder, this.world].forEach((v: Monitorable) => v.mergeAll(resize));
        Object.keys(this.payload).forEach(key => this.payload[key].mergeAll(resize));
        return this.delegate !== oldDelegate;
    }

    protected updateCoord(force: boolean) {
        [this.cursor, this.anchor, this.xborder, this.yborder, this.world].forEach((v: Monitorable) => v.merge("coord", this.coord));
        if (force) {
            [this.cursor, this.anchor, this.xborder, this.yborder, this.world].forEach((v: Monitorable) => v.mergeAll(forceMove));
        }
        Object.keys(this.payload).forEach(key => {
            this.payload[key].merge("coord", this.coord);
            if (force) {
                this.payload[key].mergeAll(forceMove);
            }
            if (this.payload[key].isChanged()) {
                this.changed = true;
            }
        });
    }

    setMenuState(menu: MenuStatus) {
        if (this.menu !== menu) {
            [this.cursor, this.anchor].forEach(v => {
                v.merge("menu", menu);
            })
            this.menu = menu;
        }
    }

    isChanged() {
        return this.changed;
    }

    render() {
        Object.keys(this.payload).forEach(key => 
            this.payload[key].forEachChild((c) => c.beforeRender())
        );
        Object.keys(this.payload).forEach(key => this.payload[key].reset());
        this.changed = false;
    }

    unmount() {
        Object.keys(this.payload).reverse().forEach(key => this.payload[key].end());
        [this.cursor, this.anchor, this.xborder, this.yborder, this.world].forEach((v: Monitorable) => v.end());
    }
}
