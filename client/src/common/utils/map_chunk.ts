import { useMemo } from "react";

export type Chunk = [number, number, number, number];

/**
 * 現在地点からの相対距離リスト
 */
const makeCube = (dlg: number, length: number) => {
  const buffer: Chunk[] = [];
  for (let i = 0; i <= length; i++) {
    for (let j = 0; j <= length; j++) {
      for (let k = 0; k <= length; k++) {
        buffer.push([i, j, k, dlg]);
      }
    }
  }
  return buffer;
};

const makeValidCube = (
  cube: Chunk[],
  [originChX, originChY, originFlatS, delegate]: Chunk,
  minScale: number,
  maxScale: number
) =>
  cube
    .map(
      ([dx, dy, ds]) =>
        [originChX + dx, originChY + dy, originFlatS + ds, delegate] as Chunk
    )
    .filter(
      ([chX, chY, flatS]) =>
        !(chX >> (flatS - minScale)) &&
        !(chY >> (flatS - minScale)) &&
        flatS <= maxScale
    );

type Handler = [Chunk, Chunk[]];

/**
 * ```
 * const [current, cube] = useCoreMapChunk(x, y, s, dlg, minScale, maxScale)
 * ```
 */
const useCoreMapChunk = (
  coordX: number,
  coordY: number,
  coordS: number,
  delegate: number,
  minScale: number,
  maxScale: number,
  length: number = 1
): Handler => {
  /**
   * 最小拡大率の単位でチャンク番号を決める
   */
  const chX = useMemo(() => Math.floor(coordX >> minScale), [coordX, minScale]);
  const chY = useMemo(() => Math.floor(coordY >> minScale), [coordY, minScale]);
  const flatS = useMemo(() => Math.floor(coordS), [coordS]);
  const current = useMemo(() => [chX, chY, flatS, delegate] as Chunk, [
    chX,
    chY,
    flatS,
    delegate
  ]);

  const cube = useMemo(() => makeCube(delegate, length), [delegate, length]);

  /**
   * 現在地点 `(chX, chY)` 周辺(左上起点) の左右上下前後 `2^length` チャンク。
   * ただし、マップ外のチャンクは除く
   */
  const validCube = useMemo(
    () => makeValidCube(cube, current, minScale, maxScale),
    [cube, current, minScale, maxScale]
  );

  return [current, validCube];
};

export default useCoreMapChunk;
