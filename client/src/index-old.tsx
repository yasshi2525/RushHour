import "typeface-roboto";
import * as React from "react";
import * as ReactDOM from "react-dom";
import { Provider } from "react-redux";
import { CircularProgress } from "@material-ui/core";
import { ThemeProvider } from "@material-ui/styles";
import "./index.css";
import GameContainer from "common";
import { jwtToUserInfo } from "state";
import store from "store";
import RushHourTheme from "components/Application";
import AppBar from "components/AppBar";
import ActionMenu from "components/ActionMenu";
import Canvas from "components/Canvas";
import { loadImages } from "sagas/model";
import { fetchMap } from "sagas/map";
import { fetchPlayers } from "sagas/player";

interface ModelProperty {
  children: JSX.Element;
}

let props = document.getElementById("properties");
let appbar = document.getElementById("appbar");
let actionmenu = document.getElementById("actionmenu");
let canvas = document.getElementById("canvas");
let loading = document.getElementById("loading");

if (
  props !== null &&
  appbar !== null &&
  actionmenu !== null &&
  canvas !== null &&
  loading !== null
) {
  let token = localStorage.getItem("jwt");
  let my = jwtToUserInfo(token);

  const Wrapper = (props: ModelProperty) => (
    <Provider store={store({ my })}>
      <ThemeProvider theme={RushHourTheme}>{props.children}</ThemeProvider>
    </Provider>
  );

  ReactDOM.render(
    <ThemeProvider theme={RushHourTheme}>
      <CircularProgress />
    </ThemeProvider>,
    loading
  );

  ReactDOM.render(
    <Wrapper>
      <AppBar />
    </Wrapper>,
    appbar
  );

  let myid = my !== undefined ? my.id : 0;
  let game = new GameContainer(myid);

  loadImages(game)
    .then(() => {
      if (loading !== null) {
        loading.remove();
      }
      if (my !== undefined) {
        ReactDOM.render(
          <Wrapper>
            <ActionMenu model={game.model} />
          </Wrapper>,
          actionmenu
        );
      }
      ReactDOM.render(
        <Wrapper>
          <Canvas readOnly={my === undefined} model={game.model} />
        </Wrapper>,
        canvas
      );
      return fetchPlayers(game.model.gamemap);
    })
    .then(() => fetchMap(game.model));
}
