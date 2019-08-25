import { PIXIModel, PIXIContainer } from "./pixi";
import { Monitorable, MonitorContainer } from "../interfaces/monitor";
import { Chunk, getChunkByPos, Point } from "../interfaces/gamemap";
import Cursor from "./cursor";
import Anchor from "./anchor";

const defaultValues: {
    pos: Point, cursor: Cursor | undefined
} = {pos: {x: 0, y: 0}, cursor: undefined};

export abstract class PointModel extends PIXIModel implements Monitorable {
    refferedCursor: Cursor | undefined;
    refferedAnchor: Anchor | undefined;

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
        this.addUpdateCallback("visible", (v) => {
            if (!v) {
                this.unreferCursor();
            }
        })
    }

    setupAfterCallback() {
        super.setupAfterCallback();
        this.addAfterCallback(() => this.unreferCursor());
    }

    standOnChunk(chunk: Chunk) {
        if (!this.props.visible) {
            return false;
        }
        let my = getChunkByPos(this.props.pos, chunk.scale);
        return chunk.x === my.x && chunk.y === my.y;
    }

    position(): Point | undefined {
        return this.current;
    }

    protected smoothMove() { 
        super.smoothMove();
        if (this.refferedCursor !== undefined) {
            this.refferedCursor.updateDisplayInfo();
        }
        if (this.refferedAnchor !== undefined) {
            this.refferedAnchor.updateDisplayInfo();
        }
    }

    protected unreferCursor() {
        if (this.refferedCursor !== undefined) {
            this.refferedCursor.unlinkSelected();
            this.refferedCursor.selectObject(this);
        }
        if (this.refferedAnchor !== undefined) {
            this.refferedAnchor.updateAnchor(true);
        }
    }

    protected setDisplayPosition() {
        let object = this.getPIXIObject();
        if (this.current !== undefined) {
            object.visible = true;
            object.x = this.current.x;
            object.y = this.current.y;
        } else {
            object.visible = false;
        }
    }

    protected followPointModel(following: PointModel | undefined, offset: number) {
        let object = this.getPIXIObject();
        if (following !== undefined) {
            object.visible = true;
            object.x = following.getPIXIObject().x + offset;
            object.y = following.getPIXIObject().y + offset;
            return true;
        } else {
            object.visible = false;
            return false;
        }
    }
}

const containerOpts: { 
    cursorClient: Point | undefined,
    anchorObj: PointModel | undefined,
    cursorObj: PointModel | undefined,
} = { 
    cursorClient: {x: 0, y: 0},
    anchorObj: undefined,
    cursorObj: undefined
};
export abstract class PointContainer<T extends PointModel> extends PIXIContainer<T> implements MonitorContainer {        
    setupDefaultValues() {
        super.setupDefaultValues();
        this.addDefaultValues(containerOpts);
    }
}