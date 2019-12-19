import { useRef, useCallback, useEffect } from "react";
import * as PIXI from "pixi.js";
import { Pluggnable } from "./map_core";

interface Container<T extends Pluggnable> {
  props: T;
  container: PIXI.Container;
}

type Handler<T extends Pluggnable> = [
  (props: T) => Container<T>,
  (id: number) => void
];

const usePixiGraphics = <T extends Pluggnable>(
  stage: PIXI.Container,
  x: number,
  y: number,
  s: number,
  xs: number,
  ys: number
): Handler<T> => {
  const storage = useRef(new Map<number, Container<T>>());

  const add = useCallback(
    (props: T) => {
      const container = new PIXI.Container();
      stage.addChild(container);
      const body = { props, container };
      storage.current.set(props.id, body);
      return body;
    },
    [storage, stage]
  );

  const remove = useCallback(
    (id: number) => {
      const container = storage.current.get(id)?.container;
      if (container) {
        stage.removeChild(container);
      }
    },
    [storage, stage]
  );

  // visiblity
  useEffect(() => {
    for (const obj of storage.current.values()) {
      obj.container.visible = obj.props.scale == Math.floor(s);
    }
  }, [s]);

  // srv coord -> clt coord
  useEffect(() => {
    const xLen = Math.pow(2, xs);
    const yLen = Math.pow(2, ys);
    for (const obj of storage.current.values()) {
      obj.container.x = (obj.props.pos.x - x) / xLen;
      obj.container.y = (obj.props.pos.y - y) / yLen;
    }
  }, [x, y, s]);

  return [add, remove];
};

export default usePixiGraphics;
