import { MenuStatus, AnchorStatus } from "../../state";
import { Monitorable } from "../interfaces/monitor";
import { AnimatedSpriteProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel } from "./sprite";
import { PointModel } from "./point";
import Cursor from "./cursor";

const defaultValues: {
    menu: MenuStatus,
    anchor: AnchorStatus | undefined
} = {
    menu: MenuStatus.IDLE,
    anchor: { type: "", pos: {x: 0, y: 0}, cid: 0 }
};

export default class extends AnimatedSpriteModel implements Monitorable {
    object: PointModel | undefined;
    cursor: Cursor | undefined;

    constructor(options: AnimatedSpriteProperty & { offset: number } ) { 
        super(options);
        this.object = undefined;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        this.props.anchor = undefined;
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("menu", (v: MenuStatus) => {
            switch (v) {
                case MenuStatus.IDLE:
                    this.merge("anchor", undefined);
            }
        })
        this.addUpdateCallback("coord", () => this.updateAnchor(false));
        this.addUpdateCallback("anchor", () => this.updateAnchor(true));
    }

    updateAnchor(force: boolean) {
        if (!force && this.object !== undefined) {
            return
        } 

        if (this.props.anchor !== undefined) {
            if (this.object !== undefined) {
                this.object.refferedAnchor = undefined;
            }
            this.object = this.model.gamemap.getOnChunk(this.props.anchor.type, this.props.anchor.pos, this.model.myid) as PointModel | undefined;
            
            if (this.object !== undefined) {
                this.object.refferedAnchor = this;
            }
            if (this.cursor !== undefined) {
                this.cursor.selectObject();
                this.cursor.moveDestination();
            }
        } else {
            this.object = undefined;
        }
        this.model.gamemap.merge("anchorObj", this.object);
    }

    updateDisplayInfo() {
        this.followPointModel(this.object);
    }
}