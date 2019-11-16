import { AnimatedSpriteModel, AnimatedSpriteContainer } from "./sprite";
import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import { AnimatedSpriteContainerProperty } from "../interfaces/pixi";

const defaultValues: {[index:string]: {}} = {
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

export class ResidenceContainer extends AnimatedSpriteContainer<Residence> implements MonitorContainer {
    constructor(options: AnimatedSpriteContainerProperty) {
        super(options, Residence, {});
    }
}

export class Company extends AnimatedSpriteModel implements Monitorable  {
    getResourceName() {
        return "companies";
    }
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class CompanyContainer extends AnimatedSpriteContainer<Company> implements MonitorContainer {
    constructor(options: AnimatedSpriteContainerProperty) {
        super(options, Company, {});
    }
}