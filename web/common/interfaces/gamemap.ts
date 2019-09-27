export interface Point {
    x: number,
    y: number
};

export interface Edge {
    from: Point,
    to: Point
};

export interface Coordinates {
    /**
     * 中心x座標(サーバにおけるマップ座標系)
     */
    cx: number,
    /**
     * 中心y座標(サーバにおけるマップ座標系)
     */
    cy: number,
    /**
     * 拡大率(クライアントウィンドウの幅が2^scaleに対応する)
     */
    scale: number,
    /**
     * 最後にsetCoordしたときズームしたか？(1=ズームした, 0=変更なし、-1=縮小した)
     */
    zoom: number
};

export interface Chunk {
    x: number,
    y: number,
    scale: number
};

export interface ResolveError {
    hasUnresolvedOwner?: boolean
}

export function getChunkByPos(pos: Point, scale: number): Chunk {
    scale = Math.floor(scale);
    let interval = Math.pow(2, scale);
    return {
        x: Math.floor(pos.x / interval),
        y: Math.floor(pos.y / interval),
        scale: scale
    }
}

export function getChunkByScale(chunk: Chunk, offset: number): Chunk {
    return {
        x: Math.floor(chunk.x * Math.pow(2, -offset)),
        y: Math.floor(chunk.y * Math.pow(2, -offset)),
        scale: chunk.scale + offset
    }
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
        min: {x: -Math.pow(2, 15), y: -Math.pow(2, 15)}, 
        max: {x: Math.pow(2, 15), y: Math.pow(2, 15)},
        default: {x: 0, y: 0}
    },
    scale: { 
        min: 5, 
        max: 16, 
        default: 10
    }
};