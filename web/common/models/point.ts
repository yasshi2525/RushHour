import { Monitorable } from "../interfaces/monitor";
import BaseModel from "./base";

const defaultValues: {x: number, y:number} = {x: 0, y: 0};

export default abstract class extends BaseModel implements Monitorable {
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }
}
