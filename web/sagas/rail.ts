import { MenuStatus } from "../state";
import * as Action from "../actions";
import { generateRequest, http, Method } from ".";

const departURL = "api/v1/rail_nodes";
const extendURL = "api/v1/rail_nodes/extend";
const connectURL = "api/v1/rail_nodes/connect";

export async function postDepart(params: Action.PointRequest) {
    let json = await http(departURL, Method.POST, params);
    let model = params.model;
    let anchorObj = model.gamemap.mergeChild("rail_nodes", json.rn);
    if (anchorObj !== undefined) {
        anchorObj.resolve({});
        model.setMenuState(MenuStatus.EXTEND_RAIL);
        model.controllers.getAnchor().merge("anchor", {
            type: "rail_nodes", 
            pos: anchorObj.get("pos"), 
            cid: anchorObj.get("cid")
        });
    }
    return json;
}

export async function postExtend(params: Action.ExtendRequest) {
    let json = await http(extendURL, Method.POST, params);
    let model = params.model;
    let anchorObj = model.gamemap.mergeChild("rail_nodes", json.rn);
    let e1 = model.gamemap.mergeChild("rail_edges", json.e1);
    let e2 = model.gamemap.mergeChild("rail_edges", json.e2);
    if (anchorObj !== undefined && e1 !== undefined && e2 !== undefined) {
        anchorObj.resolve({});
        e1.resolve({});
        e2.resolve({});
        model.controllers.getAnchor().merge("anchor", {
            type: "rail_nodes", 
            pos: anchorObj.get("pos"), 
            cid: anchorObj.get("cid")
        });
    }
    return json;
}

export async function postConnect(params: Action.ConnectRequest) {
    let json = await http(connectURL, Method.POST, params);
    let model = params.model;
    let anchorObj = model.gamemap.get("rail_nodes", json.e1.to);
    let e1 = model.gamemap.mergeChild("rail_edges", json.e1);
    let e2 = model.gamemap.mergeChild("rail_edges", json.e2);
    if (anchorObj !== undefined && e1 !== undefined && e2 !== undefined) {
        e1.resolve({});
        e2.resolve({});
        model.controllers.getAnchor().merge("anchor", {
            type: "rail_nodes", 
            pos: anchorObj.get("pos"), 
            cid: anchorObj.get("cid")
        });
    }
    return json;
}

export function* generateDepart(action: ReturnType<typeof Action.depart.request>) {
    return yield generateRequest(() => postDepart(action.payload), action, Action.depart);
}

export function* generateExtend(action: ReturnType<typeof Action.extend.request>) {
    return yield generateRequest(() => postExtend(action.payload), action, Action.extend);
}

export function* generateConnect(action: ReturnType<typeof Action.connect.request>) {
    return yield generateRequest(() => postConnect(action.payload), action, Action.connect);
}