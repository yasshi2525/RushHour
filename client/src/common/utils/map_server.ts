import { useCallback, useState, useMemo, useEffect } from "react";
import { MultiError, Errors } from "interfaces/error";
import {
  FetchMap,
  fetchMap,
  FetchMapResponse,
  FetchMapResponseKeys
} from "interfaces/endpoint";
import { GraceHandler } from "./flash";
import { useMultiHttpGetTask, useHttpGetTask } from "./http_get";
import { OkResponse, ErrorResponse } from "./http_common";
import { newInstance, CoreMap } from "./map_core";

const X = 0;
const Y = 1;
const S = 2;
const D = 3;

export const encode = (hash: string): FetchMap => {
  const v = hash.split("_").map(v => parseInt(v));
  return {
    x: v[X],
    y: v[Y],
    scale: v[S],
    delegate: v[D]
  };
};

const decode = (v: FetchMap) => `${v.x}_${v.y}_${v.scale}_${v.delegate}`;

const toCoreMap = (payload: FetchMapResponse, scale: number) => {
  const coreMap = newInstance();
  FetchMapResponseKeys.forEach(key => {
    Object.values(payload[key]).forEach(obj => {
      coreMap[key][obj.id] = {
        ...obj,
        scale,
        memberOf: []
      };
    });
  });
  return coreMap;
};

const makeBulkHandler = (
  bulkInsert: (payloadList: OkResponse<FetchMap, FetchMapResponse>[]) => void,
  setError: (v: MultiError) => void
) => ({
  onOK: bulkInsert,
  onError: (payloadList: ErrorResponse<FetchMap>[]) =>
    setError(new MultiError(payloadList.map(v => v.error)))
});

const makeHandler = (
  insert: (payload: OkResponse<FetchMap, FetchMapResponse>) => void,
  setError: (v: MultiError) => void
) => ({
  onOK: (payload: FetchMapResponse, args: FetchMap) =>
    insert({ payload, args }),
  onError: (e: Errors) => setError(new MultiError([e]))
});

type Handlers = [MultiError | undefined, () => void, () => void];

/**
 * ```
 * const [error, onBulkGrace] = useServerMap(key, keyAll, put, handler);
 * ```
 */
const useServerMap = (
  key: string,
  keyAll: string[],
  put: (key: string, value: CoreMap) => void,
  graceHandler: GraceHandler
): Handlers => {
  const register = useCallback(
    (response: OkResponse<FetchMap, FetchMapResponse>) => {
      put(
        decode(response.args),
        toCoreMap(response.payload, response.args.scale)
      );
    },
    [put]
  );

  /**
   * サーバからの複数レスポンスをストレージに登録、結果を一つにまとめる
   */
  const bulkInsert = useCallback(
    (responseList: OkResponse<FetchMap, FetchMapResponse>[]) =>
      responseList.forEach(register),
    [register]
  );

  const [error, setError] = useState<MultiError>();

  const bulkHandler = useMemo(() => makeBulkHandler(bulkInsert, setError), [
    bulkInsert
  ]);

  const [bulkFetch, bulkFetchCancel] = useMultiHttpGetTask(
    fetchMap,
    bulkHandler
  );

  /**
   * リクエストパラメタとしてストレージにデータがないチャンクを指定する
   * キャッシュにない coord 周辺チャンクのデータをサーバから取得する
   */
  const onBulkGrace = useCallback(
    (required: string[]) => {
      console.info(`request bulk fetch (cancel prev request): ${required}`);
      bulkFetchCancel();
      bulkFetch(required.map(encode));
    },
    [bulkFetch, bulkFetchCancel]
  );

  const handler = useMemo(() => makeHandler(register, setError), [register]);

  const [singleFetch, singleFetchCancel] = useHttpGetTask(fetchMap, handler);

  useEffect(() => {
    graceHandler.prepared = true;
    graceHandler.send = onBulkGrace;
  }, [graceHandler, onBulkGrace]);

  const reload = useCallback(() => {
    singleFetchCancel();
    singleFetch(encode(key));
  }, [singleFetch, key]);

  const bulkReload = useCallback(() => {
    bulkFetchCancel();
    bulkFetch(keyAll.map(encode));
  }, [bulkFetch, keyAll]);

  return [error, reload, bulkReload];
};

export default useServerMap;
