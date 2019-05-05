import { config, Point } from "../interfaces/gamemap";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { ContainerProperty } from "../interfaces/pixi";
import { PIXIModel, PIXIContainer } from "./pixi";
const defaultValues: {x: number, y:number} = {x: 0, y: 0};

export abstract class PointModel extends PIXIModel implements Monitorable {
    /**
     * smoothMove後、描画する座標(クライアント座標系)
     */
    destination: Point;
    /**
     * 描画する座標(クライアント座標系)
     */
    current: Point;
    /**
     * (x, y)が変化したとき、destination に移動するまでの残り時間。
     */
    protected latency: number;

    protected frame: number;

    protected smoothMoveFn: () => void;

    constructor(options: ContainerProperty) {
        super(options);
        this.destination = {x: 0, y: 0};
        this.current = {x: 0, y: 0};
        this.latency = 0;
        this.frame = 1000 / this.app.ticker.FPS
        this.smoothMoveFn = () => this.smoothMove();
    }

    setInitialValues(initialValues: {[index: string]: {}}) {
        super.setInitialValues(initialValues);
        this.current = this.toView(this.props.x, this.props.y);
        this.destination = this.toView(this.props.x, this.props.y);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupBeforeCallback() {
        super.setupBeforeCallback();
        this.addBeforeCallback(() => {
            this.app.ticker.add(this.smoothMoveFn)
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            this.app.ticker.remove(this.smoothMoveFn);
        })
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        ["x", "y"].forEach(v => this.addUpdateCallback(v, () => this.updateDestination()));
    }

    updateDestination() {
        this.destination = this.toView(this.props.x, this.props.y);
        this.latency = config.latency;
    }

    protected smoothMove() {
        if (this.latency > this.frame) {
            this.current.x = (this.current.x + this.destination.x) / 2;
            this.current.y = (this.current.y + this.destination.y) / 2;
            this.latency -= this.app.ticker.elapsedMS;
        } else {
            this.current = this.destination;
            this.latency = 0;
        }
        this.beforeRender();
    }
}

export abstract class PointContainer<T extends PointModel> extends PIXIContainer<T> implements MonitorContrainer {    

    setupUpdateCallback() {
        super.setupUpdateCallback();
        ["cx", "cy", "scale"].forEach(v => this.addUpdateCallback(v, () => this.updateDestination()));
    }

    protected updateDestination() {
        this.forEachChild(c => c.updateDestination());
    }
}