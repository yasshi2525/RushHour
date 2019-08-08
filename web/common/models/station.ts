import { SpriteModel, SpriteContainer } from "./sprite";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { SpriteContainerProperty } from "../interfaces/pixi";


const defaultValues: {[index:string]: {}} = {
    alpha: 1
};

export class Station extends SpriteModel implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class StationContainer extends SpriteContainer<Station> implements MonitorContrainer {
    constructor(options: SpriteContainerProperty) {
        super(options, Station, {});
    }
}