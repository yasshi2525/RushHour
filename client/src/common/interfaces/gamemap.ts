export interface Point {
  x: number;
  y: number;
}

export interface Edge {
  from: Point;
  to: Point;
}

export interface Coordinates {
  /**
   * 中心x座標(サーバにおけるマップ座標系)
   */
  cx: number;
  /**
   * 中心y座標(サーバにおけるマップ座標系)
   */
  cy: number;
  /**
   * 拡大率(クライアントウィンドウの幅が2^scaleに対応する)
   */
  scale: number;
  /**
   * 最後にsetCoordしたときズームしたか？(1=ズームした, 0=変更なし、-1=縮小した)
   */
  zoom: number;
}

export interface Chunk {
  x: number;
  y: number;
  scale: number;
}

export interface ResolveError {
  hasUnresolvedOwner?: boolean;
}

export function getChunkByPos(pos: Point, scale: number): Chunk {
  scale = Math.floor(scale);
  let interval = Math.pow(2, scale);
  return {
    x: Math.floor(pos.x / interval),
    y: Math.floor(pos.y / interval),
    scale: scale
  };
}

export function getChunkByScale(chunk: Chunk, offset: number): Chunk {
  return {
    x: Math.floor(chunk.x * Math.pow(2, -offset)),
    y: Math.floor(chunk.y * Math.pow(2, -offset)),
    scale: chunk.scale + offset
  };
}

export const config = {
  background: 0x263238,
  /**
   * サーバからマップ情報を取得する間隔 (ms)
   */
  interval: 1000,
  /**
   * 座標に変化があった際、位置合わせをするまでの遅延フレーム数
   */
  latency: 30,
  /**
   * 繰り返しアニメーションのフレーム数
   */
  round: 240,
  gamePos: {
    min: { x: -Math.pow(2, 15), y: -Math.pow(2, 15) },
    max: { x: Math.pow(2, 15), y: Math.pow(2, 15) },
    default: { x: 0, y: 0 }
  },
  scale: {
    min: 5,
    max: 16,
    default: 10
  },
  hsv: {
    saturation: 0.75,
    value: 1.0
  }
};

export function hueToRgb(H: number) {
  let S = config.hsv.saturation;
  let V = config.hsv.value;
  let C = V * S;
  let Hp = H / 60;
  let X = C * (1 - Math.abs((Hp % 2) - 1));

  let R = 0,
    G = 0,
    B = 0;
  if (0 <= Hp && Hp < 1) {
    [R, G, B] = [C, X, 0];
  }
  if (1 <= Hp && Hp < 2) {
    [R, G, B] = [X, C, 0];
  }
  if (2 <= Hp && Hp < 3) {
    [R, G, B] = [0, C, X];
  }
  if (3 <= Hp && Hp < 4) {
    [R, G, B] = [0, X, C];
  }
  if (4 <= Hp && Hp < 5) {
    [R, G, B] = [X, 0, C];
  }
  if (5 <= Hp && Hp < 6) {
    [R, G, B] = [C, 0, X];
  }

  let m = V - C;
  [R, G, B] = [R + m, G + m, B + m];

  R = Math.floor(R * 255);
  G = Math.floor(G * 255);
  B = Math.floor(B * 255);

  return [R, G, B];
}
