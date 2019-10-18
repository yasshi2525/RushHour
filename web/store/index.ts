import { createStore, applyMiddleware } from "redux";
import createSagaMiddleware from "redux-saga";
import { UserInfo } from "../common/interfaces";
import { defaultState } from "../state";
import rootReducer from "../reducers"; 
import { rushHourSaga } from "../sagas";

export default function(my: UserInfo | undefined) {
    let sagaMiddleware = createSagaMiddleware();
    let store = createStore(rootReducer, defaultState(my), applyMiddleware(sagaMiddleware));
    sagaMiddleware.run(rushHourSaga);
    return store;
}

