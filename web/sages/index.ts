import { put, call , takeLatest } from "redux-saga/effects";
import { ActionType } from "../actions";
import { fetchMap } from "./gamemap";

/**
 * 非同期処理呼び出す ActionType を指定する。
 * ここで定義した ActionTypeをキャッチした際、個々のtsで定義した非同期メソッドが呼び出される
 */
export default function* rushHourSaga() {
    yield takeLatest(ActionType.FETCH_MAP_REQUESTED, fetchMap);
}

/**
 * 非同期タスクの呼び出しを一般化したもの。各 ts ファイルで実装を定義する
 * @param func 
 * @param okEvent ActionType
 * @param failEvent ActionType
 * @param tagName 
 */
export function* aync(func: () => Promise<any>, okEvent: ActionType, failEvent: ActionType, tagName: string) {
    try {
        const payload = yield call(func);
        return yield put({type: okEvent, [tagName] : payload});
    } catch (e) {
        return yield put({type: failEvent, message: e.message});
    }
}
