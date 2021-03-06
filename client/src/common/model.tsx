import React, { createContext, useEffect, useReducer, useContext } from "react";
import { createAction } from "typesafe-actions";
import * as PIXI from "pixi.js";
import { config } from "interfaces/gamemap";
import { UserInfo } from "interfaces/user";
import { UnhandledError } from "interfaces/error";
import { ComponentProperty } from "interfaces/component";
import AuthContext from "./auth";
import { LoadingCircle } from "./loading";
import GameModel from "models";

const sheets = [
  "cursor",
  "anchor",
  "residence",
  "company",
  "rail_node",
  "rail_edge",
  "destroy"
];

const app = new PIXI.Application({
  width: window.innerWidth,
  height: window.innerHeight,
  backgroundColor: config.background,
  autoStart: true,
  antialias: true,
  resolution: window.devicePixelRatio,
  autoDensity: true
});
app.stage.sortableChildren = true;

const modelOptions = {
  app,
  cx: config.gamePos.default.x,
  cy: config.gamePos.default.y,
  scale: config.scale.default,
  zoom: 0
};

interface LoadingStatus {
  completed: boolean;
  error: UnhandledError | null;
  model: GameModel;
}

const ModelContext = createContext<GameModel>(
  new GameModel({ ...modelOptions, my: 0 })
);

const RESET = "RESET";
const ERR = "ERR";
const OK = "OK";

const reset = createAction(RESET)();
const err = createAction(ERR, (e: UnhandledError) => e)();
const ok = createAction(OK)();
type Actions =
  | ReturnType<typeof reset>
  | ReturnType<typeof err>
  | ReturnType<typeof ok>;

const reducer = (state: LoadingStatus, action: Actions) => {
  switch (action.type) {
    case RESET:
      return { ...state, completed: false };
    case ERR:
      return { ...state, error: action.payload };
    case OK:
      return { ...state, completed: true };
  }
};

const loadImages = (model: GameModel) =>
  new Promise<GameModel>((resolve, reject) => {
    sheets.forEach(key => {
      model.app.loader.add(
        key,
        `assets/bundle/spritesheet/${key}@${Math.floor(
          model.renderer.resolution
        )}x.json`
      );
    });
    model.app.loader.load(() => {
      model.init();
      return resolve(model);
    });
    model.app.loader.onError = () => reject(model);
  });

const useModel = (my?: UserInfo | null) => {
  const [state, dispatch] = useReducer(reducer, {
    completed: false,
    error: null,
    model: new GameModel({
      ...modelOptions,
      my: my ? my.id : 0
    })
  });
  useEffect(() => {
    console.info("effect useModel");
    (async () => {
      await loadImages(state.model).catch(e =>
        dispatch(err(new UnhandledError(e)))
      );
      dispatch(ok());
    })();
    return () => {
      console.info("cleanup useModel");
      state.model.app.loader.reset();
      dispatch(reset());
    };
  }, [my]);
  return [state.completed, state.error, state.model] as [
    boolean,
    UnhandledError | null,
    GameModel
  ];
};

export const ModelProvider = (props: ComponentProperty) => {
  const [[, my]] = useContext(AuthContext);
  const [completed, err, model] = useModel(my);

  if (!completed) {
    return <LoadingCircle />;
  } else if (err) {
    <>
      <div>画像データの読み込みに失敗しました。</div>
      <div>画面を更新してください</div>
      {err?.messages.map(msg => (
        <div>{msg}</div>
      ))}
    </>;
  }
  return (
    <ModelContext.Provider value={model}>
      {props.children}
    </ModelContext.Provider>
  );
};

export default ModelContext;
