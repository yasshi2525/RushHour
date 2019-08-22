import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { store } from "./store";
import GameContainer from "./common";
import GameBoard from "./components/GameBoard";

const game = new GameContainer();

ReactDOM.render(
    <Provider store={store}>
        <GameBoard readOnly={true} game={game} />
    </Provider>
, document.getElementById("rushhourContainer"));
