import SpriteModel from "./sprite";
import { Monitorable } from "../interfaces/monitor";

const defaultValues: {[index:string]: {}} = {
    alpha: 0.5
};

export class Residence extends SpriteModel implements Monitorable {
    constructor(options: {container: PIXI.Container, loader: PIXI.loaders.Loader}) {
        super({name: "residence", ...options});
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}

export class Company extends SpriteModel implements Monitorable  {
    constructor(options: {container: PIXI.Container, loader: PIXI.loaders.Loader}) {
        super({name: "company", ...options});
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}