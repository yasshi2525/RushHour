import * as PIXI from "pixi.js";
import cursor from "./cursor";
import anchor from "./anchor";
import rail_node from "./rail_node";
import rail_edge from "./rail_edge";

let args = new URL(document.URL);
let type = args.searchParams.get("type");
let resolution = parseInt(args.searchParams.get("resolution") as string);

let app = new PIXI.Application({
    width: 40, height: 40, preserveDrawingBuffer: true
});

function toDataURL(obj: PIXI.DisplayObject) {
    return app.renderer.plugins.extract.canvas(obj).toDataURL();
}

window.addEventListener("load", () => {
    for(var offset = 0; offset < 240; offset++) {
        let anchorElm = document.createElement("a");
        anchorElm.id = `offset${offset}`;
        switch(type) {
            case "cursor":
                anchorElm.href = toDataURL(cursor(resolution, offset));
                break;
            case "anchor":
                anchorElm.href = toDataURL(anchor(resolution, offset));
                break;
            case "rail_node":
                anchorElm.href = toDataURL(rail_node(resolution, offset));
                break;
            case "rail_edge":
                anchorElm.href = toDataURL(rail_edge(resolution, offset));
        }
        document.body.appendChild(anchorElm);
    }
});
