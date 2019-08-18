import { put, call } from "redux-saga/effects";
import * as Action from "../actions";

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
    fetch(url + "?" + buildQuery(params)).then(response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }
        return response;
    }).then(response => response.json())
    .then(response => {
        if (!response.status) {
            throw Error(response.results);
        }
        params.model.mergeAll(response.results);
        params.model.timestamp = response.timestamp;
        if (params.model.isChanged()) {
            params.model.render();
        }
        return response;
    })
    .catch(error => error);

export function* fetchMap(action: ReturnType<typeof Action.fetchMap.request>) {
    try {
        const response = yield call(request, fetch_url, action.payload);
        return yield put(Action.fetchMap.success(response));
    } catch (e) {
        return yield put(Action.fetchMap.failure(e));
    }
}

export function* diffMap(action: ReturnType<typeof Action.diffMap.request>) {
    try {
        const response: Action.GameMapResponse = yield call(request, diff_url, action.payload);
        return yield put(Action.diffMap.success(response));
    } catch (e) {
        return yield put(Action.diffMap.failure(e));
    }
}
