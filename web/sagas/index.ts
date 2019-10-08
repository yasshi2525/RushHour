import { takeEvery, takeLatest, put, call } from "redux-saga/effects";
import * as Action from "../actions";
import { generatePIXI } from "./model";
import { generateMap } from "./map";
import { generateDepart, generateExtend, generateConnect } from "./rail";
import { generatePlayers } from "./player";
import { generateSetMenu } from "./menu";

export function* generateRequest(
    request: any, 
    args: { [index:string]: any, payload: any },
    callbacks: { success: any, failure: any }) {
    try {
        const response = yield call(request, args.payload);
        return yield put(callbacks.success(response));
    } catch (e) {
        return yield put(callbacks.failure(e));
    }
}

async function isResponseOK(response: Response) {
    if (!response.ok) {
        throw Error(response.statusText);
    }
    return response;
}

async function isStatusOK(json: Action.ActionPayload) {
    if (!json.status) {
        throw Error(json.results);
    }
    return json;
};

async function validateResponse(rawRes: Response) {
    let res = await isResponseOK(rawRes);
    let rawJson = await res.json();
    return await isStatusOK(rawJson);
}

export async function httpGET(url: string) {
    let rawRes = await fetch(url);
    return await validateResponse(rawRes);
}

export async function httpPOST(url: string, params: {[index: string]: any}) {
    let rawRes = await fetch(url, { 
        method: "POST", 
        body: JSON.stringify(params, (key, value) => {
            return (key == "model") ? undefined : value
        }), 
        headers: new Headers({ "Content-type" : "application/json" })
    });
    return await validateResponse(rawRes);
}

/**
 * 非同期処理呼び出す ActionType を指定する。
 * ここで定義した ActionTypeをキャッチした際、個々のtsで定義した非同期メソッドが呼び出される
 */
export function* rushHourSaga() {
    yield takeLatest(Action.initPIXI.request, generatePIXI);
    yield takeLatest(Action.fetchMap.request, generateMap);
    yield takeLatest(Action.players.request, generatePlayers);
    yield takeLatest(Action.depart.request, generateDepart);
    yield takeLatest(Action.extend.request, generateExtend);
    yield takeLatest(Action.connect.request, generateConnect);
    yield takeEvery(Action.setMenu.request, generateSetMenu);
};
