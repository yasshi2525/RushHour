import { useEffect, useCallback } from "react";
import * as PIXI from "pixi.js";

type Ticker = (delta?: number) => void;

export const useTicker = (app: PIXI.Application, fn: Ticker) => {
  const cancel = useCallback(() => app.ticker.remove(fn), [app]);

  useEffect(() => {
    app.ticker.add(fn);
    return () => {
      app.ticker.remove(fn);
    };
  }, [app]);

  return cancel;
};

export const useOnceTicker = (app: PIXI.Application, fn: Ticker) => {
  useEffect(() => {
    app.ticker.addOnce(fn);
    return () => {
      app.ticker.remove(fn);
    };
  }, [app]);
};
