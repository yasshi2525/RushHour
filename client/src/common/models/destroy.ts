import { Monitorable } from "../interfaces/monitor";
import { Point } from "../interfaces/gamemap";
import { AnimatedSpriteProperty } from "../interfaces/pixi";
import { MenuStatus } from "../../state";
import { AnimatedSpriteModel } from "./sprite";
import { PointModel } from "./point";

const graphicsOpts = {
    tint: 0xf44336
}

const defaultOpts: {
    menu: MenuStatus,
    client: Point,
    enabled: boolean
 } = {
    menu: MenuStatus.IDLE,
    client: {x: 0, y: 0},
    enabled: false,
}

export default class extends AnimatedSpriteModel implements Monitorable {
    selected: PointModel | undefined;

    constructor(options: AnimatedSpriteProperty) {
        super(options);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultOpts);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        this.props.client = undefined;
        this.props.pos = undefined;
        this.updateDestination();
        this.moveDestination();
        this.merge("tint", graphicsOpts.tint);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("client", (v) => {
            this.merge("pos", this.toServer(v));
            this.selectObject();
            this.moveDestination();
        });
        this.addUpdateCallback("coord", () => {
            this.merge("pos", this.toServer(this.props.client))
            this.selectObject();
            this.updateDestination();
        });
        this.addUpdateCallback("menu", (v: MenuStatus) => {
            this.merge("enabled", v === MenuStatus.DESTROY);
        });
    }

    protected calcDestination() {
        return (this.selected !== undefined)
        ? this.toView(this.selected.get("pos"))
        : this.toView(this.toServer(this.props.client));
    }

    updateDisplayInfo() {
        if (!this.props.enabled) {
            this.sprite.visible = false;
            return;
        }
        if (!this.followPointModel(this.selected)) {
            super.updateDisplayInfo();
        }
    }

    selectObject(except: PointModel | undefined = undefined) {
        let objOnChunk = this.getObjectOnChunk(except);
        if (objOnChunk === undefined) {
            this.unlinkSelected();
        } else if (this.selected !== objOnChunk) {
            this.selected = objOnChunk;
            objOnChunk.refferedDestroyer = this;
            this.updateDestination();
        } 
    }

    unlinkSelected() {
        if (this.selected !== undefined) {
            this.selected = undefined;
            this.updateDestination();
        }
    }

    protected getObjectOnChunk(except: Monitorable | undefined = undefined) {
        let selected = this.model.gamemap.getOnChunk("rail_nodes", this.props.pos, this.model.myid);

        return selected === except ? undefined : selected as PointModel | undefined;
    }
}