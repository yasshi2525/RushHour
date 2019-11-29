import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { Monitorable, MonitorContainer } from "interfaces/monitor";
import {
  AnimatedSpriteProperty,
  AnimatedSpriteContainerProperty
} from "interfaces/pixi";

const defaultValues: { [index: string]: {} } = {
  alpha: 1
};

export class Residence extends AnimatedSpriteModel implements Monitorable {
  getResourceName() {
    return "residences";
  }
  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }
}

export class ResidenceContainer
  extends AnimatedSpriteContainer<Residence, AnimatedSpriteProperty>
  implements MonitorContainer {
  constructor(options: AnimatedSpriteContainerProperty) {
    super(options, Residence);
  }

  protected getChildOptions() {
    return this.getBasicChildOptions();
  }
}

export class Company extends AnimatedSpriteModel implements Monitorable {
  getResourceName() {
    return "companies";
  }
  setupDefaultValues() {
    super.setupDefaultValues();
    this.addDefaultValues(defaultValues);
  }
}

export class CompanyContainer
  extends AnimatedSpriteContainer<Company, AnimatedSpriteProperty>
  implements MonitorContainer {
  constructor(options: AnimatedSpriteContainerProperty) {
    super(options, Company);
  }

  protected getChildOptions() {
    return this.getBasicChildOptions();
  }
}
