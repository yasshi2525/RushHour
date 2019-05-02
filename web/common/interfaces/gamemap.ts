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
    interval: 1000, // ms
    gamePos: { 
        min: {x: -1000, y: -1000}, 
        max: {x: 1000, y: 1000},
        default: {x: 0, y: 0}
    },
    scale: { 
        min: 0, 
        max: 16, 
        default: 10
    }
};