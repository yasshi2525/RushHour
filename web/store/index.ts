import { createStore, applyMiddleware } from "redux";
import createSagaMiddleware from "redux-saga";
import { RushHourStatus } from "../state";
import rootReducer from "../reducers"; 
import rushHourSaga from "../sages";

const initState: RushHourStatus = {
    readOnly: true,
    map: {
        "residences": [],
        "companies": []
    }
};
const sagaMiddleware = createSagaMiddleware();
export const store = createStore(rootReducer, initState, applyMiddleware(sagaMiddleware));
sagaMiddleware.run(rushHourSaga);

