import { useCallback } from "react";
import * as PIXI from "pixi.js";
import { Pluggnable } from "./map_core";
import usePixiGraphics from "./pixi_graphics";

interface RailNode extends Pluggnable {}

const useRailNode = (
  stage: PIXI.Container,
  offset: number,
  x: number,
  y: number,
  s: number,
  xs: number,
  ys: number
) => {
  const [_add, remove] = usePixiGraphics<RailNode>(stage, x, y, s, xs, ys);

  const add = useCallback(
    (props: RailNode) => {
      const body = _add(props);
    },
    [_add]
  );
};
