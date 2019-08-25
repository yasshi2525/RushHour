import * as PIXI from "pixi.js";
import { MenuStatus } from "../../state";
import { Monitorable } from "../interfaces/monitor";
import { Point } from "../interfaces/gamemap";
import { PIXIProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel } from "./sprite";
import { GraphicsAnimationGenerator } from "./animate";
import { PointModel } from "./point";
import { RailNode } from "./rail";
import Anchor from "./anchor";

const graphicsOpts = {
    padding: 20,
    width: 1,
    alpha: 0.2,
    color: 0x607d8B,
    offset: { server: 3, client: 2 },
    tint: {
        info: 0xffffff,
        error: 0xf44336,
    },
    radius: 20
};

const defaultValues: {
    menu: MenuStatus,
    client: Point
} = {
    menu: MenuStatus.IDLE,
    client: {x: 0, y: 0}
};

export default class extends AnimatedSpriteModel implements Monitorable {
    selected: PointModel | undefined;
    anchor: Anchor;

    constructor(options: PIXIProperty & { offset: number, anchor: Anchor } ) {
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(graphicsOpts.width, graphicsOpts.color);
        graphics.beginFill(graphicsOpts.color, graphicsOpts.alpha);
        graphics.drawCircle(
            graphicsOpts.padding + graphicsOpts.radius,
            graphicsOpts.padding + graphicsOpts.radius,
            graphicsOpts.radius 
        );
        graphics.endFill();
        graphics.tint = graphicsOpts.tint.info;

        let generator = new GraphicsAnimationGenerator(options.app, graphics);

        let rect = graphics.getBounds().clone();
        rect.x -= graphicsOpts.padding - 1;
        rect.y -= graphicsOpts.padding - 1;
        rect.width += graphicsOpts.padding * 2;
        rect.height += graphicsOpts.padding * 2;

        let animation = generator.record(rect);

        super({ animation, ...options });
        this.anchor = options.anchor;
        this.anchor.cursor = this;
    }
    
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        this.props.client = undefined;
        this.props.pos = undefined;
        this.updateDestination();
        this.moveDestination();
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("client", (v) => {
            this.merge("pos", this.toServer(v, graphicsOpts.offset.client))
            this.selectObject();
            this.moveDestination();
        });
        this.addUpdateCallback("coord", () => {
            this.merge("pos", this.toServer(this.props.client, graphicsOpts.offset.client))
            this.selectObject();
            this.updateDestination();
        });
    }

    protected calcDestination() {
        return (this.selected !== undefined) 
        ? this.toView(this.selected.get("pos"))
        : this.toView(this.toServer(this.props.client, graphicsOpts.offset.client));
    }


    updateDisplayInfo() {
        if (!this.isVisible()) {
            this.sprite.visible = false;
            return;
        }
        if (!this.followPointModel(this.selected, graphicsOpts.offset.server)) {
            super.updateDisplayInfo();
        }
    }

    selectObject(except: PointModel | undefined = undefined) {
        let objOnChunk = this.getObjectOnChunk(except);
        let tint = this.getTint(objOnChunk);
        this.merge("tint", tint);

        if (objOnChunk === undefined) {
            this.unlinkSelected();
        } else if (objOnChunk !== this.anchor.object) {
            this.selected = objOnChunk;
            objOnChunk.refferedCursor = this;
            this.updateDestination();
        } else {
            this.unlinkSelected();
        }
        this.selectObjectCursor();
    }

    protected getObjectOnChunk(except: Monitorable | undefined = undefined) {
        let selected;
        switch(this.props.menu) {
            case MenuStatus.SEEK_DEPARTURE:
            case MenuStatus.EXTEND_RAIL:
                selected = this.model.gamemap.getOnChunk("rail_nodes", this.props.pos, 2);
                break;
        }
        return selected === except ? undefined : selected as PointModel;
    }

    protected getTint(objOnChunk: Monitorable | undefined) {
        switch(this.props.menu) {
            case MenuStatus.EXTEND_RAIL:
                if (this.anchor.object === objOnChunk) {
                    return graphicsOpts.tint.error
                }
                break;
        }
        return graphicsOpts.tint.info;
    }

    protected selectObjectCursor() {
        if (this.model.controllers.getCursor().selected === undefined) {
            let pos = this.props.client !== undefined ? {
                x: this.props.client.x / this.model.renderer.resolution,
                y: this.props.client.y / this.model.renderer.resolution
            } : undefined;
            this.model.gamemap.merge("cursorClient", pos);
        } else {
            this.model.gamemap.merge("cursorClient", undefined);
        }
    }

    unlinkSelected() {
        if (this.selected !== undefined) {
            this.selected = undefined;
            this.updateDestination();
        }
    }

    genAnchorStatus() {
        if (this.selected === undefined) {
            return undefined;
        } else {
            let res = { pos: this.selected.get("pos"), type: "", cid: this.selected.get("cid") };
            if (this.selected instanceof RailNode) {
                res.type = "rail_nodes"
                return res
            } 
            return undefined;
        }
    }

    protected isVisible() {
        return this.props.menu === MenuStatus.SEEK_DEPARTURE || this.props.menu === MenuStatus.EXTEND_RAIL;
    }
}