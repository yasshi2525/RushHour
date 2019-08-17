import { PIXIModel, PIXIContainer } from "./pixi";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { Chunk, getChunk, Point } from "../interfaces/gamemap";
import CursorModel from "./cursor";

const defaultValues: {
    pos: Point, cursor: CursorModel | undefined
} = {pos: {x: 0, y: 0}, cursor: undefined};

export abstract class PointModel extends PIXIModel implements Monitorable {
    refferedCursor: CursorModel | undefined;

    setInitialValues(initialValues: {[index: string]: {}}) {
        super.setInitialValues(initialValues);
        this.current = this.toView(this.props.pos);
        this.destination = this.toView(this.props.pos);
    }

    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(defaultValues);
    }

    setupUpdateCallback() {
        super.setupUpdateCallback();
        this.addUpdateCallback("pos", () => this.updateDestination());
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => {
            if (this.refferedCursor !== undefined) {
                this.refferedCursor.unlinkSelected();
                this.refferedCursor.selectObject();
            }
        })
    }

    standOnChunk(chunk: Chunk) {
        let my = getChunk(this.props.pos, chunk.scale);
        return chunk.x === my.x && chunk.y === my.y;
    }
}

export abstract class PointContainer<T extends PointModel> extends PIXIContainer<T> implements MonitorContrainer {        
}