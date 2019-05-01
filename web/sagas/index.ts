import { takeLatest } from "redux-saga/effects";
import * as Action from "../actions";
import { fetchMap, diffMap } from "./gamemap";

/**
 * 非同期処理呼び出す ActionType を指定する。
 * ここで定義した ActionTypeをキャッチした際、個々のtsで定義した非同期メソッドが呼び出される
 */
export default function* rushHourSaga() {
    yield takeLatest(Action.fetchMap.request, fetchMap);
    yield takeLatest(Action.diffMap.request, diffMap);
};
