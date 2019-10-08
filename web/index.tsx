import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { ThemeProvider } from "@material-ui/styles";
import GameContainer from "./common";
import { store } from "./store";
import RushHourTheme from "./components";
import AppBar from "./components/AppBar";
import ActionMenu from "./components/ActionMenu";
import Canvas from "./components/Canvas";
import { loadImages } from "./sagas/model";
import { fetchMap } from "./sagas/map";
import { fetchPlayers } from "./sagas/player";
import { CircularProgress } from "@material-ui/core";

let props = document.getElementById("properties")
let appbar = document.getElementById("appbar")
let actionmenu = document.getElementById("actionmenu")
let canvas = document.getElementById("canvas")
let loading = document.getElementById("loading")

function wrap(props: any){
    return function(Component: React.ComponentType){
        return (<Provider store={store}>
            <ThemeProvider theme={RushHourTheme}>
                <Component {...props} />
            </ThemeProvider>
        </Provider>)
    }
}

if (props !== null && appbar !== null && actionmenu !== null && canvas !== null && loading !== null) {
    let opts = !props.dataset.loggedin ? { readOnly: true } : {
        readOnly: false,
        displayName: props.dataset.displayname,
        image: props.dataset.image,
    };
    
    ReactDOM.render(
        <ThemeProvider theme={RushHourTheme}>
            <CircularProgress />
        </ThemeProvider>, loading)

    ReactDOM.render(wrap(opts)(AppBar), appbar);

    let myid = (props.dataset.oid !== undefined) ? parseInt(props.dataset.oid) : 0
    let game = new GameContainer(myid);
    
    loadImages(game)
    .then(() => {
        if (loading !== null) {
            loading.remove();
        }
        if (!opts.readOnly) {
            ReactDOM.render(wrap({ model: game.model })(ActionMenu), actionmenu);
        }
        ReactDOM.render(wrap({
            readOnly: opts.readOnly,
            model: game.model
        })(Canvas), canvas);
        return fetchPlayers(game.model.gamemap);
    }).then(() => fetchMap(game.model));
} 