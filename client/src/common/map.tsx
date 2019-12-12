import React, {
  FC,
  createContext,
  useContext,
  useEffect,
  useMemo
} from "react";
import { MultiError } from "interfaces/error";
import useCoreMapStorage from "./utils/map_storage";
import useCoreMap, { CoreMap, newInstance } from "./utils/map_core";
import useServerMap from "./utils/map_server";
import useServerMapReload from "./utils/map_reload";
import ConfigContext from "./config";
import DelegateContext from "./delegate";
import CoordContext from "./coord";

type Handler = [
  boolean,
  MultiError | undefined,
  CoreMap,
  () => void,
  () => void
];

/**
 * ```
 * const [isInitied, error, serverGameMap, reload, bulkReload] = useMap();
 * ```
 */
const useMap = (): Handler => {
  const [config] = useContext(ConfigContext);
  const [coord] = useContext(CoordContext);
  const delegate = useContext(DelegateContext);

  const [getCore, putCore, keys, current, expired] = useCoreMapStorage(
    coord,
    delegate,
    config.min_scale,
    config.max_scale
  );

  const coreMap = useCoreMap(getCore, keys);

  const [isInited, error, put, putAll] = useServerMap(putCore, expired);

  const [reload, bulkReload] = useServerMapReload(current, keys, put, putAll);

  return [isInited, error, coreMap, reload, bulkReload];
};

type Contents = [CoreMap, () => void, () => void];

const GameMapContext = createContext<Contents>([
  newInstance(),
  () => console.warn("not initialized GameMapContext"),
  () => console.warn("not initialized GameMapContext")
]);
GameMapContext.displayName = "GameMapContext";

export const GameMapProvider: FC = props => {
  const [isInited, error, data, reload, bulkReload] = useMap();
  const contents = useMemo<Contents>(() => [data, reload, bulkReload], [
    data,
    reload,
    bulkReload
  ]);

  useEffect(() => {
    console.info("after GameMapProvider");
  }, []);

  return useMemo(() => {
    if (error) {
      return (
        <div>
          <p>マップ情報の読み込みに失敗しました</p>
          <p>画面を更新してください</p>
        </div>
      );
    } else {
      return (
        <GameMapContext.Provider value={contents}>
          {props.children}
        </GameMapContext.Provider>
      );
    }
  }, [isInited, error, contents]);
};

/**
 * ```
 * const [data, reload, bulkReload] = useContext(GameMapContext);
 * reload();
 * ```
 */
export default GameMapContext;
