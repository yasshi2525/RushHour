import React, {
  FC,
  createContext,
  useContext,
  useEffect,
  useMemo,
  useRef
} from "react";
import { MultiError } from "interfaces/error";
import useCoreMapStorage from "./utils/map_storage";
import useCoreMap, { CoreMap, newInstance } from "./utils/map_core";
import useServerMap from "./utils/map_server";
import useCoreMapChunk from "./utils/map_chunk";
import ConfigContext from "./config";
import DelegateContext from "./delegate";
import CoordContext from "./coord";
import { GraceHandler } from "./utils/flash";

type Handler = [MultiError | undefined, CoreMap, () => void, () => void];

/**
 * ```
 * const [isInitied, error, serverGameMap, reload, bulkReload] = useMap();
 * ```
 */
const useMap = (): Handler => {
  const [{ min_scale, max_scale }] = useContext(ConfigContext);
  const [coordX, coordY, coordS] = useContext(CoordContext);
  const delegate = useContext(DelegateContext);

  const [core, add, sub] = useCoreMap();

  const [current, cube] = useCoreMapChunk(
    coordX,
    coordY,
    coordS,
    delegate,
    min_scale,
    max_scale
  );

  const graceHandler = useRef<GraceHandler>({
    prepared: false,
    send: () => console.warn("not initialized grace handler")
  });

  const [, put, , key, keyAll] = useCoreMapStorage(
    current,
    cube,
    min_scale,
    graceHandler.current,
    add,
    sub
  );

  const [error, reload, bulkReload] = useServerMap(
    key,
    keyAll,
    put,
    graceHandler.current
  );

  return [error, core, reload, bulkReload];
};

type Contents = [CoreMap, () => void, () => void];

const GameMapContext = createContext<Contents>([
  newInstance(),
  () => console.warn("not initialized GameMapContext"),
  () => console.warn("not initialized GameMapContext")
]);
GameMapContext.displayName = "GameMapContext";

export const GameMapProvider: FC = props => {
  const [error, data, reload, bulkReload] = useMap();
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
  }, [error, contents]);
};

/**
 * ```
 * const [data, reload, bulkReload] = useContext(GameMapContext);
 * reload();
 * ```
 */
export default GameMapContext;
