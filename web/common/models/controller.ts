import * as PIXI from "pixi.js";
import { ZIndex } from "../interfaces/pixi";
import GroupModel from "./group";
import Anchor from "./anchor";
import Cursor from "./cursor";
import WorldBorder from "./world";
import { XBorderContainer, YBorderContainer } from "./border";

export default class extends GroupModel {
    init() {
        let textures = this.model.app.loader.resources;
        let props = { 
            model: this.model, 
            app: this.model.app, 
            offset: this.model.offset, 
            delegate: this.model.delegate
        };
        
        let anchor_ss = textures["anchor"].spritesheet;
        let cursor_ss = textures["cursor"].spritesheet;

        let anchor_anim: PIXI.Texture[];
        let cursor_anim: PIXI.Texture[];

        if (anchor_ss !== undefined && cursor_ss !== undefined) {
            anchor_anim = anchor_ss.animations["anchor"]
            cursor_anim = cursor_ss.animations["cursor"];
        } else {
            anchor_anim = [PIXI.Texture.EMPTY];
            cursor_anim = [PIXI.Texture.EMPTY];
        }

        let anchor = new Anchor({ ...props, zIndex: ZIndex.ANCHOR, animation: anchor_anim });
        this.objects.anchor = anchor;
        this.objects.cursor = new Cursor({ ...props, anchor, zIndex: ZIndex.CURSOR, animation: cursor_anim });
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
        (this.objects.anchor as Anchor).updateAnchor(false);
    }
}
