import { SpriteModel, SpriteContainer } from "./sprite";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { SpriteContainerProperty } from "../interfaces/pixi";

const defaultValues: {[index:string]: {}} = {
    alpha: 1
};

export class Residence extends SpriteModel implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class ResidenceContainer extends SpriteContainer<Residence> implements MonitorContrainer {
    constructor(options: SpriteContainerProperty) {
        super(options, Residence, {});
    }
}

export class Company extends SpriteModel implements Monitorable  {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class CompanyContainer extends SpriteContainer<Company> implements MonitorContrainer {
    constructor(options: SpriteContainerProperty) {
        super(options, Company, {});
    }
}