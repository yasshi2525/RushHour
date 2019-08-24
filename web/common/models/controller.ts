import { ZIndex } from "../interfaces/pixi";
import GroupModel from "./group";
import Anchor from "./anchor";
import Cursor from "./cursor";
import WorldBorder from "./world";
import { XBorderContainer, YBorderContainer } from "./border";

export default class extends GroupModel {
    init() {
        let props = { model: this.model, app: this.model.app, offset: this.model.offset, delegate: this.model.delegate };
        let anchor = new Anchor({ ...props, zIndex: ZIndex.ANCHOR });
        this.objects.anchor = anchor;
        this.objects.cursor = new Cursor({ ...props, anchor, zIndex: ZIndex.CURSOR });
        this.containers.xborder = new XBorderContainer({ ...props, zIndex: ZIndex.NORMAL_BORDER });
        this.containers.yborder = new YBorderContainer({ ...props, zIndex: ZIndex.NORMAL_BORDER });
        this.objects.world = new WorldBorder({ ...props, zIndex: ZIndex.WORLD_BORDER });
        
        super.init();
    }

    getCursor() {
        return (this.objects.cursor as Cursor);
    }

    getAnchor() {
        return (this.objects.anchor as Anchor);
    }

    updateAnchor() {
        (this.objects.anchor as Anchor).updateAnchor();
    }
}
