import { useCallback, useEffect, useState, useRef } from "react";
import useCache, { CacheStatus } from "./cache";

export type GraceHandler = {
  prepared: boolean;
  send: (required: string[]) => void;
};

export enum FlashStatus {
  PRIMARY_ACTIVE,
  PRIMARY_GRACE,
  SECONDARY_ACTIVE,
  SECONDARY_GRACE
}

const ACTIVE = 60 * 1000;
const GRACE = 10 * 1000;
const BULK_WAIT = 1 * 1000;

type Handlers<T> = [
  (key: string) => T | undefined,
  (key: string, value: T) => void,
  (key: string) => FlashStatus | undefined
];
/**
 * 主要な領域(`primary`)は長時間オブジェクトを保管し、
 * そうでない領域は短時間オブジェクトを保管する。
 * 主要な領域のオブジェクトは消える前に`grace`リストに追加される
 *
 *
 *            primary storage            == secondary ==
 * [onPut] -> Active -> Grace  -[save]-> Active -> Grace -> [onRemove]
 *                             -[move]->
 *                            <-[restore]-
 *
 * ```
 * const [get, put, status, onBulkGrace] = useFlash(keySet, handler);
 * ```
 */
const useFlash = <T>(
  boundary: string[],
  graceHandler: GraceHandler,
  onPut: (k: string, v: T) => void = () => {},
  onRemove: (k: string) => void = () => {},
  activeTime: number = ACTIVE,
  graceTime: number = GRACE
): Handlers<T> => {
  /**
   * `boundary` にないデータは短時間のみキャッシュする。
   * (`boundary`が頻繁に変わっても、外部からのリソース読み込み回数を減らすため)
   */
  const [
    getSecondary,
    save,
    removeSecondary,
    getInSecondary,
    ,
    statusSecondary
  ] = useCache<T>(() => {}, onRemove);

  const [graceRevision, setGraceRevision] = useState(0);
  const notifyGrace = useCallback(() => {
    setGraceRevision(g => g + 1);
  }, []);

  /**
   * データの内、primary の中にあるものは長時間保持する。
   * 期限切れになった場合、`secondary`に移動する
   */
  const [
    getPrimary,
    putPrimary,
    move, // 削除すると onExpired で save コールされ、secondaryに移動する
    ,
    bulkMove,
    statusPrimary
  ] = useCache<T>(notifyGrace, save, activeTime, graceTime);

  const restore = useCallback(
    (key: string) => {
      const cache = removeSecondary(key, false); // onRemoveを呼び出さないため
      if (cache !== undefined) {
        putPrimary(key, cache);
      }
    },
    [removeSecondary, putPrimary]
  );

  const override = useCallback(
    (key: string, data: T) => {
      if (boundary.includes(key)) {
        // delete old data
        move(key, false); // secondary に移動させないため
        removeSecondary(key, false); // 下の命令で二重に通知することを防ぐため
        onRemove(key);
        // put new data
        putPrimary(key, data);
      } else {
        // delete old data
        removeSecondary(key);
        // put new data
        save(key, data);
      }
    },
    [move, removeSecondary, onRemove, putPrimary, onPut]
  );

  const get = useCallback(
    (key: string) => {
      const data = getPrimary(key);
      return data !== undefined ? data : getSecondary(key);
    },
    [getPrimary, getSecondary]
  );

  /**
   * when on put [p=primary, s=secondary, A=Active, G=Grace]
   *
   * bnd | before | after | event
   * ----+-------+-------+-------
   * in  | none   | p,A   | onPut
   * in  | p,A    | p,A   | override(onRemove) & onPut
   * in  | p,G    | p,A   | override(onRemove) & onPut
   * in  | s,A    | p,A   | override(onRemove) & onPut
   * in  | s,G    | p,A   | override(onRemove) & onPut
   * out | none   | s,A   | onPut
   * out | s,A    | s,A   | override(onRemove) & onPut
   * out | s,G    | s,A   | override(onRemove) & onPut
   */
  const put = useCallback(
    (key: string, data: T) => {
      if (get(key) !== undefined) {
        override(key, data);
      } else {
        if (boundary.includes(key)) {
          putPrimary(key, data);
        } else {
          save(key, data);
        }
      }
      onPut(key, data);
    },
    [boundary, putPrimary, save, onPut]
  );

  const status = useCallback(
    (key: string) => {
      const primary = getPrimary(key);
      if (primary) {
        switch (statusPrimary(key)) {
          case CacheStatus.ACTIVE:
            return FlashStatus.PRIMARY_ACTIVE;
          case CacheStatus.GRACE:
            return FlashStatus.PRIMARY_GRACE;
          default:
            return undefined;
        }
      } else {
        const secondary = getSecondary(key);
        if (secondary) {
          switch (statusSecondary(key)) {
            case CacheStatus.ACTIVE:
              return FlashStatus.SECONDARY_ACTIVE;
            case CacheStatus.GRACE:
              return FlashStatus.SECONDARY_GRACE;
          }
        }
        return undefined;
      }
    },
    [getPrimary, getSecondary, statusPrimary, statusSecondary]
  );

  /**
   * when change boundary [p=primary, s=secondary, A=Active, G=Grace]
   *
   * key            | place          | fired
   * before | after | before | after | event
   * -------+-------+--------+-------+-------
   * in     | in    | *      | *     |
   * in     | out   | p,A    | s,A   | move
   * in     | out   | p,G    | s,A   | move
   * in     | out   | s,*    | s,*   |
   * out    | in    | s,*    | p,A   | restore
   * out    | out   | s,*    | s,*   |
   */
  useEffect(() => {
    console.info("clean up out bound object");
    bulkMove(boundary);
    getInSecondary(boundary).forEach(([key]) => restore(key));
  }, [boundary, bulkMove, getInSecondary, restore]);

  /**
   * 更新が必要なデータが発生するイベントをキャプチャし、
   * 一定時間後に通知する
   * boundary変更と、primaryの時間経過が該当
   *
   * データの更新が必要なキー一覧は以下。
   *
   * boundary | primary | secondary | update?
   * ---------+---------+-----------+---------
   * in       | no key  | no key    | true
   * in       | no key  | active    | true
   * in       | no key  | grace     | true
   * in       | active  | *         | false
   * in       | grace   | *         | true
   * out      | *       | *         | false
   */
  useEffect(() => {
    console.info(`notify grace list after ${BULK_WAIT} ms`);

    const bulk = setTimeout(() => {
      if (graceHandler.prepared) {
        const list = boundary.filter(
          key => statusPrimary(key) !== CacheStatus.ACTIVE
        );
        if (list.length) {
          graceHandler.send(list);
        }
      } else {
        console.warn(`graceHandler is not prepared, pending to request`);
      }
    }, BULK_WAIT);
    return () => {
      console.info(`canceled to notify grace list`);
      clearTimeout(bulk);
    };
  }, [graceRevision, boundary, graceHandler, statusPrimary]);

  return [get, put, status];
};

export default useFlash;
