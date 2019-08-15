import * as PIXI from "pixi.js";
import { AnimatedSpriteModel } from "./sprite";
import { Monitorable } from "../interfaces/monitor";
import { MenuStatus } from "../../state";
import { ApplicationProperty } from "../interfaces/pixi";
import { GraphicsAnimationGenerator } from "./animate";
import { Point } from "../interfaces/gamemap";

const graphicsOpts = {
    padding: 20,
    width: 1,
    alpha: 0.2,
    color: 0x607d8B,
    radius: 20
};

const defaultValues: {
    client: Point,
    menu: MenuStatus,
    enable: boolean
} = {
    client: {x: 0, y: 0},
    menu: MenuStatus.IDLE,
    enable: false
};

export default class extends AnimatedSpriteModel implements Monitorable {

    constructor(options: ApplicationProperty & { offset: number } ) { 
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
        ["x", "y", "menu"].forEach(key => 
            this.addUpdateCallback(key,
                () => this.merge("visible", this.isVisible())));
    }

    beforeRender() {
        super.beforeRender();
        this.sprite.x = this.props.x;
        this.sprite.y = this.props.y;
    }

    protected isVisible() {
        return this.props.menu === MenuStatus.SEEK_DEPARTURE
            && this.props.x != -1 && this.props.y != -1;
    }
}