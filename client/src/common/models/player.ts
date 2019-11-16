import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import BaseModel from "./base";
import BaseContainer from "./container";
import { ModelProperty } from "../interfaces/pixi";
import { hueToRgb } from "../interfaces/gamemap";

const defaultValues: { color: number, hue: number } = { color: 0, hue: 0 };

export class Player extends BaseModel implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setInitialValues(props: {[index: string]: any}) {
        super.setInitialValues(props);
        let rgb = hueToRgb(props.hue);
        this.props.color = (rgb[0] << 16) + (rgb[1] << 8) + rgb[2];
    }
}


export class PlayerContainer extends BaseContainer<Player> implements MonitorContainer {
    constructor(props: ModelProperty) {
        super(props.model, Player, {});
    }
}