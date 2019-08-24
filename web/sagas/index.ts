import { takeLatest, put, call } from "redux-saga/effects";
import * as Action from "../actions";
import { initPIXI } from "./model";
import { fetchMap, diffMap } from "./gamemap";
import { depart } from "./rail";
import { players } from "./player";

export interface requestPayload {
    request: any,
    url: string,
    args: { [index:string]: any, payload: any },
    callbacks: { success: any, failure: any }
}

export function* requestURL(params: requestPayload) {
    try {
        const response = yield call(params.request, params.url, params.args.payload);
        return yield put(params.callbacks.success(response));
    } catch (e) {
        return yield put(params.callbacks.failure(e));
    }
}

export const isOK = (response: Response) => {
    if (!response.ok) {
        throw Error(response.statusText);
    }
    return response;
}

/**
 * 非同期処理呼び出す ActionType を指定する。
 * ここで定義した ActionTypeをキャッチした際、個々のtsで定義した非同期メソッドが呼び出される
 */
export function* rushHourSaga() {
    yield takeLatest(Action.initPIXI.request, initPIXI);
    yield takeLatest(Action.fetchMap.request, fetchMap);
    yield takeLatest(Action.diffMap.request, diffMap);
    yield takeLatest(Action.players.request, players);
    yield takeLatest(Action.depart.request, depart);
};
