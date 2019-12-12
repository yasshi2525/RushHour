import { useRef, useEffect, useCallback, useState, useMemo } from "react";

type Storage<T> = { [index: string]: { payload: T; timer: NodeJS.Timeout } };
type Handlers<T> = [
  (key: string) => T | undefined,
  (key: string, data: T) => void,
  (key: string) => void,
  string[]
];

/**
 * ```
 * const [get, put, remove, keys] = useCache();
 * ```
 */
const useCache = <T>(timeout: number = 5000): Handlers<T> => {
  const [lastUpdated, setLastUpdated] = useState(0);
  const storage = useRef<Storage<T>>({});

  const put = useCallback(
    (key: string, payload: T) => {
      if (key in storage.current) {
        console.info(`override ${key}, delete old timer`);
        clearTimeout(storage.current[key].timer);
      }
      console.info(`add ${key}`);
      storage.current[key] = {
        payload,
        timer: setTimeout(() => {
          console.info(`delete ${key} by timeout`);
          delete storage.current[key];
          setLastUpdated(new Date().getTime());
        }, timeout)
      };
      setLastUpdated(new Date().getTime());
    },
    [storage]
  );

  const get = useCallback(
    (key: string) =>
      key in storage.current ? storage.current[key].payload : undefined,
    [storage]
  );

  const remove = useCallback(
    (key: string) => {
      if (key in storage.current) {
        console.info(`delete ${key}`);
        clearTimeout(storage.current[key].timer);
        delete storage.current[key];
        setLastUpdated(new Date().getTime());
      }
    },
    [storage]
  );

  const keys = useMemo(() => Object.keys(storage.current), [
    storage,
    lastUpdated
  ]);

  useEffect(() => {
    const debug = setInterval(() => {
      if (keys.length > 0) {
        console.info(keys);
      }
    }, 1000);
    return () => {
      clearInterval(debug);
    };
  }, [keys]);

  return [get, put, remove, keys];
};

export default useCache;
