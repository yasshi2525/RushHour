import { put, call } from "redux-saga/effects";
import { ActionType } from "../actions";

const url = "http://5be50b8595e4340013f89011.mockapi.io/gamemap/1/";

const requestGameMap = () => 
    fetch(url).then( response => {
        if (!response.ok) {
            throw Error(response.statusText);
        }
        return response;
    })
    .then( response => response.json())
    .catch( error => error);

export default function* fetchGameMap() {
    try {
        const result = yield call(requestGameMap);

        yield put({type: ActionType.FETCH_MAP_SUCCEEDED, gamemap : result});
    } catch (e) {
        yield put({type: ActionType.FETCH_MAP_FAILED, message: e.message});
    }
}