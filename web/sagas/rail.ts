import { MenuStatus } from "../state";
import * as Action from "../actions";
import { requestURL } from ".";

const url = "api/v1/dept";

function buildQuery(params: Action.PointRequest): string {
    let res = new URLSearchParams();
    res.set("oid", params.oid.toString());
    res.set("x", params.x.toString());
    res.set("y", params.y.toString());
    res.set("scale", params.scale.toString());
    return res.toString();
}

const request = (url: string, params: Action.PointRequest) => 
    fetch(url, {
        method: "POST",
        body: buildQuery(params),
        headers : new Headers({"Content-type" : "application/x-www-form-urlencoded" })
    }).then(response => response.json())
    .then((response) => {
        let anchorObj = params.model.gamemap.mergeChild("rail_nodes", response.results.rn);
        if (anchorObj !== undefined) {
            params.model.controllers.getAnchor().merge("anchor", {
                type: "rail_nodes", 
                pos: anchorObj.get("pos"), 
                cid: anchorObj.get("cid")
            });
            params.model.setMenuState(MenuStatus.EXTEND_RAIL);
        }
        return response;
    })
    .catch(error => error);


export function* depart(action: ReturnType<typeof Action.depart.request>) {
    return yield requestURL({ request, url, args: action, callbacks: Action.depart });
}