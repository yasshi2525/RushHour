import { createStore, applyMiddleware } from "redux";
import createSagaMiddleware from "redux-saga";
import { defaultState } from "../state";
import rootReducer from "../reducers"; 
import rushHourSaga from "../sagas";

const sagaMiddleware = createSagaMiddleware();
export const store = createStore(rootReducer, defaultState, applyMiddleware(sagaMiddleware));
sagaMiddleware.run(rushHourSaga);

