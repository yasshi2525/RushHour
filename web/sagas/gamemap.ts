import * as Action from "../actions";
import { requestURL, isOK } from ".";

const fetch_url = "api/v1/gamemap";
const diff_url = "api/v1/gamemap/diff";

function buildQuery(opts: Action.GameMapRequest): string {
    let params = new URLSearchParams();
    params.set("cx", opts.model.coord.cx.toString());
    params.set("cy", opts.model.coord.cy.toString());
    params.set("scale", (opts.model.coord.scale + 1).toString());
    params.set("delegate", opts.model.delegate.toString());
    return params.toString();
}

const request = (url: string, params: Action.GameMapRequest): Promise<any> => 
    fetch(url + "?" + buildQuery(params))
    .then(isOK)
    .then(response => response.json())
    .then(response => {
        if (!response.status) {
            throw Error(response.results);
        }
        let error = params.model.gamemap.mergeAll(response.results);
        params.model.timestamp = response.timestamp;
        if (params.model.gamemap.isChanged()) {
            params.model.gamemap.updateDisplayInfo();
        }
        return { ...response, ...error };
    })
    .catch(error => error);

export function* fetchMap(action: ReturnType<typeof Action.fetchMap.request>) {
    return yield requestURL({ request, url: fetch_url, args: action, callbacks: Action.fetchMap });
}

export function* diffMap(action: ReturnType<typeof Action.diffMap.request>) {
    return yield requestURL({ request, url: diff_url, args: action, callbacks: Action.diffMap });
}
