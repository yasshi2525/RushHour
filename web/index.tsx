import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { CircularProgress } from "@material-ui/core";
import { ThemeProvider } from "@material-ui/styles";
import GameContainer from "./common";
import { UserInfo } from "./common/interfaces";
import store from "./store";
import RushHourTheme from "./components";
import AppBar from "./components/AppBar";
import ActionMenu from "./components/ActionMenu";
import Canvas from "./components/Canvas";
import { loadImages } from "./sagas/model";
import { fetchMap } from "./sagas/map";
import { fetchPlayers } from "./sagas/player";

let props = document.getElementById("properties")
let appbar = document.getElementById("appbar")
let actionmenu = document.getElementById("actionmenu")
let canvas = document.getElementById("canvas")
let loading = document.getElementById("loading")

if (props !== null && appbar !== null && actionmenu !== null && canvas !== null && loading !== null) {
    let my: UserInfo | undefined = props.dataset.loggedin ? {
            name: props.dataset.displayname as string,
            image: props.dataset.image as string,
            hue: parseInt(props.dataset.hue as string)
        } : undefined;

    let isAdmin = props.dataset.admin !== undefined;
    let maintenance = props.dataset.maintenance !== undefined;

    function wrap(props: any = {}){
        return function(Component: React.ComponentType){
            return (<Provider store={store({my, isAdmin, maintenance})}>
                <ThemeProvider theme={RushHourTheme}>
                    <Component {...props} />
                </ThemeProvider>
            </Provider>)
        }
    }
    
    ReactDOM.render(
        <ThemeProvider theme={RushHourTheme}>
            <CircularProgress />
        </ThemeProvider>, loading)

    ReactDOM.render(wrap()(AppBar), appbar);

    let myid = (props.dataset.oid !== undefined) ? parseInt(props.dataset.oid) : 0
    let game = new GameContainer(myid);
    
    loadImages(game)
    .then(() => {
        if (loading !== null) {
            loading.remove();
        }
        if (my !== undefined && !maintenance) {
            ReactDOM.render(wrap({ model: game.model })(ActionMenu), actionmenu);
        }
        ReactDOM.render(wrap({
            readOnly: my === undefined,
            model: game.model
        })(Canvas), canvas);
        return fetchPlayers(game.model.gamemap);
    }).then(() => fetchMap(game.model));
} 