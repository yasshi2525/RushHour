import React, { useState, useCallback, useContext, useEffect } from "react";
import ModelContext from "./model";
import { useReload } from "./map";
import GameModel from "models";

type Point = [number, number];
const ZERO: Point = [0, 0];
const X = 0;
const Y = 1;

const mouse = (e: React.MouseEvent, model: GameModel): Point => [
  e.clientX * model.renderer.resolution,
  e.clientY * model.renderer.resolution
];

const touch = (e: React.TouchEvent, model: GameModel): Point => [
  e.targetTouches[0].clientX * model.renderer.resolution,
  e.targetTouches[0].clientY * model.renderer.resolution
];

const update = (model: GameModel, start: Point, now: Point) => {
  const size = Math.max(model.renderer.width, model.renderer.height);
  const zoom = Math.pow(2, model.coord.scale);
  const dx = ((now[X] - start[X]) / size) * zoom;
  const dy = ((now[Y] - start[Y]) / size) * zoom;

  model.setCoord(
    model.coord.cx - dx,
    model.coord.cy - dy,
    model.coord.scale,
    true
  );
};

export const useDrag = () => {
  const model = useContext(ModelContext);
  const reload = useReload();
  const [pressed, setPressed] = useState(false);
  const [start, setStart] = useState(ZERO);
  const [now, setNow] = useState(ZERO);

  const onMouseDown = useCallback(
    (e: React.MouseEvent) => {
      if (pressed) {
        console.warn("mousedown in pressed state useDrag");
      } else {
        console.info("start useDrag");
        const pos = mouse(e, model);
        setStart(pos);
        setNow(pos);
        setPressed(true);
      }
    },
    [pressed]
  );

  const onMouseMove = useCallback(
    (e: React.MouseEvent) => {
      if (pressed) {
        setNow(mouse(e, model));
      }
    },
    [pressed]
  );

  const onMouseUp = useCallback(
    (_: React.MouseEvent) => {
      if (!pressed) {
        console.warn("mouseup in unpressed state useDrag");
      } else {
        console.info("up useDrag");
        if (start !== now) {
          console.info("up useDrag reload");
          reload();
        }
        setPressed(false);
        setStart(ZERO);
        setNow(ZERO);
      }
    },
    [pressed]
  );

  const onMouseOut = useCallback(
    (_: React.MouseEvent) => {
      if (pressed) {
        console.info("out useDrag");
        setPressed(false);
        setStart(ZERO);
        setNow(ZERO);
      }
    },
    [pressed]
  );

  useEffect(() => {
    if (!pressed) {
      console.info("skip effect useDrag");
    } else {
      update(model, start, now);
    }
  }, [now]);

  return {
    onMouseDown,
    onMouseMove,
    onMouseUp,
    onMouseOut
  };
};

export const useSwipe = () => {
  const model = useContext(ModelContext);
  const reload = useReload();
  const [pressed, setPressed] = useState(false);
  const [start, setStart] = useState(ZERO);
  const [now, setNow] = useState(ZERO);

  const onTouchStart = useCallback(
    (e: React.TouchEvent) => {
      if (pressed) {
        console.info("touchstart in pressed state useSwipe");
      } else if (e.targetTouches.length > 1) {
        console.info("skip by multi touch useSwipe");
      } else {
        console.info("start useSwipe");
        const pos = touch(e, model);
        setStart(pos);
        setNow(pos);
        setPressed(true);
      }
    },
    [pressed]
  );

  const onTouchMove = useCallback(
    (e: React.TouchEvent) => {
      if (pressed) {
        setNow(touch(e, model));
      }
    },
    [pressed]
  );

  const onTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      if (!pressed) {
        console.info("touchend in unpressed state useSwipe");
      } else if (e.targetTouches.length >= 1) {
        console.warn("skip by multi touch useSwipe");
      } else {
        if (start !== now) {
          console.info("up useSwipe reload");
          reload();
        }
        setPressed(false);
        setStart(ZERO);
        setNow(ZERO);
      }
    },
    [pressed]
  );

  useEffect(() => {
    if (!pressed) {
      console.info("skip effect useSwipe");
    } else {
      update(model, start, now);
    }
  }, [now]);

  return {
    onTouchStart,
    onTouchMove,
    onTouchEnd
  };
};
