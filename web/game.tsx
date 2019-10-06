import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { ThemeProvider } from "@material-ui/styles";
import GameContainer from "./common";
import { store } from "./store";
import RushHourTheme from "./components";
import GameBar from "./components/GameBar";
import ToolBar from "./components/Toolbar";
import Canvas from "./components/Canvas";
import { loadImages } from "./sagas/model";
import { fetchMap } from "./sagas/map";
import { fetchPlayers } from "./sagas/player";

let props = document.getElementById("properties")
let gamebar = document.getElementById("gamebar")
let toolbar = document.getElementById("toolbar")
let canvas = document.getElementById("canvas")

function wrap(props: any){
    return function(Component: React.ComponentType){
        return (<Provider store={store}>
            <ThemeProvider theme={RushHourTheme}>
                <Component {...props} />
            </ThemeProvider>
        </Provider>)
    }
}

if (props !== null && gamebar !== null && toolbar !== null && canvas !== null) {
    let opts = props.dataset.readOnly ? { readOnly: true } : {
        readOnly: false,
        displayName: props.dataset.displayname,
        image: props.dataset.image,
    };
    ReactDOM.render(wrap(opts)(GameBar), gamebar);

    let myid = (props.dataset.oid !== undefined) ? parseInt(props.dataset.oid) : 0
    let game = new GameContainer(myid);
    
    loadImages(game)
    .then(() => {
        if (!opts.readOnly) {
            ReactDOM.render(wrap({ model: game.model })(ToolBar), toolbar);
        }
        ReactDOM.render(wrap({
            readOnly: opts.readOnly,
            model: game.model
        })(Canvas), canvas);
        return fetchPlayers(game.model.gamemap);
    }).then(() => fetchMap(game.model));
} 