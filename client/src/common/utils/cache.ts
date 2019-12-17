import { useRef, useEffect, useCallback } from "react";

export enum CacheStatus {
  ACTIVE,
  GRACE
}

type StatefulObject<T> = {
  data: T;
  status: CacheStatus;
  timer: NodeJS.Timeout;
};

type Handlers<T> = [
  (key: string) => T | undefined,
  (key: string, data: T) => void,
  (key: string, notify?: boolean) => T | undefined,
  (key: string[]) => [string, T][],
  (key: string[]) => void,
  (key: string) => CacheStatus | undefined
];

const ACTIVE = 15 * 1000;
const GRACE = 10 * 1000;

/**
 * 一定期間のみオブジェクトを保管する。
 * `[active] = activeTime => [grace] = graceTime => (delete)`
 * 例えば、サーバからデータをとるとき、状態がgraceになったら取得する。
 * この場合、graceTimeを想定レスポンスタイムにする
 * ```
 * const [get, put, remove, getIn, removeOutOf, status] = useCache(fn, fn, activeTime, graceTime);
 * ```
 */
const useCache = <T>(
  onGrace: (key: string, data: T) => void = () => {},
  onExpired: (key: string, data: T) => void = () => {},
  activeTime: number = ACTIVE,
  graceTime: number = GRACE
): Handlers<T> => {
  const storage = useRef(new Map<string, StatefulObject<T>>());

  const expire = useCallback(
    (key: string, notify: boolean = true) => {
      console.info(`expires ${key}`);
      const obj = storage.current.get(key);
      if (obj !== undefined) {
        clearTimeout(obj.timer);
        storage.current.delete(key);
        if (notify) {
          onExpired(key, obj.data);
        }
        return obj.data;
      }
      return undefined;
    },
    [storage, onExpired]
  );

  const grace = useCallback(
    (key: string, obj: StatefulObject<T>) => {
      console.info(`grace ${key}`);
      obj.status = CacheStatus.GRACE;
      clearTimeout(obj.timer);
      obj.timer = setTimeout(() => expire(key), graceTime);
      onGrace(key, obj.data);
    },
    [expire, onGrace, graceTime]
  );

  const get = useCallback((key: string) => storage.current.get(key)?.data, [
    storage
  ]);

  const put = useCallback(
    (key: string, data: T) => {
      const prev = storage.current.get(key);
      if (prev !== undefined) {
        console.info(`override ${key}, delete old timer`);
        clearTimeout(prev.timer);
      }
      const obj: StatefulObject<T> = {
        data,
        status: CacheStatus.ACTIVE,
        timer: setTimeout(() => grace(key, obj), activeTime)
      };
      storage.current.set(key, obj);
    },
    [storage, grace, activeTime]
  );

  const getIn = useCallback(
    (boundary: string[]) => {
      const result: [string, T][] = [];
      for (const [key, data] of storage.current.entries()) {
        if (boundary.includes(key)) {
          result.push([key, data.data]);
        }
      }
      return result;
    },
    [storage]
  );

  const removeOutOf = useCallback(
    (boundary: string[]) => {
      for (const key of storage.current.keys()) {
        if (!boundary.includes(key)) {
          expire(key);
        }
      }
    },
    [storage, expire]
  );

  const status = useCallback(
    (key: string) => storage.current.get(key)?.status,
    [storage]
  );

  return [get, put, expire, getIn, removeOutOf, status];
};

export default useCache;
