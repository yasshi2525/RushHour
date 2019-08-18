import { PIXIModel, PIXIContainer } from "./pixi";
import { Monitorable, MonitorContrainer } from "../interfaces/monitor";
import { Chunk, getChunk, Point } from "../interfaces/gamemap";
import { Cursor } from "./cursor";

const defaultValues: {
    pos: Point, cursor: Cursor | undefined
} = {pos: {x: 0, y: 0}, cursor: undefined};

export abstract class PointModel extends PIXIModel implements Monitorable {
    refferedCursor: Cursor | undefined;

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
        this.addUpdateCallback("outMap", (v) => {
            if (v) {
                this.unreferCursor();
            }
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.unreferCursor());
    }

    standOnChunk(chunk: Chunk) {
        let my = getChunk(this.props.pos, chunk.scale);
        return chunk.x === my.x && chunk.y === my.y;
    }

    protected unreferCursor() {
        if (this.refferedCursor !== undefined) {
            this.refferedCursor.unlinkSelected();
            this.refferedCursor.selectObject(this);
        }
    }
}

export abstract class PointContainer<T extends PointModel> extends PIXIContainer<T> implements MonitorContrainer {        
}