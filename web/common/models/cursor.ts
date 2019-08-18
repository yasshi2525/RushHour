import * as PIXI from "pixi.js";
import { AnimatedSpriteModel } from "./sprite";
import { Monitorable } from "../interfaces/monitor";
import { MenuStatus, AnchorStatus } from "../../state";
import { ModelProperty } from "../interfaces/pixi";
import { GraphicsAnimationGenerator, RoundAnimationGenerator } from "./animate";
import { PointModel } from "./point";
import { RailNode } from "./rail";

const cursorOpts = {
    padding: 20,
    width: 1,
    alpha: 0.2,
    color: 0x607d8B,
    radius: 20
};

const cursorDefaultValues: {
    menu: MenuStatus
} = {
    menu: MenuStatus.IDLE
};

export class Cursor extends AnimatedSpriteModel implements Monitorable {
    selected: PointModel | undefined;
    anchor: Anchor;

    constructor(options: ModelProperty & { offset: number, anchor: Anchor } ) {
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(cursorOpts.width, cursorOpts.color);
        graphics.beginFill(cursorOpts.color, cursorOpts.alpha);
        graphics.drawCircle(
            cursorOpts.padding + cursorOpts.radius,
            cursorOpts.padding + cursorOpts.radius,
            cursorOpts.radius 
        );
        graphics.endFill();

        let generator = new GraphicsAnimationGenerator(options.app, graphics);

        let rect = graphics.getBounds().clone();
        rect.x -= cursorOpts.padding - 1;
        rect.y -= cursorOpts.padding - 1;
        rect.width += cursorOpts.padding * 2;
        rect.height += cursorOpts.padding * 2;

        let animation = generator.record(rect);

        super({ animation, ...options });
        this.anchor = options.anchor;
    }
    
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(cursorDefaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        this.props.pos = undefined;
        this.updateDestination();
        this.moveDestination();
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("pos", () => {
            this.selectObject();
            this.moveDestination();
        });
        this.addUpdateCallback("coord", () => {
            this.selectObject();
            this.updateDestination();
        });
    }

    protected calcDestination() {
        return (this.selected === undefined) 
        ? this.toView(this.props.pos)
        : this.toView(this.selected.get("pos"));
    }


    beforeRender() {
        if (!this.isVisible()) {
            this.sprite.visible = false;
            return;
        }
        if (this.selected !== undefined && this.selected.current !== undefined) {
            this.sprite.visible = true;
            this.sprite.x = this.selected.current.x - 2;
            this.sprite.y = this.selected.current.y - 2;
            return;
        }
        super.beforeRender();
    }

    selectObject(except: PointModel | undefined = undefined) {
        if (this.props.pos === undefined) {
            this.unlinkSelected();
            return;
        }
        var selected;
        switch(this.props.menu) {
            case MenuStatus.SEEK_DEPARTURE:
            case MenuStatus.EXTEND_RAIL:
                selected = this.model.getOnChunk("rail_nodes", this.props.pos, 1);
                break;
        }
        if (selected instanceof PointModel && selected !== this.anchor.object && selected !== except) {
            this.selected = selected;
            selected.refferedCursor = this;
            this.updateDestination();
        } else {
            this.unlinkSelected();
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

const anchorOpts = {
    padding: 20,
    width: 4,
    alpha: 1.0,
    slice: 8,
    color: 0x607d8B,
    radius: 20
};

const anchorDefaultValues: {
    menu: MenuStatus,
    oid: number,
    anchor: AnchorStatus | undefined
} = {
    menu: MenuStatus.IDLE,
    oid: 1,
    anchor: { type: "", pos: {x: 0, y: 0}, cid: 0 }
};

export class Anchor extends AnimatedSpriteModel implements Monitorable {
    object: PointModel | undefined;

    constructor(options: ModelProperty & { offset: number } ) { 
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(anchorOpts.width, anchorOpts.color, anchorOpts.alpha);

        let offset = anchorOpts.padding + anchorOpts.radius;

        for (var i = 0; i < anchorOpts.slice; i++) {
            let start = i / anchorOpts.slice * Math.PI * 2;
            let end = (i + 0.5) / anchorOpts.slice * Math.PI * 2;
            let next = (i + 1) / anchorOpts.slice * Math.PI * 2;

            graphics.lineStyle(anchorOpts.width, anchorOpts.color, anchorOpts.alpha);
            graphics.arc(offset, offset, anchorOpts.radius, start, end);
            graphics.lineStyle(anchorOpts.width, anchorOpts.color, 0);
            graphics.arc(offset, offset, anchorOpts.radius, end, next);
        }

        let generator = new RoundAnimationGenerator(options.app, graphics, new PIXI.Point(offset, offset));

        let rect = graphics.getBounds().clone();
        rect.x -= anchorOpts.padding - 1;
        rect.y -= anchorOpts.padding - 1;
        rect.width += anchorOpts.padding * 2;
        rect.height += anchorOpts.padding * 2;

        let animation = generator.record(rect);

        super({ animation, ...options });
        this.object = undefined;
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(anchorDefaultValues);
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
        this.addUpdateCallback("coord", () => this.updateAnchor());
        this.addUpdateCallback("anchor", () => this.updateAnchor());
    }

    updateAnchor() {
        if (this.props.anchor !== undefined) {
            this.object = this.model.getOnChunk(this.props.anchor.type, this.props.anchor.pos, this.props.oid);
        } else {
            this.object = undefined;
        }
    }

    beforeRender() {
        if (this.object === undefined || this.object.current === undefined) {
            this.sprite.visible = false;
        } else {
            this.sprite.visible = true;
            this.sprite.x = this.object.current.x - 5;
            this.sprite.y = this.object.current.y - 5;
        }
    }
}