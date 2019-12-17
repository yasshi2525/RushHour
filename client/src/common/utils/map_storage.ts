import { useMemo } from "react";
import useFlash, { FlashStatus, GraceHandler } from "./flash";
import { CoreMap } from "./map_core";
import { Chunk } from "./map_chunk";

const X = 0;
const Y = 1;
const S = 2;
const D = 3;

export const hash = (chunk: Chunk, min: number) =>
  `${chunk[X] << min}_${chunk[Y] << min}_${chunk[S]}_${chunk[D]}`;

const bulkHash = (targets: Chunk[], minScale: number) =>
  targets.map(ch => hash(ch, minScale));

type Handlers = [
  (key: string) => CoreMap | undefined,
  (key: string, value: CoreMap) => void,
  (key: string) => FlashStatus | undefined,
  string,
  string[]
];

/**
 * マップデータを格納する。一定時間経過後は消す。
 * 範囲外のものは短期間保持したあと消す。
 * ```
 * const [get, put, status, key, keyAll] = useCoreMapStorage(current, cube, minS, graceHandler);
 * ```
 */
const useCoreMapStorage = (
  current: Chunk,
  cube: Chunk[],
  minScale: number,
  graceHandler: GraceHandler,
  onAdd: (k: string, v: CoreMap) => void,
  onDelete: (key: string) => void
): Handlers => {
  const key = useMemo(() => hash(current, minScale), [current, minScale]);

  /**
   * 長期間保存するチャンク一覧の作成
   */
  const primaries = useMemo(() => bulkHash(cube, minScale), [cube, minScale]);

  /**
   * サーバから受信したデータを格納する。一定時間経過後は消す。
   * 範囲外のものは短期間保持したあと消す。
   * (受信まで時間がかかり、その間に`coord`が変わった場合が該当)
   */
  const [get, put, status] = useFlash<CoreMap>(
    primaries,
    graceHandler,
    onAdd,
    onDelete
  );

  return [get, put, status, key, primaries];
};

export default useCoreMapStorage;
