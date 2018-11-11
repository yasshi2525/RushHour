import { takeLatest } from 'redux-saga/effects'
import { FETCH_MAP_REQUESTED } from "../actions";
import fetchGameMap from "./gamemap";

export default function* rushHourSaga() {
    yield takeLatest(FETCH_MAP_REQUESTED, fetchGameMap);
}
