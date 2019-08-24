import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import BaseModel from "./base";
import BaseContainer from "./container";
import { ModelProperty } from "../interfaces/pixi";

const defaultValues: { color: number } = { color: 0 };

export class Player extends BaseModel implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}


export class PlayerContainer extends BaseContainer<Player> implements MonitorContainer {
    constructor(props: ModelProperty) {
        super(props.model, Player, {});
    }
}