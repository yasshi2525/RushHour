import { useEffect, useCallback, useState, useRef } from "react";
import * as PIXI from "pixi.js";

type Counter = (offset: number, delta: number) => void;

export const useCounter = (
  pixi: PIXI.Application,
  fn: Counter,
  interval: number,
  repeat: boolean = false
) => {
  const counter = useRef(0);
  const [cache, setCache] = useState(counter.current);
  const get = useCallback(() => counter.current, [counter]);
  const put = useCallback(
    (v: number) => {
      counter.current = v;
      setCache(v);
    },
    [counter]
  );

  useEffect(() => {
    console.info(`fn is changed, set counter 0`);
    put(0);
  }, [fn]);

  useEffect(() => {
    const wrapper = (delta: number) => {
      const now = get() + 1;
      if (now <= interval) {
        fn(now / interval, delta);
        if (now === interval && repeat) {
          console.info("counter is rounded");
          put(0);
        } else {
          put(now);
        }
      } else {
        console.info("counter is expired");
        pixi.ticker.remove(wrapper);
      }
    };
    pixi.ticker.add(wrapper);
    return () => {
      console.info("cleanup old counter");
      pixi.ticker.remove(wrapper);
    };
  }, [pixi, fn, interval, repeat]);

  return cache;
};
