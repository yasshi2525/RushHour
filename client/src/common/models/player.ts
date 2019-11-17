import { Monitorable, MonitorContainer } from "common/interfaces/monitor";
import { hueToRgb } from "common/interfaces/gamemap";
import BaseModel, { BaseProperty } from "./base";
import BaseContainer from "./container";

const defaultValues: { color: number; hue: number } = { color: 0, hue: 0 };

export class Player extends BaseModel implements Monitorable {
  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }

  setInitialValues(props: { [index: string]: any }) {
    super.setInitialValues(props);
    let rgb = hueToRgb(props.hue);
    this.props.color = (rgb[0] << 16) + (rgb[1] << 8) + rgb[2];
  }
}

export class PlayerContainer extends BaseContainer<Player, BaseProperty>
  implements MonitorContainer {
  constructor(props: BaseProperty) {
    super(props, Player);
  }
  protected getChildOptions(): BaseProperty {
    return this.getBasicChildOptions();
  }
}
