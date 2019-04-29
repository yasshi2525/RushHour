export interface Point {
    x: number,
    y: number
};

export interface Edge {
    from: Point,
    to: Point
};

export const config = {
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