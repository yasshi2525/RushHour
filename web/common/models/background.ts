import SpriteModel from "./sprite";
import { Monitorable } from "../interfaces/monitor";
import { LocalableProprty } from "../interfaces/pixi";

const defaultValues: {[index:string]: {}} = {
    alpha: 1
};

export class Residence extends SpriteModel implements Monitorable {
    constructor(options: LocalableProprty) {
        super({name: "residence", ...options});
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class Company extends SpriteModel implements Monitorable  {
    constructor(options: LocalableProprty) {
        super({name: "company", ...options});
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}