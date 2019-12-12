import React, {
  FC,
  createContext,
  useEffect,
  useMemo,
  useState,
  useContext,
  useCallback
} from "react";
import { Dimension, Coordinates } from "interfaces";
import { useOnceTicker } from "./utils/tick";
import WindowContext from "./windows";
import ConfigContext from "./config";
import PixiContext from "./pixi";

const X = 0;
const Y = 1;
const S = 2;
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
   * 拡大率(クライアントウィンドウの幅が2^sに対応する)
   */
  s: number;
}

const session = localStorage.getItem("loc");
const locSession: ICoordinates | undefined =
  session !== null ? JSON.parse(session) : undefined;

type CoordState = [
  Coordinates,
  Dimension,
  (x: number, y: number, s: number) => void
];
const useCoord = (): CoordState => {
  const [config] = useContext(ConfigContext);
  const window = useContext(WindowContext);
  const pixi = useContext(PixiContext);

  const mapSize = useMemo(() => 1 << (config.max_scale - config.min_scale), [
    config
  ]);
  const [coord, setCoord] = useState<Coordinates>(
    locSession === undefined
      ? [mapSize >> 1, mapSize >> 1, (config.max_scale + config.min_scale) / 2]
      : [locSession.x, locSession.y, locSession.s]
  );

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
  const scale = useMemo<Dimension>(() => {
    const [w, h] = window;
    return w > h
      ? [coord[S], coord[S] + Math.log2(h / w)]
      : [coord[S] + Math.log2(w / h), coord[S]];
  }, [coord, window]);

  const [dest, setDest] = useState<Coordinates>([coord[0], coord[1], coord[2]]);
  const [latency, setLatency] = useState(LATENCY);

  /**
   * ```
   * left   = 0
   * top    = 0
   * right  = 2 ^ (max_scale - min_scale)
   * bottom = 2 ^ (max_scale - min_scale)
   *
   * サーバの座標幅 = srv(画面幅) とすると、
   * srv(w) = 2 ^ scale[X]
   * srv(h) = 2 ^ scale[Y]
   *
   * 左上点 (x, y) の取りうる値は、
   * left < x
   * top  < y
   *        x + srv(w) < right
   *        y + srv(h) < bottom
   *
   * 0 < x
   * 0 < y
   *     x + 2 ^ scale[X] < 2 ^ (max_scale - min_scale)
   *     y + 2 ^ scale[Y] < 2 ^ (max_scale - min_scale)
   *
   * 0 < x < 2 ^ (max_scale - min_scale) - 2 ^ scale[X]
   * 0 < y < 2 ^ (max_scale - min_scale) - 2 ^ scale[Y]
   * ```
   */
  const updateCoord = useCallback(
    (x: number, y: number, s: number, force: boolean = false) => {
      {
        const xLen = 1 << scale[X];
        const left = 0;
        const right = mapSize - xLen;
        x = x < left ? left : x > right ? right : x;
      }
      {
        const yLen = 1 << scale[Y];
        const top = 0;
        const buttom = mapSize - yLen;
        y = y < top ? top : y > buttom ? buttom : y;
      }
      {
        const min = config.min_scale;
        const max = config.max_scale;
        s = s < min ? min : s > max ? max : s;
      }

      setDest([x, y, s]);
      setLatency(force ? 0 : LATENCY);
    },
    [config, scale]
  );

  const smoother = useCallback(() => {
    if (latency > 0) {
      let ratio = latency / LATENCY;
      if (ratio < 0.5) {
        ratio = 1.0 - ratio;
      }
      setCoord(now => [
        now[X] * ratio + dest[X] * (1 - ratio),
        now[Y] * ratio + dest[Y] * (1 - ratio),
        now[S] * ratio + dest[S] * (1 - ratio)
      ]);
      setLatency(v => v - 1);
    }
  }, [latency, dest]);

  useOnceTicker(pixi, smoother);

  return [coord, scale, updateCoord];
};

/**
 * chunk = coord / delegate
 * ```
 * [coord, scale, chunk, update] = useContext(CoordContext)
 * ```
 */
const CoordContext = createContext<CoordState>([
  [0, 0, 0],
  [0, 0],
  () => {
    console.warn("not initialized CoordContext");
  }
]);
CoordContext.displayName = "CoordContext";

export const CoordProvider: FC = props => {
  const [coord, scale, update] = useCoord();
  useEffect(() => {
    console.info("after CoordProvider");
  }, []);
  return useMemo(
    () => (
      <CoordContext.Provider value={[coord, scale, update]}>
        {props.children}
      </CoordContext.Provider>
    ),
    [coord, scale, update]
  );
};

export default CoordContext;
