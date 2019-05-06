import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { PIXIModel, PIXIContainer } from "./pixi";

const defaultValues = {x: 0, y: 0};

export abstract class PointModel extends PIXIModel implements Monitorable {

    setInitialValues(initialValues: {[index: string]: {}}) {
        super.setInitialValues(initialValues);
        this.current = this.toView(this.props.x, this.props.y);
        this.destination = this.toView(this.props.x, this.props.y);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        ["x", "y"].forEach(v => this.addUpdateCallback(v, () => this.updateDestination()));
    }
}

export abstract class PointContainer<T extends PointModel> extends PIXIContainer<T> implements MonitorContrainer {        
}