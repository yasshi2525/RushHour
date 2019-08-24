import * as PIXI from "pixi.js";
import { MenuStatus, AnchorStatus } from "../../state";
import { Monitorable } from "../interfaces/monitor";
import { PIXIProperty } from "../interfaces/pixi";
import { AnimatedSpriteModel } from "./sprite";
import { RoundAnimationGenerator } from "./animate";
import { PointModel } from "./point";

const graphicsOpts = {
    padding: 20,
    width: 4,
    alpha: 1.0,
    slice: 8,
    color: 0x607d8B,
    radius: 20
};

const defaultValues: {
    menu: MenuStatus,
    oid: number,
    anchor: AnchorStatus | undefined
} = {
    menu: MenuStatus.IDLE,
    oid: 1,
    anchor: { type: "", pos: {x: 0, y: 0}, cid: 0 }
};

export default class extends AnimatedSpriteModel implements Monitorable {
    object: PointModel | undefined;

    constructor(options: PIXIProperty & { offset: number } ) { 
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(graphicsOpts.width, graphicsOpts.color, graphicsOpts.alpha);

        let offset = graphicsOpts.padding + graphicsOpts.radius;

        for (var i = 0; i < graphicsOpts.slice; i++) {
            let start = i / graphicsOpts.slice * Math.PI * 2;
            let end = (i + 0.5) / graphicsOpts.slice * Math.PI * 2;
            let next = (i + 1) / graphicsOpts.slice * Math.PI * 2;

            graphics.lineStyle(graphicsOpts.width, graphicsOpts.color, graphicsOpts.alpha);
            graphics.arc(offset, offset, graphicsOpts.radius, start, end);
            graphics.lineStyle(graphicsOpts.width, graphicsOpts.color, 0);
            graphics.arc(offset, offset, graphicsOpts.radius, end, next);
        }

        let generator = new RoundAnimationGenerator(options.app, graphics, new PIXI.Point(offset, offset));

        let rect = graphics.getBounds().clone();
        rect.x -= graphicsOpts.padding - 1;
        rect.y -= graphicsOpts.padding - 1;
        rect.width += graphicsOpts.padding * 2;
        rect.height += graphicsOpts.padding * 2;

        let animation = generator.record(rect);

        super({ animation, ...options });
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
            this.object = this.model.gamemap.getOnChunk(this.props.anchor.type, this.props.anchor.pos, this.props.oid) as PointModel;
            this.object.refferedAnchor = this;
        } else {
            this.object = undefined;
        }
    }

    updateDisplayInfo() {
        this.followPointModel(this.object, 0);
    }
}