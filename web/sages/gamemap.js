import { put, call } from "redux-saga/effects";
import * as actions from "../actions";

const url = "http://5be50b8595e4340013f89011.mockapi.io/gamemap/";

function gamemap(timestamp) {
    return fetch(url + timestamp + "/")
        .then( response => {
            if (!response.ok) {
                throw Error(response.statusText);
            }
            return response;
        })
        .then( response => response.json())
        .catch( error => {error});
}

export default function* fetchGameMap(action) {
    try {
        const payload = yield call(gamemap, action.payload.timestamp);
        yield put({type: actions.FETCH_MAP_SUCCEEDED, payload : payload});
    } catch (e) {
        yield put({type: actions.FETCH_MAP_FAILED, message: e.message});
    }
}