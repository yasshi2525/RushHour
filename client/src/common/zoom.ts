import React, { useState, useCallback, useContext, useEffect } from "react";
import ModelContext from "./model";
import GameModel from "models";

const SENSITIVITY = 0.1;
type Coord = [number, number, number];
const X = 0;
const Y = 1;
const S = 2;

/**
 * posのクライアント座標系を中心にscaleを変更したとき、
 * 新たな画面中央のサーバ座標系の座標を取得します。
 */
const _zoom = (model: GameModel, server: Coord, client: Coord): Coord => {
  // 画面中心からの相対座標
  const size = Math.max(model.renderer.width, model.renderer.height);
  const zoom = Math.pow(2, server[S]);
  const diffCenter = [
    client[X] - model.renderer.width / 2,
    client[Y] - model.renderer.height / 2
  ];
  if (diffCenter[X] === 0 && diffCenter[Y] === 0) {
    return server;
  }

  // 中心座標の変更
  const diffSrv = [
    (diffCenter[X] / size) * zoom,
    (diffCenter[Y] / size) * zoom
  ];
  const distSrv = Math.sqrt(diffSrv[X] * diffSrv[X] + diffSrv[Y] * diffSrv[Y]);
  const nextDist = distSrv * Math.pow(2, client[S] - server[S]);
  const theta = Math.atan2(diffCenter[Y], diffCenter[X]);

  return [
    server[X] - (nextDist - distSrv) * Math.cos(theta),
    server[Y] - (nextDist - distSrv) * Math.sin(theta),
    client[S]
  ];
};

export const useZoom = () => {
  const model = useContext(ModelContext);
  return useCallback((x: number, y: number, _scale?: number) => {
    const scale = _scale === undefined ? model.coord.scale : _scale;
    const server: Coord = [model.coord.cx, model.coord.cy, model.coord.scale];
    return _zoom(model, server, [x, y, scale]);
  }, []);
};

const client = (e: React.WheelEvent, model: GameModel): Coord => {
  const dS = e.deltaY > 0 ? SENSITIVITY : -SENSITIVITY;
  const nextScale = model.coord.scale + dS;
  const roundNext = Math.round(nextScale * 10) / 10;
  return [
    e.clientX * model.renderer.resolution,
    e.clientY * model.renderer.resolution,
    roundNext
  ];
};

export const useWheel = () => {
  const model = useContext(ModelContext);
  const zoom = useZoom();
  const onWheel = useCallback((e: React.WheelEvent) => {
    model.setCoord(...zoom(...client(e, model)));
  }, []);
  return onWheel;
};

/**
 * [`x`, `y`, `dist`]
 */
type Gravity = [number, number, number];
const ZERO: Gravity = [0, 0, 0];

const distance = (e: React.TouchEvent, model: GameModel) => {
  const dx = e.targetTouches[0].clientX - e.targetTouches[1].clientX;
  const dy = e.targetTouches[0].clientY - e.targetTouches[1].clientY;
  return Math.sqrt(dx * dx + dy * dy) * model.renderer.resolution;
};

const gravity = (e: React.TouchEvent, model: GameModel): Gravity => {
  const pos = [0, 0];
  for (let i = 0; i < e.targetTouches.length; i++) {
    pos[X] += e.targetTouches[i].clientX;
    pos[Y] += e.targetTouches[i].clientY;
  }
  return [
    (pos[X] / e.targetTouches.length) * model.renderer.resolution,
    (pos[Y] / e.targetTouches.length) * model.renderer.resolution,
    distance(e, model)
  ];
};

export const usePinch = () => {
  const model = useContext(ModelContext);
  const [pressed, setPressed] = useState(false);
  const [start, setStart] = useState(ZERO);
  const [now, setNow] = useState(ZERO);
  const zoom = useZoom();

  const onTouchStart = useCallback(
    (e: React.TouchEvent) => {
      if (pressed) {
        console.info("touchstart in pressed state usePinch");
      } else if (e.targetTouches.length === 1) {
        console.info("skip by single touch usePinch");
      } else {
        console.info("start usePinch");
        const grv = gravity(e, model);
        setStart(grv);
        setNow(grv);
        setPressed(true);
      }
    },
    [pressed]
  );

  const onTouchMove = useCallback(
    (e: React.TouchEvent) => {
      if (pressed) {
        setNow(gravity(e, model));
      }
    },
    [pressed]
  );

  const onTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      if (!pressed) {
        console.info("touchend in unpressed state usePinch");
      } else if (e.targetTouches.length >= 2) {
        console.info("touchend in touch >= 2 state usePinch");
      } else {
        setPressed(false);
        setStart(ZERO);
        setNow(ZERO);
      }
    },
    [pressed]
  );

  useEffect(() => {
    if (!pressed) {
      console.info("skip effect usePinch");
    } else {
      const size = Math.max(model.renderer.width, model.renderer.height);
      const client: Coord = [
        now[X],
        now[Y],
        model.coord.scale -
          ((now[S] - start[S]) / size) * model.renderer.resolution
      ];
      model.setCoord(...zoom(...client));
    }
  }, [now]);

  return {
    onTouchStart,
    onTouchMove,
    onTouchEnd
  };
};
