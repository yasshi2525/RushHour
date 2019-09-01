import { MenuStatus } from "../state";
import * as Action from "../actions";
import { requestURL } from ".";

const departUrl = "api/v1/dept";
const extendUrl = "api/v1/extend"

function buildDepartQuery(params: Action.PointRequest): string {
    let res = new URLSearchParams();
    res.set("oid", params.oid.toString());
    res.set("x", params.x.toString());
    res.set("y", params.y.toString());
    res.set("scale", params.scale.toString());
    return res.toString();
}

function buildExtendQuery(params: Action.ExtendRequest): string {
    let res = new URLSearchParams();
    res.set("oid", params.oid.toString());
    res.set("x", params.x.toString());
    res.set("y", params.y.toString());
    res.set("scale", params.scale.toString());
    res.set("rnid", params.rnid.toString());
    return res.toString();
}


const requestDepart = (url: string, params: Action.PointRequest) => 
    fetch(url, {
        method: "POST",
        body: buildDepartQuery(params),
        headers: new Headers({ "Content-type" : "application/x-www-form-urlencoded" })
    }).then(response => response.json())
    .then(response => {
        let model = params.model;
        let anchorObj = model.gamemap.mergeChild("rail_nodes", response.results.rn);
        if (anchorObj !== undefined) {
            anchorObj.resolve({});
            model.setMenuState(MenuStatus.EXTEND_RAIL);
            model.controllers.getAnchor().merge("anchor", {
                type: "rail_nodes", 
                pos: anchorObj.get("pos"), 
                cid: anchorObj.get("cid")
            });
        }
        return response;
    })
    .catch(error => error);

const requestExtend = (url: string, params: Action.ExtendRequest) =>
    fetch(url, {
        method: "POST",
        body: buildExtendQuery(params),
        headers: new Headers({ "Content-type" : "application/x-www-form-urlencoded" })
    }).then(response => response.json())
    .then(response => {
        let model = params.model;
        let anchorObj = model.gamemap.mergeChild("rail_nodes", response.results.rn);
        let e1 = model.gamemap.mergeChild("rail_edges", response.results.e1);
        let e2 = model.gamemap.mergeChild("rail_edges", response.results.e2);
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
        return response;
    })
    .catch(error => error);

export function* depart(action: ReturnType<typeof Action.depart.request>) {
    return yield requestURL({ request: requestDepart, url: departUrl, args: action, callbacks: Action.depart });
}

export function* extend(action: ReturnType<typeof Action.extend.request>) {
    return yield requestURL({ request: requestExtend, url: extendUrl, args: action, callbacks: Action.extend });
}