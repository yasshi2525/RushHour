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
     * 最後にsetScaleしたときズームしたか？(1=ズームした, 0=変更なし、-1=縮小した)
     */
    zoom: number
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
        default: 10,
        delegate: 3
    },
    zIndices: [
        "residences",
        "companies",
        "rail_edges",
        "rail_nodes",
        "stations",
        "gates",
        "platforms",
        "rail_lines",
        "line_tasks",
        "players",
        "humans",
    ]
};