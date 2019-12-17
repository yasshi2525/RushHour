import React, {
  FC,
  createContext,
  useEffect,
  useMemo,
  useState,
  useContext,
  useCallback
} from "react";
import { useCounter } from "./utils/tick";
import WindowContext from "./windows";
import ConfigContext from "./config";
import PixiContext from "./pixi";

/**
 * 座標に変化があった際、位置合わせをするまでの遅延フレーム数
 */
const LATENCY = 30;

interface ICoordinates {
  /**
   * 左上x座標(サーバにおけるマップ座標系)
   */
  x: number;
  /**
   * 左上y座標(サーバにおけるマップ座標系)
   */
  y: number;
  /**
   * 拡大率(クライアントウィンドウの幅がサーバマップ座標系の2^sに対応する)
   */
  s: number;
}
type Dimension = [number, number];
type Coordinates = [number, number, number];

/**
 * 座標系の中心を求める
 */
const defaultCoordinates = (min: number, max: number): Coordinates => {
  const scale = Math.floor((max + min) / 2);
  const pos = 1 << (scale - 1);
  return [pos, pos, scale];
};

/**
 * 画面横幅、縦幅が `scale` 表記でどの値になるか求める。
 * 長辺の幅が `2 ^ scale`
 * ```
 * 求める値 (x, y) は 画面幅 (w, h) に対し、
 *        w  = 2 ^ x
 *        h  = 2 ^ y
 * max(w, h) = 2 ^ scale
 *
 * <=>     x = log 2 w
 *         y = log 2 h
 *     scale = log 2 max(w, h)
 *
 * h = a * w (0 < a) とすると
 * y = log 2 h
 *   = log 2 (a * w)
 *   = log 2 a + log 2 w
 *   = log 2 a + x
 *   = log 2 (h / w) + x
 *
 * w > h ... x =                  scale, y = log 2 (h / w) + scale
 * w = h ... x =                  scale, y =                 scale
 * w < h ... x = -log 2 (h / w) + scale, y =                 scale
 * ```
 */
const scalize = (scale: number, width: number, height: number): Dimension =>
  width > height
    ? [scale, scale + Math.log2(height / width)]
    : [scale + Math.log2(width / height), scale];

/**
 * ```
 * left   = 0
 * top    = 0
 * right  = 2 ^ maxScale
 * bottom = 2 ^ maxScale
 *
 * サーバの座標幅 = srv(画面幅) とすると、
 * srv(w) = 2 ^ xScale
 * srv(h) = 2 ^ yScale
 *
 * 左上点 (x, y) の取りうる値は、
 * left < x
 * top  < y
 *        x + srv(w) < right
 *        y + srv(h) < bottom
 *
 * 0 < x
 * 0 < y
 *     x + 2 ^ xScale < 2 ^ maxScale
 *     y + 2 ^ yScale < 2 ^ maxScale
 *
 * 0 < x < 2 ^ maxScale - 2 ^ xScale
 * 0 < y < 2 ^ maxScale - 2 ^ yScale
 * ```
 */
const validate = (
  x: number,
  y: number,
  s: number,
  xScale: number,
  yScale: number,
  minScale: number,
  maxScale: number
): Coordinates => {
  const mapSize = 1 << maxScale;
  console.info(`mapSize=${mapSize} xS=${xScale} yS=${yScale}`);
  {
    const xLen = Math.pow(2, xScale);
    const left = 0;
    const right = mapSize - xLen;
    x = x < left ? left : x > right ? right : x;
    console.info(`xLen=${xLen}`);
  }
  {
    const yLen = Math.pow(2, yScale);
    const top = 0;
    const buttom = mapSize - yLen;
    y = y < top ? top : y > buttom ? buttom : y;
    console.info(`yLen=${yLen}`);
  }
  {
    const min = minScale;
    const max = maxScale;
    s = s < min ? min : s > max ? max : s;
  }

  return [x, y, s];
};

const DELTA = 0.01;
const round = (current: number, want: number) =>
  Math.abs(want - current) > DELTA ? current : want;
const smooth = (dest: number) => (delta: number) => (prev: number) =>
  round(prev * delta + dest * (1 - delta), dest);

const session = localStorage.getItem("loc");
const locSession: ICoordinates | undefined =
  session !== null ? JSON.parse(session) : undefined;

const initCoordinates = (
  width: number,
  height: number,
  minScale: number,
  maxScale: number
): Coordinates => {
  if (locSession === undefined) {
    return defaultCoordinates(minScale, maxScale);
  } else {
    const [xScale, yScale] = scalize(locSession.s, width, height);
    return validate(
      locSession.x,
      locSession.y,
      locSession.s,
      xScale,
      yScale,
      minScale,
      maxScale
    );
  }
};
type CoordState = [
  number,
  number,
  number,
  number,
  number,
  (x: number, y: number, s: number) => void
];
const useCoord = (): CoordState => {
  const [{ min_scale, max_scale }] = useContext(ConfigContext);
  const [w, h] = useContext(WindowContext);
  const pixi = useContext(PixiContext);

  const [initX, initY, initS] = useMemo(
    () => initCoordinates(w, h, min_scale, max_scale),
    [w, h, min_scale, max_scale]
  );

  const [x, setX] = useState(initX);
  const [y, setY] = useState(initY);
  const [s, setS] = useState(initS);

  useEffect(() => {
    console.info(`coord=(${x},${y},${s})`);
  }, [x, y, s]);

  const [xScale, yScale] = useMemo(() => scalize(s, w, h), [s, w, h]);

  const [dstX, setDstX] = useState(initX);
  const [dstY, setDstY] = useState(initY);
  const [dstS, setDstS] = useState(initS);

  const smoothX = useCallback((delta: number) => setX(smooth(dstX)(delta)), [
    dstX
  ]);
  const smoothY = useCallback((delta: number) => setY(smooth(dstY)(delta)), [
    dstY
  ]);
  const smoothS = useCallback((delta: number) => setS(smooth(dstS)(delta)), [
    dstS
  ]);

  const [latency, setLatency] = useState(LATENCY);

  const updateCoord = useCallback(
    (newX: number, newY: number, newS: number, force: boolean = false) => {
      console.info(`udpateCoord=${force ? 0 : LATENCY}`);
      const [vx, vy, vs] = validate(
        newX,
        newY,
        newS,
        xScale,
        yScale,
        min_scale,
        max_scale
      );
      console.info(`validate ${vx}, ${vy}, ${vs}`);
      setDstX(vx);
      setDstY(vy);
      setDstS(vs);
      setLatency(force ? 1 : LATENCY);
    },
    [xScale, yScale, min_scale, max_scale]
  );

  const smoother = useCallback(
    (offset: number) => {
      console.info(`smoother ${offset.toFixed(2)}`);
      const delta = Math.min(offset, 1 - offset);
      smoothX(delta);
      smoothY(delta);
      smoothS(delta);
    },
    [smoothX, smoothY, smoothS]
  );

  useCounter(pixi, smoother, latency);

  return [x, y, s, xScale, yScale, updateCoord];
};

/**
 * ```
 * [x, y, s, xScale, yScale, update] = useContext(CoordContext)
 * ```
 */
const CoordContext = createContext<CoordState>([
  0,
  0,
  0,
  0,
  0,
  () => {
    console.warn("not initialized CoordContext");
  }
]);
CoordContext.displayName = "CoordContext";

export const CoordProvider: FC = props => {
  const context = useCoord();
  useEffect(() => {
    console.info("after CoordProvider");
  }, []);
  return useMemo(
    () => (
      <CoordContext.Provider value={context}>
        {props.children}
      </CoordContext.Provider>
    ),
    [context]
  );
};

export default CoordContext;
