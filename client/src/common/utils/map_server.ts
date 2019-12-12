import { useCallback, useEffect, useState } from "react";
import { MultiError } from "interfaces/error";
import {
  FetchMap,
  fetchMap,
  FetchMapResponse,
  FetchMapResponseKeys
} from "interfaces/endpoint";
import { useMultiHttpGet } from "./http_get";
import { newInstance, Chunk, CoreMap } from "./map_core";
import { OkResponse } from "./http_common";

const X = 0;
const Y = 1;
const S = 2;
const D = 3;

export const encode = (chunk: Chunk): FetchMap => ({
  x: chunk[X],
  y: chunk[Y],
  scale: chunk[S],
  delegate: chunk[D]
});

const decode = (args: FetchMap): Chunk => [
  args.x,
  args.y,
  args.scale,
  args.delegate
];

const hash = (x: number, y: number, scale: number, dlg: number) =>
  `${x}_${y}_${scale}_${dlg}`;

const toHash = (args: FetchMap) =>
  hash(args.x, args.y, args.scale, args.delegate);

const toCoreMap = (args: FetchMap, payload: FetchMapResponse) => {
  const coreMap = newInstance();
  FetchMapResponseKeys.forEach(key => {
    Object.values(payload[key]).forEach(obj => {
      coreMap[key][obj.id] = {
        ...obj,
        scale: args.scale,
        memberOf: [toHash(args)]
      };
    });
  });
  return coreMap;
};

type Handlers = [
  boolean,
  MultiError | undefined,
  (payload: OkResponse<FetchMap, FetchMapResponse>) => void,
  (payloadList: OkResponse<FetchMap, FetchMapResponse>[]) => void
];

/**
 * ```
 * const [response, put, putAll] = useServerMap(put, expired);
 * ```
 */
const useServerMap = (
  put: (key: Chunk, value: CoreMap) => void,
  expired: Chunk[]
): Handlers => {
  /**
   * リクエストパラメタとしてストレージにデータがないチャンクを指定する
   * キャッシュにない coord 周辺チャンクのデータをサーバから取得する
   */
  const response = useMultiHttpGet({
    ...fetchMap,
    argsList: expired.map(encode)
  });

  const insert = useCallback(
    (payload: OkResponse<FetchMap, FetchMapResponse>) =>
      put(decode(payload.args), toCoreMap(payload.args, payload.payload)),
    [put]
  );

  /**
   * 受信したレスポンスデータをストレージに格納する
   */
  const bulkInsert = useCallback(
    (payloadList: OkResponse<FetchMap, FetchMapResponse>[]) =>
      payloadList.forEach(insert),
    [put, insert]
  );

  const [error, setError] = useState<MultiError>();
  const [isInited, setInited] = useState(false);

  /**
   * サーバからの複数レスポンスをストレージに登録、結果を一つにまとめる
   */
  useEffect(() => {
    if (response) {
      console.info(`insert server response to flash memory`);
      bulkInsert(response.payloadList);
      if (response.errorList.length) {
        setError(
          new MultiError(response.errorList.map(payload => payload.error))
        );
      }
      setInited(true);
    }
  }, [response, bulkInsert]);

  return [isInited, error, insert, bulkInsert];
};

export default useServerMap;
