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
     * 拡大率(クライエントウィンドウの幅が2^scaleに対応する)
     */
    scale: number,
}

export const config = {
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
        min: 0, 
        max: 16, 
        default: 10
    }
};