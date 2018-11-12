import { takeLatest } from "redux-saga/effects";
import { ActionType } from "../actions";
import fetchGameMap from "./gamemap";

export default function* rushHourSaga() {
    yield takeLatest(ActionType.FETCH_MAP_REQUESTED, fetchGameMap);
}
