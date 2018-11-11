import { createStore, applyMiddleware } from "redux";
import createSagaMiddleware from 'redux-saga';
import rootReducer from "./reducers";
import rushHourSaga from './sages';
import React from 'react';
import ReactDOM from 'react-dom'
import { Provider } from "react-redux";
import GameBoard from './components/GameBoard';

const sagaMiddleware = createSagaMiddleware();
const store = createStore(rootReducer, applyMiddleware(sagaMiddleware));
sagaMiddleware.run(rushHourSaga);

export default class RushHourApp extends React.Component {
    constructor(props) {
        super(props);
        this.readOnly = props.readOnly;
    }

    render() {
        return <GameBoard store={store} />;
    }
}

ReactDOM.render(
    <Provider store={store}>
        <GameBoard readOnly="true" />
    </Provider>
, document.getElementById("rushhourContainer"));
