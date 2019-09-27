import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { store } from "./store";
import GameContainer from "./common";
import GameBoard from "./components/GameBoard";

const root = document.getElementById("rushhourContainer")

if (root !== null ) {
    const myid = (root.dataset.oid !== undefined) ? parseInt(root.dataset.oid) : 0
    const game = new GameContainer(myid);
    ReactDOM.render(
        <Provider store={store}>
            <GameBoard 
                readOnly={!root.dataset.loggedin}
                displayName={root.dataset.displayname}
                image={root.dataset.image}
                game={game} />
        </Provider>
    , root);
} 


