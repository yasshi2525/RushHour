import { Monitorable, MonitorContainer } from "common/interfaces/monitor";
import {
  SpriteProperty,
  SpriteContainerProperty
} from "common/interfaces/pixi";
import { SpriteModel, SpriteContainer } from "./sprite";

const defaultValues: { [index: string]: {} } = {
  alpha: 1
};

export class Station extends SpriteModel implements Monitorable {
  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }
}

export class StationContainer extends SpriteContainer<Station, SpriteProperty>
  implements MonitorContainer {
  constructor(options: SpriteContainerProperty) {
    super(options, Station);
  }

  protected getChildOptions() {
    return this.getBasicChildOptions();
  }
}
