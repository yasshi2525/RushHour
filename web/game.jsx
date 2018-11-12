import React from "react";
import ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { store } from "./store";
import GameBoard from "./components/GameBoard";

ReactDOM.render(
    <Provider store={store}>
        <GameBoard readOnly={true} />
    </Provider>
, document.getElementById("rushhourContainer"));
