import { useCallback, useEffect, useMemo } from "react";
import useCache from "./cache";

type Handlers<T> = [
  (key: string) => T | undefined,
  (key: string, value: T) => void,
  string[]
];

/**
 * ```
 * const [get, put, expired] = useFlash(keySet);
 * ```
 */
const useFlash = <T>(
  primaries: string[],
  primaryTimeout: number = 30 * 1000
): Handlers<T> => {
  /**
   * データの内、primary の中にあるものは長時間保持する
   */
  const [getStorage, putStorage, removeStorage, keysStorage] = useCache<T>(
    primaryTimeout
  );

  /**
   * primaryにないデータも短時間キャッシュする。
   * (primaryが頻繁に変わっても、外部からのリソース読み込み回数を減らすため)
   */
  const [getTmp, putTmp, , keysTmp] = useCache<T>();

  const get = useCallback(
    (key: string) => {
      const d = getStorage(key);
      return d !== undefined ? d : getTmp(key);
    },
    [getStorage, getTmp]
  );

  const put = useCallback(
    (key: string, data: T) => {
      putTmp(key, data);
      if (primaries.includes(key)) {
        putStorage(key, data);
      }
    },
    [putTmp, putStorage, primaries]
  );

  const expiredKeys = useMemo(
    () =>
      primaries.filter(
        key => !keysStorage.includes(key) && !keysTmp.includes(key)
      ),
    [primaries, keysStorage, keysTmp]
  );

  useEffect(() => {
    console.info(`clean up flash repository`);
    keysStorage.forEach(key => {
      if (primaries.includes(key)) {
        removeStorage(key);
      }
    });
  }, [removeStorage, keysStorage, primaries]);

  return [get, put, expiredKeys];
};

export default useFlash;
