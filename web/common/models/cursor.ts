import { MenuStatus } from "../../state";
import { Monitorable } from "../interfaces/monitor";
import { Point } from "../interfaces/gamemap";
import { AnimatedSpriteProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel } from "./sprite";
import { PointModel } from "./point";
import { RailNode } from "./rail";
import Anchor from "./anchor";

const graphicsOpts = {
    tint: {
        info: 0xffffff,
        error: 0xf44336,
    }
};

const defaultValues: {
    menu: MenuStatus,
    client: Point,
    activation: boolean
} = {
    menu: MenuStatus.IDLE,
    client: {x: 0, y: 0},
    activation: true
};

export default class extends AnimatedSpriteModel implements Monitorable {
    selected: PointModel | undefined;
    anchor: Anchor;

    constructor(options: AnimatedSpriteProperty & { offset: number, anchor: Anchor } ) {
        super(options);
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
            this.merge("pos", this.toServer(v))
            this.selectObject();
            this.moveDestination();
        });
        this.addUpdateCallback("coord", () => {
            this.merge("pos", this.toServer(this.props.client))
            this.selectObject();
            this.updateDestination();
        });
    }

    protected calcDestination() {
        return (this.selected !== undefined) 
        ? this.toView(this.selected.get("pos"))
        : this.toView(this.toServer(this.props.client));
    }


    updateDisplayInfo() {
        if (!this.isVisible()) {
            this.sprite.visible = false;
            return;
        }
        if (!this.followPointModel(this.selected)) {
            super.updateDisplayInfo();
        }
    }

    selectObject(except: PointModel | undefined = undefined) {
        let objOnChunk = this.getObjectOnChunk(except);
        this.activate(objOnChunk);
        let tint = this.getTint();
        this.merge("tint", tint);

        if (objOnChunk === undefined) {
            this.unlinkSelected();
        } else if (objOnChunk !== this.anchor.object) {
            this.selected = objOnChunk;
            objOnChunk.refferedCursor = this;
            this.model.gamemap.merge("cursorObj", objOnChunk);
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
    protected activate(objOnChunk: Monitorable | undefined) {
        let activation = true;
        switch(this.props.menu) {
            case MenuStatus.EXTEND_RAIL:
                if (this.anchor.object === objOnChunk) {
                    activation = false;
                }
                if (this.anchor.object instanceof RailNode && this.selected instanceof RailNode) {
                    let anchor = this.anchor.object;
                    let selected = this.selected;
                    let link = Object.keys(anchor.out).find(eid => eid != "cursorIn" && anchor.out[eid].to == selected);
                    if (link !== undefined) {
                        activation = false;
                    }
                }
                break;
        }
        this.merge("activation", activation);
    }

    protected getTint() {
        return this.props.activation ? graphicsOpts.tint.info : graphicsOpts.tint.error
    }

    protected selectObjectCursor() {
        if (this.selected === undefined) {
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
            this.model.gamemap.merge("cursorObj", undefined);
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