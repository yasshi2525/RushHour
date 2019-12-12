import { useMemo, useCallback } from "react";
import { Coordinates } from "interfaces";
import { Chunk, CoreMap } from "./map_core";
import useFlash from "./flash";

const X = 0;
const Y = 1;
const S = 2;

const encode = (chunk: Chunk) => chunk.join("_");
const decode = (hash: string) => hash.split("_").map(v => parseInt(v)) as Chunk;

type Handlers = [
  (key: Chunk) => CoreMap | undefined,
  (key: Chunk, value: CoreMap) => void,
  Chunk[],
  Chunk,
  Chunk[]
];

/**
 * マップデータを格納する。一定時間経過後は消す。
 * 範囲外のものは短期間保持したあと消す。
 * ```
 * const [get, put, keys, currentKey, expiredKeys] = useCoreMapStorage(coord, dlg, minS, maxS);
 * ```
 */
const useCoreMapStorage = (
  coord: Coordinates,
  delegate: number,
  minScale: number,
  maxScale: number,
  length: number = 1
): Handlers => {
  /**
   * チャンクとして保存するため、小数点以下を切り捨て
   */
  const flatCoord = useMemo<Coordinates>(
    () => [Math.floor(coord[X]), Math.floor(coord[Y]), Math.floor(coord[S])],
    [coord]
  );

  /**
   * 現在地点からの相対距離リスト
   */
  const cube = useMemo<Coordinates[]>(() => {
    const buffer: Coordinates[] = [];
    for (let x = 0; x <= length; x++) {
      for (let y = 0; y <= length; y++) {
        for (let s = 0; s <= length; s++) {
          buffer.push([x, y, s]);
        }
      }
    }
    return buffer;
  }, [length]);

  /**
   * 現在地点 `coord` 周辺(左上起点) の左右上下前後 `2^length` チャンク。
   * ただし、マップ外のチャンクは除く
   */
  const targets = useMemo<Coordinates[]>(
    () =>
      cube
        .map(
          ([dx, dy, ds]) =>
            [
              flatCoord[X] + dx,
              flatCoord[Y] + dy,
              flatCoord[S] + ds
            ] as Coordinates
        )
        .filter(
          ([x, y, s]) =>
            !(x >> (maxScale - minScale)) &&
            !(y >> (maxScale - minScale)) &&
            s <= maxScale
        ),
    [maxScale, minScale, cube, flatCoord]
  );

  /**
   * 長期間保存するチャンク一覧の作成
   */
  const targetKeys = useMemo(
    () => targets.map(([x, y, s]) => encode([x, y, s, delegate])),
    [targets]
  );

  /**
   * サーバから受信したデータを格納する。一定時間経過後は消す。
   * 範囲外のものは短期間保持したあと消す。
   * (受信まで時間がかかり、その間に`coord`が変わった場合が該当)
   */
  const [get, put, expired] = useFlash<CoreMap>(targetKeys);

  const getChunk = useCallback((chunk: Chunk) => get(encode(chunk)), [get]);
  const putChunk = useCallback(
    (chunk: Chunk, data: CoreMap) => put(encode(chunk), data),
    [put]
  );

  const chunkList = useMemo(
    () => targets.map(([x, y, s]) => [x, y, s, delegate] as Chunk),
    [targets]
  );

  const currentChunk = useMemo<Chunk>(
    () => [flatCoord[X], flatCoord[Y], flatCoord[S], delegate],
    [flatCoord, delegate]
  );
  const expiredChunk = useMemo(() => expired.map(decode), [expired]);

  return [getChunk, putChunk, chunkList, currentChunk, expiredChunk];
};

export default useCoreMapStorage;
