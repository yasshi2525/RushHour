import * as PIXI from "pixi.js";
import { AnimatedSpriteModel } from "./sprite";
import { Monitorable } from "../interfaces/monitor";
import { MenuStatus } from "../../state";
import { ModelProperty } from "../interfaces/pixi";
import { GraphicsAnimationGenerator } from "./animate";
import { PointModel } from "./point";

const graphicsOpts = {
    padding: 20,
    width: 1,
    alpha: 0.2,
    color: 0x607d8B,
    radius: 20
};

const defaultValues: {
    menu: MenuStatus,
    enable: boolean
} = {
    menu: MenuStatus.IDLE,
    enable: false
};

export default class extends AnimatedSpriteModel implements Monitorable {
    selected: PointModel | undefined;

    constructor(options: ModelProperty & { offset: number } ) { 
        let graphics = new PIXI.Graphics();
        graphics.lineStyle(graphicsOpts.width, graphicsOpts.color);
        graphics.beginFill(graphicsOpts.color, graphicsOpts.alpha);
        graphics.drawCircle(
            graphicsOpts.padding + graphicsOpts.radius,
            graphicsOpts.padding + graphicsOpts.radius,
            graphicsOpts.radius 
        );
        graphics.endFill();

        let generator = new GraphicsAnimationGenerator(options.app, graphics);

        let rect = graphics.getBounds().clone();
        rect.x -= graphicsOpts.padding - 1;
        rect.y -= graphicsOpts.padding - 1;
        rect.width += graphicsOpts.padding * 2;
        rect.height += graphicsOpts.padding * 2;

        let animation = generator.record(rect);

        super({ animation, ...options });
    }
    
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("pos", () => {
            this.selectObject();
            this.moveDestination();
        });
        this.addUpdateCallback("coord", () => {
            this.selectObject();
            this.moveDestination();
        });
    }

    protected calcDestination() {
        return (this.selected === undefined) 
        ? this.toView(this.props.pos)
        : this.props.client;
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

    selectObject() {
        if (this.props.pos === undefined) {
            this.unlinkSelected();
            return;
        }
        var selected;
        switch(this.props.menu) {
            case MenuStatus.SEEK_DEPARTURE:
                selected = this.model.getOnChunk("rail_nodes", this.props.pos, 1);
                break;
        }
        if (selected instanceof PointModel) {
            this.selected = selected;
            selected.refferedCursor = this;
            this.updateDestination();
        } else {
            this.unlinkSelected();
        }
    }

    unlinkSelected() {
        this.selected = undefined;
        this.updateDestination();
    }

    protected isVisible() {
        return this.props.menu === MenuStatus.SEEK_DEPARTURE;
    }
}