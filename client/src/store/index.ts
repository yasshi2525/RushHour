import { createStore, applyMiddleware } from "redux";
import createSagaMiddleware from "redux-saga";
import { defaultState, DefaultProp } from "state";
import rootReducer from "reducers";
import { rushHourSaga } from "sagas";

export default function(opts: DefaultProp) {
  let sagaMiddleware = createSagaMiddleware();
  let store = createStore(
    rootReducer,
    defaultState(opts),
    applyMiddleware(sagaMiddleware)
  );
  sagaMiddleware.run(rushHourSaga);
  return store;
}
