import React, { FC, createContext, useMemo, useEffect } from "react";
import * as PIXI from "pixi.js";

const BACKGROUND = 0x263238;

const app = new PIXI.Application({
  width: window.innerWidth,
  height: window.innerHeight,
  backgroundColor: BACKGROUND,
  autoStart: true,
  antialias: true,
  resolution: window.devicePixelRatio,
  autoDensity: true
});
app.stage.sortableChildren = true;

const PixiContext = createContext<PIXI.Application>(app);
PixiContext.displayName = "PixiContext";

export const PixiProvider: FC = props => {
  useEffect(() => {
    console.info("after PixiProvider");
  }, []);

  return useMemo(
    () => (
      <PixiContext.Provider value={app}>{props.children}</PixiContext.Provider>
    ),
    []
  );
};
export default PixiContext;
