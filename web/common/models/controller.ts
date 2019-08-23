import GroupModel from "./group";
import { Anchor, Cursor } from "./cursor";
import { XBorderContainer, YBorderContainer, WorldBorder } from "./border";

export default class extends GroupModel {
    init() {
        let props = { model: this.model, app: this.model.app, offset: this.model.offset, delegate: this.model.delegate };
        let anchor = new Anchor(props);
        this.objects.anchor = anchor;
        this.objects.cursor = new Cursor({ ...props, anchor: anchor });
        this.containers.xborder = new XBorderContainer(props);
        this.containers.yborder = new YBorderContainer(props);
        this.objects.world = new WorldBorder(props);
        
        super.init();
    }

    tick() {
        super.tick();
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
